package helpers

import (
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
func NewEvent(rabbit gateways.RabbitI) Event {
	return &event{
		R: rabbit,
	}
}

func (e *event) Send(sEvent interface{}) error {
	return e.R.Produce(sEvent)
}
