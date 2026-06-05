package main

import (
	"fmt"
	"os"

	"github.com/resend/resend-go/v3"
	"github.com/joho/godotenv"
)

func main() {
	
	err := godotenv.Load()

	if err != nil {
		panic(err)
	}

	apiKey := os.Getenv("RESEND_API_KEY")
	emailDomain := os.Getenv("EMAIL_DOMAIN")

	client := resend.NewClient(apiKey)

	params := &resend.SendEmailRequest{
		From:    emailDomain,
		To:      []string{"felipemaejima@gmail.com"},
		Subject: "Teste",
		Html:    "hELLO!",
	}

	sent, err := client.Emails.Send(params)

	if err != nil {
		panic(err)
	}

	fmt.Println(sent.Id)
}
