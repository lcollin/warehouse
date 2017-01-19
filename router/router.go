package router

import (
	"fmt"

	"gopkg.in/gin-gonic/gin.v1"

	"github.com/lcollin/expresso-inventory/handlers"
	"github.com/lcollin/expresso-inventory/gateways"
)

type Inventory struct {
	router *gin.Engine
	subscription handlers.SubscriptionIfc
	// orders handlers.OrdersI
}

func New() (*Inventory, error) {
	sql, err := gateways.NewSql()

	if err != nil {
		fmt.Println("ERROR: could not connect to mysql.")
		fmt.Println(err.Error())
		return nil, err
	}

	s := &Inventory{
		inventory : handlers.NewInventory(sql),
	}

	InitRouter(s)
	return s, nil
}

/*InitRouter connects the handlers to endpoints with gin*/
func InitRouter(s *Bloodlines) {
	s.router = gin.Default()


	inventory := s.router.Group("/api/inventory")
	{
		subscription.POST("", s.inventory.New)
		subscription.GET("", s.inventory.ViewAllInventory)
		subscription.GET("/:inventoryId", s.inventory.ViewByName)
		subscription.POST("/:inventoryId", s.inventory.Update)
		subscription.POST("/:inventoryId/deactivate", s.inventory.Deactivate)
		subscription.DELETE("/:inventoryId", s.inventory.Cancel)
	}
	
}

func (s *Inventory) Start(port string) {
	s.router.Run(port)
}