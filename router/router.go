package router

import (
	"fmt"

	"gopkg.in/alexcesaro/statsd.v2"
	"gopkg.in/gin-gonic/gin.v1"

	"github.com/ghmeier/bloodlines/config"
	"github.com/ghmeier/bloodlines/gateways"
	"github.com/ghmeier/bloodlines/handlers"
	"github.com/lcollin/warehouse/handlers"
)

type Inventory struct {
	router   *gin.Engine
	item     handlers.ItemIfc
	order    handlers.OrderIfc
	suborder handlers.SubOrderIfc
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

	item := tc.router.Group("/api/item")
	{
		item.POST("", tc.item.New)
		item.GET("", tc.item.ViewAll)
		item.GET("/:itemId", tc.item.View)
		item.PUT("/:itemId", tc.item.Update)
		item.DELETE("/:itemId", tc.item.Delete)
	}

	order := tc.router.Group("/api/order")
	{
		order.POST("", tc.order.New)
		order.GET("", tc.order.ViewAll)
		order.GET("/:orderId", tc.order.View)
		order.PUT("/:orderId", tc.order.Update)
		order.DELETE("/:orderId", tc.order.Delete)
	}

	suborder := tc.router.Group("/api/suborder")
	{
		suborder.POST("", tc.suborder.New)
		suborder.GET("", tc.suborder.ViewAll)
		suborder.GET("/:suborderId", tc.suborder.View)
		suborder.PUT("/:suborderId", tc.suborder.Update)
		suborder.DELETE("/:suborderId", tc.suborder.Delete)
	}

}

func (s *Inventory) Start(port string) {
	s.router.Run(port)
}
