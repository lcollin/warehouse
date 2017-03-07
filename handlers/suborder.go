package handlers

import (
	"github.com/ghmeier/bloodlines/handlers"
	"github.com/lcollin/warehouse/helpers"
	"github.com/lcollin/warehouse/models"
	"gopkg.in/alexcesaro/statsd.v2"
	"gopkg.in/gin-gonic/gin.v1"
)

type SubOrderIfc interface {
	New(ctx *gin.Context)
	ViewAll(ctx *gin.Context)
	View(ctx *gin.Context)
	Update(ctx *gin.Context)
	Delete(ctx *gin.Context)
	Time() gin.HandlerFunc
	GetJWT() gin.HandlerFunc
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

	suborder := models.NewSubOrder(json.OrderID, json.ItemID, json.Quantity)
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
	suborderID := ctx.Param("suborderID")

	suborder, err := i.Helper.GetByID(suborderID)
	if err != nil {
		i.ServerError(ctx, err, suborderID)
		return
	}

	i.Success(ctx, suborder)
}

func (i *SubOrder) Update(ctx *gin.Context) {
	suborderID := ctx.Param("suborderID")

	var json models.SubOrder
	err := ctx.BindJSON(&json)
	if err != nil {
		i.UserError(ctx, "Error: Unable to parse json", err)
		return
	}

	err = i.Helper.Update(&json, suborderID)
	if err != nil {
		i.ServerError(ctx, err, suborderID)
		return
	}

	i.Success(ctx, json)
}

func (i *SubOrder) Delete(ctx *gin.Context) {
	suborderID := ctx.Param("suborderID")

	err := i.Helper.Delete(suborderID)
	if err != nil {
		i.ServerError(ctx, err, suborderID)
		return
	}

	i.Success(ctx, nil)
}
