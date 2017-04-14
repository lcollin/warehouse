package helpers

import (
	"github.com/stripe/stripe-go"

	"github.com/ghmeier/bloodlines/gateways"
)

/*Event helps with manipulating roaster properties*/
type Event interface {
	Send(interface{}) error
}

type event struct {
	R gateways.RabbitI
}

/*NewEvent initializes and returns a roaster with the given gateways*/
func NewEvent(rabbit g.RabbitI) Event {
	return &event{
		R: rabbit,
	}
}

func (e *event) Send(e interface{}) error {
	return e.R.Produce(e)
}
