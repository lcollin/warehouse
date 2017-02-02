package router

import (
	"fmt"

	"gopkg.in/alexcesaro/statsd.v2"
	"gopkg.in/gin-gonic/gin.v1"
	"github.com/ghmeier/bloodlines/config"
	g "github.com/ghmeier/bloodlines/gateways"
	h "github.com/ghmeier/bloodlines/handlers"
	"github.com/lcollin/warehouse/handlers"
)

type Inventory struct {
	router   *gin.Engine
	item     handlers.ItemIfc
	order    handlers.OrderIfc
	suborder handlers.SubOrderIfc
}

func New(config *config.Root) (*Inventory, error) {
	sql, err := g.NewSQL(config.SQL)
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

	i := &Inventory{
		item:     handlers.NewItem(ctx),
		order:    handlers.NewOrder(ctx),
		suborder: handlers.NewSubOrder(ctx),
	}

	InitRouter(i)
	return i, nil
}

/*InitRouter connects the handlers to endpoints with gin*/
func InitRouter(i *Inventory) {
	i.router = gin.Default()

	item := i.router.Group("/api/item")
	{
		item.POST("", i.item.New)
		item.GET("", i.item.ViewAll)
		item.GET("/:itemId", i.item.View)
		item.PUT("/:itemId", i.item.Update)
		item.DELETE("/:itemId", i.item.Delete)
	}

	order := i.router.Group("/api/order")
	{
		order.POST("", i.order.New)
		order.GET("", i.order.ViewAll)
		order.GET("/:orderId", i.order.View)
		order.PUT("/:orderId", i.order.Update)
		order.DELETE("/:orderId", i.order.Delete)
	}

	suborder := i.router.Group("/api/suborder")
	{
		suborder.POST("", i.suborder.New)
		suborder.GET("", i.suborder.ViewAll)
		suborder.GET("/:suborderId", i.suborder.View)
		suborder.PUT("/:suborderId", i.suborder.Update)
		suborder.DELETE("/:suborderId", i.suborder.Delete)
	}

}

func (s *Inventory) Start(port string) {
	s.router.Run(port)
}
