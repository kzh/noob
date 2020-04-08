package message

import (
	"log"
	"net/url"
	"os"

	"github.com/streadway/amqp"
)

var rabbitmq *amqp.Connection

func init() {
	log.Println("Connecting to RabbitMQ...")

	var addr url.URL
	addr.Scheme = "amqp"
	addr.User = url.UserPassword(
		"user",
		os.Getenv("RABBITMQ_PASSWORD"),
	)
	addr.Host = "noob-rabbitmq:5672"
	addr.Path = "/"

	u := addr.String()

	var err error
	rabbitmq, err = amqp.Dial(u)
	if err != nil {
		panic(err)
	}

	log.Println("Connected to RabbitMQ.")
}
