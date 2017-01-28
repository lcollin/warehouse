package handlers

import (
	"github.com/ghmeier/bloodlines/handlers"
	"github.com/lcollin/warehouse/containers"
	"gopkg.in/alexcesaro/statsd.v2"
	"gopkg.in/gin-gonic/gin.v1"
)

type Item interface {
	New(ctx *gin.Context)
	GetById(ctx *gin.Context)
	ViewAllItems(ctx *gin.Context)
	CreateOrder(ctx *gin.Context)
	ShipOrder(ctx *gin.Context)
}

type Item struct {
	*handlers.BaseHandler
	Helper helpers.ItemI
}

func NewItem(ctx *handlers.GatewayContext) ItemIfc {
	stats := ctx.Stats.Clone(statsd.Prefix("api.item"))
	return &Item{
		BaseHandler: &handlers.BaseHandler{Stats: stats},
		Helper:      helpers.NewItem(ctx.Sql),
	}
}

func (s *Item) New(ctx *gin.Context) {
	s.Success(ctx, nil)
}

//Get inventory of specific coffee
func (s *Item) GetByName(ctx *gin.Context) {
	id := ctx.Param("name")
	if id == "" {
		ctx.JSON(500, errResponse("name is a required parameter"))
		return
	}

	rows, err := s.sql.Select("SELECT * FROM item WHERE name=?")
	if err != nil {
		ctx.JSON(500, errResponse(err.Error()))
		return
	}
	inventory, err := containers.FromSql(rows)
	if err != nil {
		ctx.JSON(500, errResponse(err.Error()))
		return
	}

	ctx.JSON(200, gin.H{"data": inventory})
}

//Get entire inventory
func (s *Item) ViewAllItem(ctx *gin.Context) {
	rows, err := s.sql.Select("SELECT * FROM inventory")
	if err != nil {

		ctx.JSON(500, errResponse(err.Error()))
		return
	}
	inventories, err := containers.FromSql(rows)
	if err != nil {
		ctx.JSON(500, errResponse(err.Error()))
		return
	}

	s.Success(ctx, inventories)
}

func (s *Item) AddItem(ctx *gin.Context) {
	s.Success(ctx, nil)
}

func (s *Item) RemoveItem(ctx *gin.Context) {
	s.Success(ctx, nil)
}

func (s *Item) ViewImageURL(ctx *gin.Context) {
	s.Success(ctx, nil)
}

func (s *Item) GetPrice(ctx *gin.Context) {
	s.Success(ctx, nil)
}

func (s *Item) GetOzInBag(ctx *gin.Context) {
	s.Success(ctx, nil)
}
