package resend

import (
	"testing"

	"github.com/felipemaejima/message-dispatcher/internal/domain"
	"github.com/joho/godotenv"
)

func TestSendEmail(t *testing.T) {
	err := godotenv.Load("../../../.env")

	if err != nil {
		t.Fatalf("erro carregando env: %v", err)
	}

	notifier := NewNotifier()

	message := domain.Message{
		To: []string{
			"felipemaejima@gmail.com",
		},
		Subject: "Integration Test",
		Body:    "<h1>Teste de integração</h1>",
	}

	err = notifier.Send(message)

	if err != nil {
		t.Fatalf("erro ao enviar email: %v", err)
	}
}
