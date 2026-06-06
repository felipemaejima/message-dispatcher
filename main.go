package main

import (
	"encoding/json"
	"log"
	"os"

	"github.com/joho/godotenv"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/resend/resend-go/v3"
)

type Email struct {
	From    string   `json:"from"`
	To      []string `json:"to"`
	Subject string   `json:"subject"`
	Bcc     []string `json:"bcc"`
	Cc      []string `json:"cc"`
	ReplyTo string   `json:"replyTo"`
	Html    string   `json:"html"`
	Text    string   `json:"text"`
}

func main() {

	err := godotenv.Load()

	if err != nil {
		panic(err)
	}

	brokerUrl := os.Getenv("BROKER_URL")
	messagesQueue := os.Getenv("MESSAGES_QUEUE")

	conn, err := amqp.Dial(brokerUrl)

	if err != nil {
		log.Fatal(err)
	}

	defer conn.Close()

	ch, err := conn.Channel()

	if err != nil {
		log.Fatal(err)
	}

	defer ch.Close()

	msgs, err := ch.Consume(
		messagesQueue, // nome da fila
		"",            // consumer name
		false,         // auto ack
		false,
		false,
		false,
		nil,
	)

	if err != nil {
		log.Fatal(err)
	}

	log.Println("Aguardando mensagens...")

	for msg := range msgs {
		log.Printf("%s\n", msg.Body)

		var email Email

		err := json.Unmarshal(msg.Body, &email)

		if err != nil {
			log.Println(err)
			msg.Nack(false, true)
			return
		}

		log.Println(email)

		// Email{
		// 	To: []string{
		// 		"felipemaejima@gmail.com",
		// 	},
		// 	Subject: "Teste",
		// 	Html:    "<h1>Hello</h1>",
		// }

		_, err = dispatchEmail(email)

		if err != nil {
			log.Println(err)
			return
		}

		msg.Ack(false)
	}

}

// todo email service
func dispatchEmail(data Email) (string, error) {
	apiKey := os.Getenv("RESEND_API_KEY")
	emailDomain := os.Getenv("EMAIL_DOMAIN")

	client := resend.NewClient(apiKey)

	params := &resend.SendEmailRequest{
		From:    emailDomain,
		To:      data.To,
		Subject: data.Subject,
		Html:    data.Html,
		Text:    data.Text,
	}

	sent, err := client.Emails.Send(params)

	if err != nil {
		return "", err
	}

	return sent.Id, nil
}
