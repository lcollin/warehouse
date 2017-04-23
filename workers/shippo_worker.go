package workers

import (
	"encoding/json"
	"fmt"

	shippo "github.com/coldbrewcloud/go-shippo/models"

	"github.com/ghmeier/bloodlines/gateways"
	"github.com/ghmeier/bloodlines/handlers"
	"github.com/ghmeier/bloodlines/models"
	"github.com/ghmeier/bloodlines/workers"
	tc "github.com/jakelong95/TownCenter/gateways"
	"github.com/lcollin/warehouse/helpers"
	wmodels "github.com/lcollin/warehouse/models"
)

var Events = map[string]bool{
	"track_updated": true,
}

type shippoWorker struct {
	RB         gateways.RabbitI
	Item       helpers.ItemI
	Order      helpers.OrderI
	Bloodlines gateways.Bloodlines
	TownCenter tc.TownCenterI
}

func NewShippoWorker(ctx *handlers.GatewayContext) workers.Worker {
	worker := &shippoWorker{
		RB:         ctx.Rabbit,
		Item:       helpers.NewItem(ctx.Sql, ctx.S3, ctx.Coinage, ctx.Covenant, ctx.Bloodlines, ctx.TownCenter),
		Order:      helpers.NewOrder(ctx.Sql, ctx.TownCenter, ctx.Bloodlines, ctx.Shippo),
		TownCenter: ctx.TownCenter,
		Bloodlines: ctx.Bloodlines,
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

	order, err := s.Order.GetByTransactionID(t.ObjectID)
	if err != nil {
		fmt.Println(err)
		return
	}
	if order == nil {
		fmt.Println("ERROR: no order for id " + t.ObjectID)
		return
	}

	err = order.SetStatus(t.TrackingStatus.Status)
	if err != nil {
		fmt.Println(err)
		return
	}

	err = s.Order.Update(order)
	if err != nil {
		fmt.Println(err)
		return
	}

	if order.Status != wmodels.SHIPPED && order.Status != wmodels.TRANSIT {
		fmt.Println(s)
		return
	}

	item, err := s.Item.GetBySubscription(order.SubscriptionID)
	if err != nil {
		return
	}

	user, err := s.TownCenter.GetUser(order.UserID)
	if err != nil {
		return
	}
	if user == nil {
		return
	}

	s.Bloodlines.ActivateTrigger("order_status_update", &models.Receipt{
		UserID: user.ID,
		Values: map[string]string{
			"first_name":   user.FirstName,
			"last_name":    user.LastName,
			"item_name":    item.Name,
			"tracking_url": order.TrackingURL,
		},
	})

}
