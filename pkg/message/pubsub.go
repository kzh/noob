package message

import (
	"encoding/json"

	"github.com/streadway/amqp"

	"github.com/kzh/noob/pkg/model"
)

func Subscribe(pubsub string) (<-chan amqp.Delivery, error) {
	ch, err := rabbitmq.Channel()
	if err != nil {
		return nil, err
	}

	err = ch.ExchangeDeclare(
		pubsub, "fanout",
		true, false, false, false, nil,
	)
	if err != nil {
		return nil, err
	}

	q, err := ch.QueueDeclare(
		"",
		false, false, true, false, nil,
	)
	if err != nil {
		return nil, err
	}

	err = ch.QueueBind(
		q.Name,
		"", pubsub, false, nil,
	)
	if err != nil {
		return nil, err
	}

	msgs, err := ch.Consume(
		q.Name, "",
		true, false, false, false, nil,
	)
	return msgs, err
}

func Publish(pubsub string, res model.SubmissionResult) error {
	ch, err := rabbitmq.Channel()
	if err != nil {
		return err
	}

	err = ch.ExchangeDeclare(
		pubsub, "fanout",
		true, false, false, false, nil,
	)
	if err != nil {
		return err
	}

	body, err := json.Marshal(res)
	if err != nil {
		return err
	}

	publishing := amqp.Publishing{
		DeliveryMode: amqp.Persistent,
		ContentType:  "text/plain",
		Body:         body,
	}
	return ch.Publish(pubsub, "", false, false, publishing)
}
