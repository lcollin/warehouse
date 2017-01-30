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

func (i *SubOrder) New(ctx *gin.Context) {
	var json models.SubOrder
	err := ctx.BindJSON(&json)
	if err != nil {
		i.UserError(ctx, "Error: Unable to parse json", err)
		return
	}

	suborder := models.NewSubOrder(json.ShopID, json.Name, json.Picture, json.Type, json.InStockBags,
		json.ProviderPrice, json.ConsumerPrice, json.OZInBag)
	err = i.Helper.Insert(suborder)
	if err != nil {
		i.ServerError(ctx, err, json)
		return
	}

	i.Success(ctx, suborder)
}

func (i *SubOrder) ViewAll(ctx *gin.Context) {
	offset, limit := i.GetPaging(ctx)

	suborders, err := i.Helper.GetAll(offset, limit)
	if err != nil {
		i.ServerError(ctx, err, suborders)
		return
	}

	i.Success(ctx, suborders)
}

func (i *SubOrder) View(ctx *gin.Context) {
	suborderId := ctx.Param("suborderId")

	suborder, err := i.Helper.GetByID(suborderId)
	if err != nil {
		i.ServerError(ctx, err, suborderId)
		return
	}

	i.Success(ctx, suborder)
}

func (i *SubOrder) Update(ctx *gin.Context) {
	suborderId := ctx.Param("suborderId")

	var json models.SubOrder
	err := ctx.BindJSON(&json)
	if err != nil {
		i.UserError(ctx, "Error: Unable to parse json", err)
		return
	}

	err = i.Helper.Update(&json, suborderId)
	if err != nil {
		i.ServerError(ctx, err, suborderId)
		return
	}

	i.Success(ctx, json)
}

func (i *SubOrder) Delete(ctx *gin.Context) {
	suborderId := ctx.Param("suborderId")

	err := i.Helper.Delete(suborderId)
	if err != nil {
		i.ServerError(ctx, err, suborderId)
		return
	}

	i.Success(ctx, nil)
}