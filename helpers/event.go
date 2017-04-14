package helpers

import (
	shippo "github.com/coldbrewcloud/go-shippo/models"

	"github.com/ghmeier/bloodlines/gateways"
)

/*Event helps with manipulating roaster properties*/
type Event interface {
	Send(*shippo.Transaction) error
}

type event struct {
	R gateways.RabbitI
}

/*NewEvent initializes and returns a roaster with the given gateways*/
func NewEvent(rabbit gateways.RabbitI) Event {
	return &event{
		R: rabbit,
	}
}

func (e *event) Send(t *shippo.Transaction) error {
	return e.R.Produce(t)
}
