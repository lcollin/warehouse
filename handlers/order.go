package handlers

import (
	"github.com/ghmeier/bloodlines/handlers"
	"github.com/lcollin/warehouse/helpers"
	"github.com/lcollin/warehouse/models"
	"gopkg.in/alexcesaro/statsd.v2"
	"gopkg.in/gin-gonic/gin.v1"
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
		o.UserError(ctx, "Error: Unable to parse json", err)
		return
	}

	order := models.NewOrder(json.UserID, json.SubscriptionID)
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

func (i *Order) View(ctx *gin.Context) {
	orderID := ctx.Param("orderID")

	order, err := i.Helper.GetByID(orderID)
	if err != nil {
		i.ServerError(ctx, err, orderID)
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

	err = i.Helper.Update(&json, orderID)
	if err != nil {
		i.ServerError(ctx, err, orderID)
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
