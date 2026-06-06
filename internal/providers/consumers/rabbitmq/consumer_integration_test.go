package rabbitmq

import (
	"os"
	"testing"
	"time"

	"github.com/joho/godotenv"
	amqp "github.com/rabbitmq/amqp091-go"
)

func TestConsumeMessage(t *testing.T) {

	err := godotenv.Load("../../../../.env")

	if err != nil {
		t.Fatalf("erro carregando env: %v", err)
	}

	queue := "test.notifications"

	conn, err := amqp.Dial(
		os.Getenv("RABBITMQ_BROKER_URL"),
	)

	if err != nil {
		t.Fatalf("erro conexão: %v", err)
	}

	defer conn.Close()

	ch, err := conn.Channel()

	if err != nil {
		t.Fatalf("erro channel: %v", err)
	}

	defer ch.Close()

	_, err = ch.QueueDeclare(
		queue,
		true,
		false,
		false,
		false,
		nil,
	)

	if err != nil {
		t.Fatalf("erro queue: %v", err)
	}

	payload := []byte(`{
		"channel":"email",
		"subject":"teste"
	}`)

	err = ch.Publish(
		"",
		queue,
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        payload,
		},
	)

	if err != nil {
		t.Fatalf("erro publish: %v", err)
	}

	consumer := NewConsumer(queue)

	received := make(chan []byte)

	go func() {

		err := consumer.Consume(
			func(body []byte) error {

				received <- body

				return nil
			},
		)

		if err != nil {
			t.Error(err)
		}
	}()

	select {

	case body := <-received:

		if string(body) != string(payload) {
			t.Fatalf(
				"payload diferente\nesperado: %s\nrecebido: %s",
				string(payload),
				string(body),
			)
		}

	case <-time.After(5 * time.Second):

		t.Fatal("timeout aguardando mensagem")
	}
}
