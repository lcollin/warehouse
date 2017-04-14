package workers

import (
	"encoding/json"
	"fmt"

	//"github.com/pborman/uuid"

	"github.com/ghmeier/bloodlines/gateways"
	"github.com/ghmeier/bloodlines/handlers"
	"github.com/ghmeier/bloodlines/workers"
	"github.com/lcollin/warehouse/helpers"
	//wmodels "github.com/lcollin/warehouse/models"
)

var Events = map[string]bool{
	"track_updated": true,
}

type shippoWorker struct {
	RB   gateways.RabbitI
	Item helpers.ItemI
}

func NewShippoWorker(ctx *handlers.GatewayContext) workers.Worker {
	worker := &eventWorker{
		RB: ctx.Rabbit,
		W:  ctx.Warehouse,
	}

	return &workers.BaseWorker{
		HandleFunc: b.HandleFunc(worker.handle),
		RB:         ctx.Rabbit,
		Item:       helpers.NewItem(ctx.Sql, ctx.S3, ctx.Coinage),
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
