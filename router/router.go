package router

import (
	"fmt"

	"gopkg.in/gin-gonic/gin.v1"

	"github.com/ghmeier/bloodlines/config"
	"github.com/ghmeier/bloodlines/gateways"
	"github.com/ghmeier/bloodlines/handlers"
	"github.com/lcollin/expresso-inventory/handlers"
)

type Inventory struct {
	router       *gin.Engine
	subscription handlers.InventoryIfc
	order 			 handlers.OrderIfc
	// orders handlers.OrdersI
}

func New(config *config.Root) (*Inventory, error) {
	sql, err := gateways.NewSQL(config.SQL)
	if err != nil {
		fmt.Println("ERROR: could not connect to mysql.")
		fmt.Println(err.Error())
		return nil, err
	}

	stats, err := statsd.New(
		statsd.Address(config.Statsd.Host+":"+config.Statsd.Port),
		statsd.Prefix(config.Statsd.Prefix),
	)
	if err != nil {
		fmt.Println(err.Error())
	}

	ctx := &h.GatewayContext{
		Sql:   sql,
		Stats: stats,
	}

	s := &Inventory{
		inventory: handlers.NewInventory(ctx),
	}

	InitRouter(s)
	return s, nil
}

/*InitRouter connects the handlers to endpoints with gin*/
func InitRouter(s *Bloodlines) {
	s.router = gin.Default()

	inventory := s.router.Group("/api/inventory")
	{
		//subscription.POST("", s.inventory.New)
		subscription.GET("", s.inventory.ViewAllInventory)
		subscription.GET("", s.inventory.ViewAllInventory)
		subscription.GET("/:inventoryId", s.inventory.ViewByName)
		subscription.POST("/:inventoryId", s.inventory.Update)
		subscription.POST("/:inventoryId/deactivate", s.inventory.Deactivate)
		subscription.DELETE("/:inventoryId", s.inventory.Cancel)

		/**
		receipt.GET("", b.receipt.ViewAll)
		receipt.POST("/send", b.receipt.Send)
		receipt.GET("/:receiptId", b.receipt.View)
		*/
	}

}

func (s *Inventory) Start(port string) {
	s.router.Run(port)
}
