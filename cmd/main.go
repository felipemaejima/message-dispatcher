package main

import (
	"encoding/json"
	"log"
	"os"

	"github.com/felipemaejima/message-dispatcher/internal/application"
	"github.com/felipemaejima/message-dispatcher/internal/domain"
	"github.com/felipemaejima/message-dispatcher/internal/ports"
	"github.com/felipemaejima/message-dispatcher/internal/providers/consumers/rabbitmq"
	"github.com/felipemaejima/message-dispatcher/internal/providers/notifiers/resend"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load(".env")

	if err != nil {
		log.Fatal(err)
		return
	}

	consumer := rabbitmq.NewConsumer(
		os.Getenv("MESSAGES_QUEUE"),
	)

	dispatcher := application.NewDispatcher(
		map[string]ports.Notifier{
			/**
			 * Lista de canais de envio da mensagem e suas implementações
			 */
			"email": resend.NewNotifier(), // todo channel deve ser o nome do serviço (?)
		},
	)

	// todo adicionar retorno e envio em cima do retorno
	err = consumer.Consume(
		func(body []byte) error {

			var notification domain.Notification

			err := json.Unmarshal(
				body,
				&notification,
			)

			if err != nil {
				return err
			}

			return dispatcher.Dispatch(
				notification,
			)
		},
	)

	if err != nil {
		log.Fatal(err)
	}
}
