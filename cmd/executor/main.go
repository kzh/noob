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
	"sync"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/client"
	opentracing "github.com/opentracing/opentracing-go"
	otlog "github.com/opentracing/opentracing-go/log"
	"github.com/streadway/amqp"

	noobdb "github.com/kzh/noob/pkg/database"
	"github.com/kzh/noob/pkg/message"
	"github.com/kzh/noob/pkg/model"
	"github.com/kzh/noob/pkg/tracing"
)

var dock *client.Client

func buildImageContext(parent opentracing.Span, code string) (io.Reader, error) {
	span := opentracing.StartSpan(
		"buildImageContext",
		opentracing.ChildOf(parent.Context()),
	)
	defer span.Finish()

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

func buildImage(parent opentracing.Span, id string, buildContext io.Reader) error {
	span := opentracing.StartSpan(
		"buildImage",
		opentracing.ChildOf(parent.Context()),
	)
	defer span.Finish()

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

func prepareContainer(parent opentracing.Span, uid string) (string, error) {
	span := opentracing.StartSpan(
		"prepareContainer",
		opentracing.ChildOf(parent.Context()),
	)
	defer span.Finish()

	ctx := context.Background()
	resp, err := dock.ContainerCreate(
		ctx,
		&container.Config{
			Image:           uid,
			NetworkDisabled: true,
		},
		&container.HostConfig{
			Resources: container.Resources{
				MemoryReservation: 4*1024*1024 + 1,
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

func test(parent opentracing.Span, container, problem, submission string) {
	span := opentracing.StartSpan(
		"test",
		opentracing.ChildOf(parent.Context()),
	)
	defer span.Finish()

	probio, err := noobdb.IOProblem(problem)
	if err != nil {
		var result model.SubmissionResult
		result.Stage = "Internal"
		result.Status = "FAILED"
		message.Publish(submission, result)
		return
	}

	inputs := strings.Split(probio.In, "---")
	outputs := strings.Split(probio.Out, "---")

	var wg sync.WaitGroup

	for i, input := range inputs {
		wg.Add(1)

		i := i
		in, out := sanitize(input)+"\n", sanitize(outputs[i])
		go func() {
			defer wg.Done()

			resp, err := exec(span, container, in, out)

			var result model.SubmissionResult
			result.Stage = strconv.Itoa(i + 1)
			result.Status = "PASSED"
			if resp != "" || err != nil {
				result.Status = "FAILED"
				log.Printf("%s %#v\n", resp, err)
			}

			message.Publish(submission, result)
		}()
	}

	wg.Wait()
}

func exec(parent opentracing.Span, uid, in, out string) (string, error) {
	span := opentracing.StartSpan(
		"exec",
		opentracing.ChildOf(parent.Context()),
	)
	defer span.Finish()

	span.LogFields(
		otlog.String("input", in),
	)

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

func clean(parent opentracing.Span, uid string) error {
	span := opentracing.StartSpan(
		"clean",
		opentracing.ChildOf(parent.Context()),
	)
	defer span.Finish()

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
	span := opentracing.StartSpan("handleSubmission")
	defer span.Finish()

	var submission model.Submission
	err := json.Unmarshal(msg.Body, &submission)
	if err != nil {
		log.Println(err)
		return
	}

	span.LogFields(
		otlog.String("id", submission.ID),
		otlog.String("problem", submission.ProblemID),
		otlog.String("code", submission.Code),
	)

	log.Println("Incoming submission.")
	log.Println("ID: " + submission.ID)

	ctx, err := buildImageContext(span, submission.Code)
	if err != nil {
		log.Println(err)
		return
	}

	err = buildImage(span, submission.ID, ctx)
	if err != nil {
		log.Println(err)
		return
	}

	cid, err := prepareContainer(span, submission.ID)
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

	test(span, cid, submission.ProblemID, submission.ID)

	err = clean(span, submission.ID)
	if err != nil {
		log.Println(err)
		return
	}

	log.Println("Finishing handling submission.")
}

func main() {
	log.Println("Noob: Executor Worker is starting...")

	closer, err := tracing.InitJaeger()
	if err != nil {
		panic(err)
	}
	defer closer.Close()

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
