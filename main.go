package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/resend/resend-go/v3"
)

// todo adicionar o pacote do json para logs e playload
type Email struct {
	From    string
	To      []string
	Subject string
	Bcc     []string
	Cc      []string
	ReplyTo string
	Html    string
	Text    string
}

func main() {

	err := godotenv.Load()

	if err != nil {
		panic(err)
	}

	brokerUrl := os.Getenv("BROKER_URL")

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
		"contact.messages", // nome da fila
		"",                 // consumer name
		true,               // auto ack
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
		log.Printf("Recebido: %s\n", msg.Body)
	
		// if err := processMessage(msg.Body); err != nil {
		// 	msg.Nack(false, true)
		// 	continue
		// }
	
		// msg.Ack(false)
	}

	// _, err = dispatchEmail(Email{
	// 	To: []string{
	// 		"felipemaejima@gmail.com",
	// 	},
	// 	Subject: "Teste",
	// 	Html:    "<h1>Hello</h1>",
	// })

	// if err != nil {
	// 	log.Println(err)
	// }
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
	}

	sent, err := client.Emails.Send(params)

	if err != nil {
		return "", err
	}

	return sent.Id, nil
}
