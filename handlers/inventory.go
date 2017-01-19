package handlers

import(
	"fmt"

	"gopkg.in/gin-gonic/gin.v1"
	"github.com/pborman/uuid"

	"github.com/lcollin/expresso-inventory/containers"
	"github.com/lcollin/expresso-inventory/gateways"
)

type InventoryIfc interface {
	New(ctx *gin.Context)
	ViewAllInventory(ctx *gin.Context)
	ViewByName(ctx *gin.Context)
	// Update(ctx *gin.Context)
	// Deactivate(ctx *gin.Context)
	// Cancel(ctx *gin.Context)
}

type Inventory struct {
	sql *gateways.Sql
}

func NewInventory(sql *gateways.Sql) InventoryIfc {
	return &Inventory{sql: sql}
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
	ctx.JSON(200, empty())
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
	subscription, err := containers.FromSql(rows)
	if err != nil {
		ctx.JSON(500, errResponse(err.Error()))
		return
	}

	ctx.JSON(200, gin.H{"data": inventory})
}

func (s *Subscription) AddInventory(ctx *gin.Context) {
	ctx.JSON(200, empty())
}

func (s *Subscription) RemoveInventory(ctx *gin.Context) {
	ctx.JSON(200, empty())
}

func (s *Subscription) ViewImageURL(ctx *gin.Context) {
	ctx.JSON(200, empty())
}

func (s *Subscription) GetPrice(ctx *gin.Context) {
	ctx.JSON(200, empty())
}

func (s *Subscription) GetOzInBag(ctx *gin.Context) {
	ctx.JSON(200, empty())
}