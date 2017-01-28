package handlers

import (
	"gopkg.in/alexcesaro/statsd.v2"
	"gopkg.in/gin-gonic/gin.v1"

	"github.com/ghmeier/bloodlines/handlers"
	"github.com/lcollin/warehouse/containers"
	"github.com/lcollin/warehouse/helpers"
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
		o.OrderError(ctx, "Error: Unable to parse json", err)
		return
	}

	order := models.NewOrder(json.UserID)
	err = o.Helper.Insert(order)
	if err != nil {
		o.ServerError(ctx, err, json)
		return
	}

	o.Success(ctx, item)
}

//Get order of specific coffee
func (s *Order) GetByOrderID(ctx *gin.Context) {
	id := ctx.Param("id")
	if id == "" {
		ctx.JSON(500, errResponse("id is a required parameter"))
		return
	}

	rows, err := s.sql.Select("SELECT * FROM order WHERE id=?")
	if err != nil {
		ctx.JSON(500, errResponse(err.Error()))
		return
	}
	order, err := containers.FromSql(rows)
	if err != nil {
		ctx.JSON(500, errResponse(err.Error()))
		return
	}

	ctx.JSON(200, gin.H{"data": order})
}

//Get entire order
func (s *Order) ViewAllOrders(ctx *gin.Context) {
	rows, err := s.sql.Select("SELECT * FROM order")
	if err != nil {

		ctx.JSON(500, errResponse(err.Error()))
		return
	}
	orders, err := containers.FromSql(rows)
	if err != nil {
		ctx.JSON(500, errResponse(err.Error()))
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
