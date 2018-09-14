package main

import (
	"log"

	"github.com/streadway/amqp"

	"github.com/kzh/noob/lib/queue"
)

func handle(msg amqp.Delivery) {
	log.Println(string(msg.Body))
}

func main() {
	msgs, err := queue.Poll()
	if err != nil {
		panic(err)
	}

	for msg := range msgs {
		handle(msg)
	}
}
