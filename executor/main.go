package main

import (
	"archive/tar"
	"bytes"
	"context"
	"encoding/json"
	"io"
	"io/ioutil"
	"log"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/streadway/amqp"

	"github.com/kzh/noob/lib/model"
	"github.com/kzh/noob/lib/queue"
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
			NoCache:     true,
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
