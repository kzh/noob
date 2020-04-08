package message

import (
	"encoding/json"
	"log"

	"github.com/streadway/amqp"

	"github.com/kzh/noob/pkg/model"
)

func Schedule(s model.Submission) error {
	ch, err := rabbitmq.Channel()
	if err != nil {
		return err
	}
	defer ch.Close()

	_, err = ch.QueueDeclare(
		"submissions",
		true, false, false, false, nil,
	)
	if err != nil {
		return err
	}

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
	return ch.Publish("", "submissions", false, false, publishing)
}

func Poll() (<-chan amqp.Delivery, error) {
	ch, err := rabbitmq.Channel()
	if err != nil {
		return nil, err
	}

	_, err = ch.QueueDeclare(
		"submissions",
		true, false, false, false, nil,
	)
	if err != nil {
		return nil, err
	}

	msgs, err := ch.Consume(
		"submissions",
		"", true, false, false, false, nil,
	)
	return msgs, err
}
