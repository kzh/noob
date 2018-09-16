package main

import (
	"archive/tar"
	"bytes"
	"context"
	"encoding/json"
	"io/ioutil"
	"log"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/streadway/amqp"

	"github.com/kzh/noob/lib/model"
	"github.com/kzh/noob/lib/queue"
)

var dock *client.Client

func buildImage(code string) error {
	raw, err := ioutil.ReadFile("image.tar")
	if err != nil {
		return err
	}

	buf := bytes.NewBuffer(raw[:len(raw)-1024])
	w := tar.NewWriter(buf)
	header := &tar.Header{
		Name: "main.go",
		Size: int64(len(code)),
	}
	err = w.WriteHeader(header)
	if err != nil {
		return err
	}
	_, err = w.Write([]byte(code))
	if err != nil {
		return err
	}
	if err := w.Close(); err != nil {
		return err
	}

	ctx := context.Background()
	res, err := dock.ImageBuild(
		ctx,
		buf,
		types.ImageBuildOptions{
			Context:    buf,
			Dockerfile: "Dockerfile",
			Remove:     true,
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

func handle(msg amqp.Delivery) {
	var submission model.Submission
	err := json.Unmarshal(msg.Body, &submission)
	if err != nil {
		panic(err)
	}

	log.Println("Building image...")
	err = buildImage(submission.Code)
	if err != nil {
		log.Println(err)
	}
	log.Println("Finished building image?")
}

func main() {
	log.Println("Noob: Executor Worker is starting...")

	var err error
	dock, err = client.NewEnvClient()
	if err != nil {
		panic(err)
	}

	msgs, err := queue.Poll()
	if err != nil {
		panic(err)
	}

	for msg := range msgs {
		handle(msg)
	}
}
