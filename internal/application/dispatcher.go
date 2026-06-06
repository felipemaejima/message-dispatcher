package application

import (
	"errors"

	"github.com/felipemaejima/message-dispatcher/internal/domain"
	"github.com/felipemaejima/message-dispatcher/internal/ports"
)

type Dispatcher struct {
	notifiers map[string]ports.Notifier
}

func NewDispatcher(
	notifiers map[string]ports.Notifier,
) *Dispatcher {

	return &Dispatcher{
		notifiers: notifiers,
	}

}

func (d *Dispatcher) Dispatch(
	notification domain.Notification,
) error {

	notifier, exists := d.notifiers[notification.Channel]

	if !exists {
		return errors.New(
			"unsupported channel",
		)
	}

	return notifier.Send(notification.Message)
}
