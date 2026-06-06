package resend

import (
	"os"

	"github.com/felipemaejima/message-dispatcher/internal/domain"
	"github.com/resend/resend-go/v3"
)

type Notifier struct {
	client *resend.Client
	from   string
}

func NewNotifier() *Notifier {
	return &Notifier{
		client: resend.NewClient(
			os.Getenv("RESEND_API_KEY"),
		),
		from: os.Getenv("EMAIL_DOMAIN"),
	}
}

func (n *Notifier) Send(
	message domain.Message,
) error {

	params := &resend.SendEmailRequest{
		From:    n.from,
		To:      message.To,
		Subject: message.Subject,
		Html:    message.Body,
	}

	_, err := n.client.Emails.Send(params)

	return err
}
