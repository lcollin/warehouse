package workers

import (
	"encoding/json"
	"fmt"

	"github.com/pborman/uuid"
	"github.com/stripe/stripe-go"

	"github.com/ghmeier/bloodlines/gateways"
	"github.com/ghmeier/bloodlines/handlers"
	b "github.com/ghmeier/bloodlines/workers"
	warehouse "github.com/lcollin/warehouse/gateways"
	wmodels "github.com/lcollin/warehouse/models"
)

var Events = map[string]bool{
	"track_updated": true,
}

type shippoWorker struct {
	RB gateways.RabbitI
	W  waerhouse.Warehouse
}

func NewShippoWorker(ctx *handlers.GatewayContext) b.Worker {
	worker := &eventWorker{
		RB: ctx.Rabbit,
		W:  ctx.Warehouse,
	}

	return &b.BaseWorker{
		HandleFunc: b.HandleFunc(worker.handle),
		RB:         ctx.Rabbit,
	}
}

func (e *shippoWorker) handle(body []byte) {
	var event stripe.Event
	err := json.Unmarshal(body, &event)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	switch event.Type {
	case "invoice.created":
		e.invoiceCreate(&event)

	}
}
