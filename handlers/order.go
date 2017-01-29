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
	GetByOrderID(ctx *gin.Context)
	ViewAllOrders(ctx *gin.Context)
	CreateOrder(ctx *gin.Context)
	ShipOrder(ctx *gin.Context)
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

	order := models.NewOrder(json.UserID)
	err = o.Helper.Insert(order)
	if err != nil {
		o.ServerError(ctx, err, json)
		return
	}

	o.Success(ctx, order)
}

//Get order of specific coffee
func (s *Order) GetByOrderID(ctx *gin.Context) {
	id := ctx.Param("id")

	order, err := s.Helper.GetByID(id)
	if err != nil {
		s.ServerError(ctx, err, id)
		return
	}

	s.Success(ctx, order)
}

//Get entire order
func (s *Order) ViewAllOrders(ctx *gin.Context) {
	offset, limit := s.GetPaging(ctx)

	orders, err := s.Helper.GetAll(offset, limit)
	if err != nil {
		s.ServerError(ctx, err, nil)
		return
	}

	s.Success(ctx, orders)
}

// Subscriptions will create orders automatically
// Params to include: item id, store id, etc
func (s *Order) CreateOrder(ctx *gin.Context) {
	s.Success(ctx, nil)
}

// Provider will notify that an order has been shipped
// Params to include: shipping tracking number
func (s *Order) ShipOrder(ctx *gin.Context) {
	s.Success(ctx, nil)
}
