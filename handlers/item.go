package handlers

import (
	"github.com/ghmeier/bloodlines/handlers"
	"github.com/lcollin/warehouse/containers"
	"github.com/lcollin/warehouse/helpers"
	"gopkg.in/alexcesaro/statsd.v2"
	"gopkg.in/gin-gonic/gin.v1"
)

type Item interface {
	New(ctx *gin.Context)
	ViewAll(ctx *gin.Context)
	View(ctx *gin.Context)
	Update(ctx *gin.Context)
	Delete(ctx *gin.Context)
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

func (i *Item) New(ctx *gin.Context) {
	var json models.Item
	err := ctx.BindJSON(&json)
	if err != nil {
		i.ItemError(ctx, "Error: Unable to parse json", err)
		return
	}

	item := models.NewItem(json.ShopID, json.Name, json.Picture, json.Type, json.InStockBags,
		                   json.ProviderPrice, json.ConsumerPrice, json.OzInBag)
	err = i.Helper.Insert(item)
	if err != nil {
		i.ServerError(ctx, err, json)
		return
	}

	i.Success(ctx, item)
}

func (i *Item) ViewAll(ctx *gin.Context) {
	offset, limit := i.GetPaging(ctx)

	items, err := i.Helper.GetAll(offset, limit)
	if err != nil {
		i.ServerError(ctx, err, items)
		return
	}

	i.Success(ctx, items)
}

func (i *Item) View(ctx *gin.Context) {
	itemId := ctx.Param("itemId")
	
	item, err := i.Helper.GetByID(itemId)
	if err != nil {
		i.ServerError(ctx, err, itemId)
		return
	}

	i.Success(ctx, item)
}

func (i *Item) Update(ctx *gin.Context) {
	itemId := ctx.Param("itemId")

	var json models.Item
	err := ctx.BindJSON(&json)
	if err != nil {
		i.UserError(ctx, "Error: Unable to parse json", err)
		return
	}

	err = i.Helper.Update(&json, itemId)
	if err != nil {
		i.ServerError(ctx, err, itemId)
		return
	}

	i.Success(ctx, json)
}

func (i *Item) Delete(ctx *gin.Context) {
	itemId := ctx.Param("itemId")

	err := i.Helper.Delete(itemId)
	if err != nil {
		i.ServerError(ctx, err, itemId)
		return
	}

	i.Success(ctx, nil)
}
