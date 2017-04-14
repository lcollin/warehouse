package workers

import (
	"encoding/json"
	"fmt"

	//"github.com/pborman/uuid"
	shippo "github.com/coldbrewcloud/go-shippo/models"

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
	worker := &shippoWorker{
		RB:   ctx.Rabbit,
		Item: helpers.NewItem(ctx.Sql, ctx.S3, ctx.Coinage),
	}

	return &workers.BaseWorker{
		HandleFunc: workers.HandleFunc(worker.handle),
		RB:         ctx.Rabbit,
	}
}

func (e *shippoWorker) handle(body []byte) {
	var t shippo.Transaction
	err := json.Unmarshal(body, &t)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	e.updateOrder(&t)
}

func (s *shippoWorker) updateOrder(t *shippo.Transaction) {
	fmt.Println(t)
}
