package ports

import "github.com/felipemaejima/message-dispatcher/internal/domain"

type Notifier interface {
	Send(message domain.Message) error
}
