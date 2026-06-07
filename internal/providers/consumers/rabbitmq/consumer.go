package rabbitmq

import (
	"os"

	amqp "github.com/rabbitmq/amqp091-go"
)

type Consumer struct {
	queue string
}

func NewConsumer(
	queue string,
) *Consumer {

	return &Consumer{
		queue: queue,
	}
}

func (c *Consumer) Consume(
	handler func([]byte) error,
) error {

	conn, ch, err := connect()

	if err != nil {
		return err
	}

	defer ch.Close()
	defer conn.Close()

	msgs, err := ch.Consume(
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
      // todo jogar para fila dlq
			msg.Nack(false, false)
			continue
		}

		msg.Ack(false)
	}

	return nil
}

func connect() (*amqp.Connection, *amqp.Channel, error) {
	brokerUrl := os.Getenv("RABBITMQ_BROKER_URL")

	conn, err := amqp.Dial(brokerUrl)

	if err != nil {
		return nil, nil, err
	}

	ch, err := conn.Channel()

	if err != nil {
		conn.Close()
		return nil, nil, err
	}

	return conn, ch, nil
}
