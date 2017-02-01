package handlers

import (
	"gopkg.in/alexcesaro/statsd.v2"
	"gopkg.in/gin-gonic/gin.v1"
	"github.com/ghmeier/bloodlines/handlers"
	"github.com/lcollin/warehouse/helpers"
	"github.com/lcollin/warehouse/models"
)

type OrderIfc interface {
	New(ctx *gin.Context)
	ViewAll(ctx *gin.Context)
	View(ctx *gin.Context)
	Update(ctx *gin.Context)
	Delete(ctx *gin.Context)
}

type Order struct {
	*handlers.BaseHandler
	Helper helpers.OrderI
}

func NewOrder(ctx *handlers.GatewayContext) OrderIfc {
	stats := ctx.Stats.Clone(statsd.Prefix("api.order"))
	return &Order{
		BaseHandler: &handlers.BaseHandler{Stats: stats},
		Helper:      helpers.NewOrder(ctx.Sql),
	}
}

func (o *Order) New(ctx *gin.Context) {
	var json models.Order
	err := ctx.BindJSON(&json)
	if err != nil {
		i.UserError(ctx, "Error: Unable to parse json", err)
		return
	}

	order := models.NewOrder(json.UserID, json.SubscriptionID, json.RequestDate, json.ShipDate)
	err = i.Helper.Insert(order)
	if err != nil {
		i.ServerError(ctx, err, json)
		return
	}

	i.Success(ctx, order)
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

func (i *Order) View(ctx *gin.Context) {
	orderId := ctx.Param("orderId")

	order, err := i.Helper.GetByID(orderId)
	if err != nil {
		i.ServerError(ctx, err, orderId)
		return
	}

	i.Success(ctx, order)
}

func (i *Order) Update(ctx *gin.Context) {
	orderId := ctx.Param("orderId")

	var json models.Order
	err := ctx.BindJSON(&json)
	if err != nil {
		i.UserError(ctx, "Error: Unable to parse json", err)
		return
	}

	err = i.Helper.Update(&json, orderId)
	if err != nil {
		i.ServerError(ctx, err, orderId)
		return
	}

	i.Success(ctx, json)
}

func (i *Order) Delete(ctx *gin.Context) {
	orderId := ctx.Param("orderId")

	err := i.Helper.Delete(orderId)
	if err != nil {
		i.ServerError(ctx, err, orderId)
		return
	}

	i.Success(ctx, nil)
}

