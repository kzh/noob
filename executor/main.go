package main

import (
	"archive/tar"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"strconv"
	"strings"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/client"
	"github.com/streadway/amqp"

	noobdb "github.com/kzh/noob/pkg/database"
	"github.com/kzh/noob/pkg/message"
	"github.com/kzh/noob/pkg/model"
)

var dock *client.Client

func buildImageContext(code string) (io.Reader, error) {
	raw, err := ioutil.ReadFile("image.tar")
	if err != nil {
		return nil, err
	}

	trim := len(raw) - 1024
	buf := bytes.NewBuffer(raw[:trim])

	w := tar.NewWriter(buf)
	header := &tar.Header{
		Name: "main.go",
		Size: int64(len(code)),
	}
	err = w.WriteHeader(header)
	if err != nil {
		return nil, err
	}
	_, err = w.Write([]byte(code))
	if err != nil {
		return nil, err
	}
	if err := w.Close(); err != nil {
		return nil, err
	}

	return buf, nil
}

func buildImage(id string, buildContext io.Reader) error {
	ctx := context.Background()
	res, err := dock.ImageBuild(
		ctx,
		buildContext,
		types.ImageBuildOptions{
			Remove:      true,
			ForceRemove: true,
			Tags:        []string{id},
			Context:     buildContext,
			Dockerfile:  "Dockerfile",
		},
	)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	b, _ := ioutil.ReadAll(res.Body)
	log.Println(string(b))

	return nil
}

func prepareContainer(uid string) (string, error) {
	ctx := context.Background()
	resp, err := dock.ContainerCreate(
		ctx,
		&container.Config{
			Image:           uid,
			NetworkDisabled: true,
		},
		&container.HostConfig{
			Resources: container.Resources{
				Memory: 4*1024*1024 + 1,
			},
		},
		nil, uid,
	)
	if err != nil {
		return "", err
	}

	err = dock.ContainerStart(
		ctx, resp.ID,
		types.ContainerStartOptions{},
	)
	if err != nil {
		return "", err
	}

	return resp.ID, nil
}

func test(uid, in, out string) (string, error) {
	ctx := context.Background()
	resp, err := dock.ContainerExecCreate(
		ctx, uid,
		types.ExecConfig{
			Cmd:          []string{"./ex"},
			Tty:          true,
			AttachStdin:  true,
			AttachStderr: true,
			AttachStdout: true,
		},
	)
	if err != nil {
		return "", err
	}

	exec, err := dock.ContainerExecAttach(
		ctx, resp.ID,
		types.ExecConfig{Tty: true},
	)
	if err != nil {
		return "", err
	}
	defer exec.Close()

	_, err = exec.Conn.Write([]byte(in))
	if err != nil {
		return "", err
	}

	_, err = fmt.Fscanf(exec.Reader, in+out)
	if err != nil {
		line, _ := ioutil.ReadAll(exec.Reader)
		log.Println("FAILED")
		fmt.Printf("Error:\n%s\n", string(line))
		return string(line), nil
	}

	log.Println("PASSED")
	return "", err
}

func clean(uid string) error {
	ctx := context.Background()
	err := dock.ContainerRemove(
		ctx, uid,
		types.ContainerRemoveOptions{
			RemoveVolumes: true,
			Force:         true,
		},
	)
	if err != nil {
		return err
	}

	_, err = dock.ImageRemove(
		ctx, uid,
		types.ImageRemoveOptions{
			Force:         true,
			PruneChildren: true,
		},
	)

	return err
}

func sanitize(in string) string {
	return strings.Map(func(r rune) rune {
		if r == '\n' || r == '\r' {
			return -1
		}

		return r
	}, in)
}

func handle(msg amqp.Delivery) {
	var submission model.Submission
	err := json.Unmarshal(msg.Body, &submission)
	if err != nil {
		log.Println(err)
		return
	}

	log.Println("Incoming submission.")
	log.Println("ID: " + submission.ID)

	ctx, err := buildImageContext(submission.Code)
	if err != nil {
		log.Println(err)
		return
	}

	err = buildImage(submission.ID, ctx)
	if err != nil {
		log.Println(err)
		return
	}

	cid, err := prepareContainer(submission.ID)
	if err != nil {
		log.Println(err)

		_, err = dock.ImagesPrune(
			context.Background(),
			filters.Args{},
		)

		var result model.SubmissionResult
		result.Stage = "Compile"
		result.Status = "FAILED"
		message.Publish(submission.ID, result)
		return
	}

	probio, err := noobdb.IOProblem(submission.ProblemID)
	inputs := strings.Split(probio.In, "---")
	outputs := strings.Split(probio.Out, "---")

	for i, in := range inputs {
		resp, err := test(cid, sanitize(in)+"\n", sanitize(outputs[i]))

		var result model.SubmissionResult
		result.Stage = strconv.Itoa(i + 1)
		result.Status = "PASSED"
		if resp != "" || err != nil {
			result.Status = "FAILED"
			log.Printf("%s %#v\n", resp, err)
		}

		message.Publish(submission.ID, result)
	}

	err = clean(submission.ID)
	if err != nil {
		log.Println(err)
		return
	}

	log.Println("Finishing handling submission.")
}

func main() {
	log.Println("Noob: Executor Worker is starting...")

	var err error
	dock, err = client.NewEnvClient()
	if err != nil {
		panic(err)
	}

	msgs, err := message.Poll()
	if err != nil {
		panic(err)
	}

	for msg := range msgs {
		go handle(msg)
	}
}
