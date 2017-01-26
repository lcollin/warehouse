package handlers

import (
	"gopkg.in/alexcesaro/statsd.v2"
	"gopkg.in/gin-gonic/gin.v1"

	"github.com/ghmeier/bloodlines/handlers"
	"github.com/lcollin/expresso-order/containers"
)

type OrderIfc interface {
	New(ctx *gin.Context)
	GetByOrderId(ctx *gin.Context)
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

func (s *Order) New(ctx *gin.Context) {
	s.Success(ctx, nil)
}

//Get order of specific coffee
func (s *Order) GetByOrderId(ctx *gin.Context) {
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





