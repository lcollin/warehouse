package router

import (
	"fmt"

	"github.com/ghmeier/bloodlines/config"
	g "github.com/ghmeier/bloodlines/gateways"
	h "github.com/ghmeier/bloodlines/handlers"
	coinage "github.com/ghmeier/coinage/gateways"
	tcg "github.com/jakelong95/TownCenter/gateways"
	"github.com/lcollin/warehouse/handlers"
	cg "github.com/yuderekyu/covenant/gateways"

	"gopkg.in/alexcesaro/statsd.v2"
	"gopkg.in/gin-gonic/gin.v1"
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

	s3 := g.NewS3(config.S3)
	tc := tcg.NewTownCenter(config.TownCenter)
	bloodlines := g.NewBloodlines(config.Bloodlines)
	coinage := coinage.NewCoinage(config.Coinage)
	covenant := cg.NewCovenant(config.Covenant)

	ctx := &h.GatewayContext{
		Sql:        sql,
		Stats:      stats,
		TownCenter: tc,
		Coinage:    coinage,
		Covenant:   covenant,
		S3:         s3,
		Bloodlines: bloodlines,
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
	i.router.Use(h.GetCors())

	item := i.router.Group("/api")
	{
		item.Use(i.item.GetJWT())
		item.Use(i.item.Time())
		item.POST("/item", i.item.New)
		item.GET("/item", i.item.ViewAll)
		item.GET("/item/:itemID", i.item.View)
		item.PUT("/item/:itemID", i.item.Update)
		item.DELETE("/item/:itemID", i.item.Delete)
		item.POST("/item/:itemID/photo", i.item.Upload)
		item.GET("/roaster/item/:id", i.item.ViewByRoasterID)
	}

	order := i.router.Group("/api")
	{
		order.Use(i.order.GetJWT())
		order.Use(i.order.Time())
		order.POST("/order", i.order.New)
		order.GET("/order", i.order.ViewAll)
		order.GET("/order/:orderID", i.order.View)
		order.PUT("/order/:orderID", i.order.Update)
		order.DELETE("/order/:orderID", i.order.Delete)
		order.GET("/roaster/order/:id", i.order.ViewByRoasterID)
		order.GET("/user/order/:id", i.order.ViewByUserID)
		order.POST("/label", i.order.GetShippingLabel)
	}

	suborder := i.router.Group("/api")
	{
		suborder.Use(i.suborder.GetJWT())
		suborder.Use(i.suborder.Time())
		suborder.POST("/suborder", i.suborder.New)
		suborder.GET("/suborder", i.suborder.ViewAll)
		suborder.GET("/suborder/:suborderID", i.suborder.View)
		suborder.PUT("/suborder/:suborderID", i.suborder.Update)
		suborder.DELETE("/suborder/:suborderID", i.suborder.Delete)
	}

}

func (s *Inventory) Start(port string) {
	s.router.Run(port)
}
