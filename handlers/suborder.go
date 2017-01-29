package handlers

import (
	"gopkg.in/alexcesaro/statsd.v2"
	"gopkg.in/gin-gonic/gin.v1"

	"github.com/ghmeier/bloodlines/handlers"
	"github.com/lcollin/warehouse/helpers"
	"github.com/lcollin/warehouse/models"
)

type SubOrderIfc interface {
	New(ctx *gin.Context)
	GetBySubOrderId(ctx *gin.Context)
	ViewAllSubOrders(ctx *gin.Context)
	CreateSubOrder(ctx *gin.Context)
	ShipSubOrder(ctx *gin.Context)
}

type SubOrder struct {
	*handlers.BaseHandler
	Helper helpers.SubOrderI
}

func NewSubOrder(ctx *handlers.GatewayContext) SubOrderIfc {
	stats := ctx.Stats.Clone(statsd.Prefix("api.suborder"))
	return &SubOrder{
		BaseHandler: &handlers.BaseHandler{Stats: stats},
		Helper:      helpers.NewSubOrder(ctx.Sql),
	}
}

func (s *SubOrder) New(ctx *gin.Context) {
	var json models.SubOrder
	err := ctx.BindJSON(&json)
	if err != nil {
		s.UserError(ctx, "Error: Unable to parse json", err)
		return
	}

	suborder := models.NewSubOrder(json.OrderID, json.ItemID)
	err = s.Helper.Insert(suborder)
	if err != nil {
		s.ServerError(ctx, err, json)
		return
	}

	s.Success(ctx, suborder)
}

//Get suborder of specific coffee
func (s *SubOrder) GetBySubOrderId(ctx *gin.Context) {
	id := ctx.Param("id")

	suborder, err := s.Helper.GetByID(id)
	if err != nil {
		s.ServerError(ctx, err, id)
		return
	}

	s.Success(ctx, suborder)
}

//Get entire suborder
func (s *SubOrder) ViewAllSubOrders(ctx *gin.Context) {
	offset, limit := s.GetPaging(ctx)

	suborders, err := s.Helper.GetAll(offset, limit)
	if err != nil {
		s.ServerError(ctx, err, nil)
		return
	}

	s.Success(ctx, suborders)
}

// Subscriptions will create suborders automatically
// Params to include: item id, store id, etc
func (s *SubOrder) CreateSubOrder(ctx *gin.Context) {
	s.Success(ctx, nil)
}

// Provider will notify that an suborder has been shipped
// Params to include: shipping tracking number
func (s *SubOrder) ShipSubOrder(ctx *gin.Context) {
	s.Success(ctx, nil)
}
