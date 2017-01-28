package handlers

import (
	"gopkg.in/alexcesaro/statsd.v2"
	"gopkg.in/gin-gonic/gin.v1"

	"github.com/ghmeier/bloodlines/handlers"
	"github.com/lcollin/warehouse/containers"
	"github.com/lcollin/warehouse/helpers"
)

type SubSubOrderIfc interface {
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
		s.SubOrderError(ctx, "Error: Unable to parse json", err)
		return
	}

	suborder := models.NewSubOrder(json.OrderID, json.ItemID)
	err = s.Helper.Insert(suborder)
	if err != nil {
		s.ServerError(ctx, err, json)
		return
	}

	s.Success(ctx, item)
}

//Get suborder of specific coffee
func (s *SubOrder) GetBySubOrderId(ctx *gin.Context) {
	id := ctx.Param("id")
	if id == "" {
		ctx.JSON(500, errResponse("id is a required parameter"))
		return
	}

	rows, err := s.sql.Select("SELECT * FROM suborder WHERE id=?")
	if err != nil {
		ctx.JSON(500, errResponse(err.Error()))
		return
	}
	suborder, err := containers.FromSql(rows)
	if err != nil {
		ctx.JSON(500, errResponse(err.Error()))
		return
	}

	ctx.JSON(200, gin.H{"data": suborder})
}

//Get entire suborder
func (s *SubOrder) ViewAllSubOrders(ctx *gin.Context) {
	rows, err := s.sql.Select("SELECT * FROM suborder")
	if err != nil {

		ctx.JSON(500, errResponse(err.Error()))
		return
	}
	suborders, err := containers.FromSql(rows)
	if err != nil {
		ctx.JSON(500, errResponse(err.Error()))
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
