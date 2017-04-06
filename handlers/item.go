package handlers

import (
	"github.com/ghmeier/bloodlines/handlers"
	"github.com/lcollin/warehouse/helpers"
	"github.com/lcollin/warehouse/models"
	"gopkg.in/alexcesaro/statsd.v2"
	"gopkg.in/gin-gonic/gin.v1"
)

type ItemIfc interface {
	New(ctx *gin.Context)
	ViewAll(ctx *gin.Context)
	ViewByRoasterID(ctx *gin.Context)
	View(ctx *gin.Context)
	Update(ctx *gin.Context)
	Delete(ctx *gin.Context)
	Upload(ctx *gin.Context)
	Time() gin.HandlerFunc
	GetJWT() gin.HandlerFunc
}

type Item struct {
	*handlers.BaseHandler
	Helper helpers.ItemI
}

func NewItem(ctx *handlers.GatewayContext) ItemIfc {
	stats := ctx.Stats.Clone(statsd.Prefix("api.item"))
	return &Item{
		BaseHandler: &handlers.BaseHandler{Stats: stats},
		Helper:      helpers.NewItem(ctx.Sql, ctx.S3),
	}
}

func (i *Item) New(ctx *gin.Context) {
	var json models.Item
	err := ctx.BindJSON(&json)
	if err != nil {
		i.UserError(ctx, "Error: Unable to parse json", err)
		return
	}

	item := models.NewItem(json.RoasterID, json.Name, json.PictureURL, json.CoffeeType, json.InStockBags,
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
	search := models.ItemSearch(ctx)

	items, err := i.Helper.GetAll(offset, limit, search)
	if err != nil {
		i.ServerError(ctx, err, items)
		return
	}

	i.Success(ctx, items)
}

func (i *Item) ViewByRoasterID(ctx *gin.Context) {
	roasterID := ctx.Param("roasterID")

	items, err := i.Helper.GetByRoasterID(roasterID)
	if err != nil {
		i.ServerError(ctx, err, items)
		return
	}

	if items == nil {
		i.NotFoundError(ctx, "ERROR: No items available")
		return
	}

	i.Success(ctx, items)
}

func (i *Item) View(ctx *gin.Context) {
	itemID := ctx.Param("itemID")

	item, err := i.Helper.GetByID(itemID)
	if err != nil {
		i.ServerError(ctx, err, itemID)
		return
	}

	i.Success(ctx, item)
}

func (i *Item) Update(ctx *gin.Context) {
	itemID := ctx.Param("itemID")

	var json models.Item
	err := ctx.BindJSON(&json)
	if err != nil {
		i.UserError(ctx, "Error: Unable to parse json", err)
		return
	}

	err = i.Helper.Update(&json, itemID)
	if err != nil {
		i.ServerError(ctx, err, itemID)
		return
	}

	i.Success(ctx, json)
}

func (i *Item) Delete(ctx *gin.Context) {
	itemID := ctx.Param("itemID")

	err := i.Helper.Delete(itemID)
	if err != nil {
		i.ServerError(ctx, err, itemID)
		return
	}

	i.Success(ctx, nil)
}

func (i *Item) Upload(ctx *gin.Context) {
	id := ctx.Param("itemID")

	file, headers, err := ctx.Request.FormFile("photo")
	if err != nil {
		i.UserError(ctx, "ERROR: unable to find body", nil)
	}
	defer file.Close()

	err = i.Helper.Upload(id, headers.Filename, file)
	if err != nil {
		i.ServerError(ctx, err, id)
		return
	}

	i.Success(ctx, nil)
}
