package rabbitmq

import (
	amqp "github.com/rabbitmq/amqp091-go"
)

type Consumer struct {
	conn  *amqp.Connection
	ch    *amqp.Channel
	queue string
}

func NewConsumer(
	conn *amqp.Connection,
	ch *amqp.Channel,
	queue string,
) *Consumer {

	return &Consumer{
		conn:  conn,
		ch:    ch,
		queue: queue,
	}
}

func (c *Consumer) Consume(
	handler func([]byte) error,
) error {

	msgs, err := c.ch.Consume(
		c.queue,
		"",
		false,
		false,
		false,
		false,
		nil,
	)

	if err != nil {
		return err
	}

	for msg := range msgs {

		err := handler(msg.Body)

		if err != nil {
			msg.Nack(false, true)
			continue
		}

		msg.Ack(false)
	}

	return nil
}
