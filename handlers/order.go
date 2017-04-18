package handlers

import (
	"fmt"

	"github.com/ghmeier/bloodlines/handlers"
	"github.com/lcollin/warehouse/helpers"
	"github.com/lcollin/warehouse/models"

	"github.com/pborman/uuid"
	"gopkg.in/alexcesaro/statsd.v2"
	"gopkg.in/gin-gonic/gin.v1"
)

type OrderIfc interface {
	New(ctx *gin.Context)
	ViewAll(ctx *gin.Context)
	View(ctx *gin.Context)
	Update(ctx *gin.Context)
	Delete(ctx *gin.Context)
	GetShippingLabel(ctx *gin.Context)
	ViewByUserID(ctx *gin.Context)
	ViewByRoasterID(ctx *gin.Context)
	Time() gin.HandlerFunc
	GetJWT() gin.HandlerFunc
}

type Order struct {
	*handlers.BaseHandler
	Helper helpers.OrderI
	Item   helpers.ItemI
}

func NewOrder(ctx *handlers.GatewayContext) OrderIfc {
	stats := ctx.Stats.Clone(statsd.Prefix("api.order"))
	return &Order{
		BaseHandler: &handlers.BaseHandler{Stats: stats},
		Helper:      helpers.NewOrder(ctx.Sql, ctx.TownCenter, ctx.Bloodlines),
		Item:        helpers.NewItem(ctx.Sql, ctx.S3, ctx.Coinage, ctx.Covenant, ctx.Bloodlines, ctx.TownCenter),
	}
}

func (o *Order) New(ctx *gin.Context) {
	var json models.Order
	err := ctx.BindJSON(&json)
	if err != nil {
		o.UserError(ctx, "Error: Unable to parse json", err)
		return
	}
	item, err := o.Item.GetBySubscription(json.SubscriptionID)
	if err != nil {
		o.ServerError(ctx, err, json)
		return
	}
	if item == nil {
		o.NotFoundError(ctx, "Error: no item for subscription")
		return
	}

	err = o.Item.RemoveStock(item, int(json.Quantity))
	if err != nil {
		fmt.Println(err.Error())
	}

	order := models.NewOrder(json.UserID, json.SubscriptionID, json.Quantity)
	err = o.Helper.Insert(order)
	if err != nil {
		o.ServerError(ctx, err, json)
		return
	}

	o.Success(ctx, order)
}

func (i *Order) ViewAll(ctx *gin.Context) {
	offset, limit := i.GetPaging(ctx)

	orders, err := i.Helper.GetAll(offset, limit)
	if err != nil {
		i.ServerError(ctx, err, orders)
		return
	}

	i.Success(ctx, orders)
}

func (i *Order) ViewByUserID(ctx *gin.Context) {
	offset, limit := i.GetPaging(ctx)
	id := ctx.Param("id")

	orders, err := i.Helper.GetByUserID(uuid.Parse(id), offset, limit)
	if err != nil {
		i.ServerError(ctx, err, orders)
		return
	}

	i.Success(ctx, orders)
}

func (i *Order) ViewByRoasterID(ctx *gin.Context) {
	offset, limit := i.GetPaging(ctx)
	id := ctx.Param("id")

	orders, err := i.Helper.GetByRoasterID(uuid.Parse(id), offset, limit)
	if err != nil {
		i.ServerError(ctx, err, orders)
		return
	}

	i.Success(ctx, orders)
}

func (i *Order) View(ctx *gin.Context) {
	orderID := ctx.Param("orderID")

	order, err := i.Helper.GetByID(orderID)
	if err != nil {
		i.ServerError(ctx, err, nil)
		return
	}

	i.Success(ctx, order)
}

func (i *Order) GetShippingLabel(ctx *gin.Context) {
	var json models.ShipmentRequest
	err := ctx.BindJSON(&json)
	if err != nil {
		i.UserError(ctx, "Error: Unable to parse json", err)
		return
	}
	order, err := i.Helper.GetShippingLabel(&json)
	if err != nil {
		i.ServerError(ctx, err, nil)
		return
	}

	i.Success(ctx, order)
}

func (i *Order) Update(ctx *gin.Context) {
	orderID := ctx.Param("orderID")

	var json models.Order
	err := ctx.BindJSON(&json)
	if err != nil {
		i.UserError(ctx, "Error: Unable to parse json", err)
		return
	}

	json.ID = uuid.Parse(orderID)
	err = i.Helper.Update(&json)
	if err != nil {
		i.ServerError(ctx, err, json)
		return
	}

	i.Success(ctx, json)
}

func (i *Order) Delete(ctx *gin.Context) {
	orderID := ctx.Param("orderID")

	err := i.Helper.Delete(orderID)
	if err != nil {
		i.ServerError(ctx, err, orderID)
		return
	}

	i.Success(ctx, nil)
}
