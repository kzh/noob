package queue

import (
	"encoding/json"
	"log"
	"net/url"
	"os"

	"github.com/streadway/amqp"

	"github.com/kzh/noob/lib/model"
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

func queue() (ch *amqp.Channel, q amqp.Queue, e error) {
	ch, err := rabbitmq.Channel()
	if err != nil {
		e = err
		return
	}

	q, err = ch.QueueDeclare(
		"submissions",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		ch.Close()
	}

	return ch, q, err
}

func Schedule(s model.Submission) error {
	ch, q, err := queue()
	if err != nil {
		return err
	}
	defer ch.Close()

	body, err := json.Marshal(s)
	if err != nil {
		return err
	}
	log.Println(string(body))

	publishing := amqp.Publishing{
		DeliveryMode: amqp.Persistent,
		ContentType:  "text/plain",
		Body:         body,
	}
	return ch.Publish("", q.Name, false, false, publishing)
}

func Poll() (<-chan amqp.Delivery, error) {
	ch, q, err := queue()
	if err != nil {
		return nil, err
	}

	msgs, err := ch.Consume(
		q.Name,
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	return msgs, err
}
