package handlers

import (
	"gopkg.in/alexcesaro/statsd.v2"
	"gopkg.in/gin-gonic/gin.v1"

	"github.com/ghmeier/bloodlines/handlers"
	"github.com/lcollin/expresso-inventory/containers"
)

type InventoryIfc interface {
	New(ctx *gin.Context)
	ViewAllInventory(ctx *gin.Context)
	GetByName(ctx *gin.Context)
	// Update(ctx *gin.Context)
	// Deactivate(ctx *gin.Context)
	// Cancel(ctx *gin.Context)
}

type Inventory struct {
	*handlers.BaseHandler
	/*Helper ....*/
}

func NewInventory(ctx *handlers.GatewayContext) InventoryIfc {
	stats := ctx.Stats.Clone(statsd.Prefix("api.inventory"))
	return &Inventory{
		/*Helper: ..... */
		BaseHandler: &handlers.BaseHandler{Stats: stats},
	}
}

func (s *Inventory) New(ctx *gin.Context) {
	// var json containers.Inventory
	// err := ctx.BindJSON(&json)

	// if err != nil {
	// 	ctx.JSON(400, errResponse("Invalid Inventory Object"))
	// 	fmt.Printf("%s", err.Error())
	// 	return
	// }

	// //todo add in reference to orders params?
	// err = s.sql.Modify(
	// 	"INSERT INTO subscription VALUE(?, ?, ?, ?, ?, ?, ?, ?)",
	// 	uuid.New(),
	// 	json.OrderId,
	// 	json.Type,
	// 	json.UserId,
	// 	json.Status,
	// 	json.CreatedAt,
	// 	json.StartAt,
	// 	json.TotalPrice)

	// if err != nil {
	// 	ctx.JSON(500, &gin.H{"error": err, "message": err.Error()})
	// 	return
	// }
	s.Success(ctx, nil)
}

//Get inventory of specific coffee
func (s *Inventory) GetByName(ctx *gin.Context) {
	id := ctx.Param("name")
	if id == "" {
		ctx.JSON(500, errResponse("name is a required parameter"))
		return
	}

	rows, err := s.sql.Select("SELECT * FROM inventory WHERE name=?")
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
func (s *Inventory) ViewAllInventory(ctx *gin.Context) {
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

func (s *Inventory) AddInventory(ctx *gin.Context) {
	s.Success(ctx, nil)
}

func (s *Inventory) RemoveInventory(ctx *gin.Context) {
	s.Success(ctx, nil)
}

func (s *Inventory) ViewImageURL(ctx *gin.Context) {
	s.Success(ctx, nil)
}

func (s *Inventory) GetPrice(ctx *gin.Context) {
	s.Success(ctx, nil)
}

func (s *Inventory) GetOzInBag(ctx *gin.Context) {
	s.Success(ctx, nil)
}
