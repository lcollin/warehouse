package router

import (
	"fmt"

	"gopkg.in/alexcesaro/statsd.v2"
	"gopkg.in/gin-gonic/gin.v1"

	"github.com/ghmeier/bloodlines/config"
	"github.com/ghmeier/bloodlines/gateways"
	h "github.com/ghmeier/bloodlines/handlers"
	"github.com/lcollin/warehouse/handlers"
)

type Inventory struct {
	router   *gin.Engine
	item     handlers.ItemIfc
	order    handlers.OrderIfc
	suborder handlers.SubOrderIfc
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
		item:     handlers.NewItem(ctx),
		order:    handlers.NewOrder(ctx),
		suborder: handlers.NewSubOrder(ctx),
	}

	InitRouter(s)
	return s, nil
}

/*InitRouter connects the handlers to endpoints with gin*/
func InitRouter(s *Inventory) {
	s.router = gin.Default()

	item := s.router.Group("/api/item")
	{
		//subscription.POST("", s.inventory.New)
		item.GET("", s.item.ViewAll)
		item.GET("/:itemId", s.item.View)
		item.PUT("/:itemId", s.item.Update)
		item.DELETE("/:itemId", s.item.Delete)

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
