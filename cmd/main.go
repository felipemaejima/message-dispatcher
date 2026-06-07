package main

import (
	"encoding/json"
	"log"
	"os"

	"github.com/felipemaejima/message-dispatcher/internal/application"
	"github.com/felipemaejima/message-dispatcher/internal/domain"
	"github.com/felipemaejima/message-dispatcher/internal/ports"
	"github.com/felipemaejima/message-dispatcher/internal/providers/consumers/rabbitmq"
	"github.com/felipemaejima/message-dispatcher/internal/providers/logger"
	"github.com/felipemaejima/message-dispatcher/internal/providers/notifiers/resend"
	"github.com/joho/godotenv"
)

func main() {
	appLogger, err := logger.New()

	if err != nil {
		log.Fatal(err)
	}

	appLogger.Println(
		"application started",
	)

	err = godotenv.Load(".env")

	if err != nil {
		appLogger.Println(
			"error loading env:",
			err.Error(),
		)
		log.Fatal(err)
	}

	consumer := rabbitmq.NewConsumer(
		os.Getenv("MESSAGES_QUEUE"),
	)

	dispatcher := application.NewDispatcher(
		map[string]ports.Notifier{
			/**
			 * Lista de canais de envio da mensagem e suas implementações
			 */
			"resend": resend.NewNotifier(),
		},
	)

	// todo adicionar retorno e envio em cima do retorno (ver se a abordagem eh valido ou se deve usar channel)
	err = consumer.Consume(
		func(body []byte) error {

			appLogger.Printf(
				"message received: %s",
				string(body),
			)

			var notification domain.Notification

			err := json.Unmarshal(
				body,
				&notification,
			)

			if err != nil {
				appLogger.Printf(
					"invalid payload: %s",
					err.Error(),
				)
				return err
			}

			dispatchErr := dispatcher.Dispatch(notification)

			if dispatchErr == nil {
				appLogger.Printf(
					"notification dispatched channel=%s",
					notification.Channel,
				)
			} else {
				appLogger.Println(
					"error sending message:",
					dispatchErr.Error(),
				)
			}

			return dispatchErr
		},
	)

	if err != nil {
		appLogger.Println(
			"error in consumer:",
			err.Error(),
		)
		log.Fatal(err)
	}
}
