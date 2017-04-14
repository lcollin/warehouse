package handlers

import (
	"fmt"

	//"github.com/pborman/uuid"
	shippo "github.com/coldbrewcloud/go-shippo/models"
	"gopkg.in/alexcesaro/statsd.v2"
	"gopkg.in/gin-gonic/gin.v1"

	"github.com/ghmeier/bloodlines/handlers"
	"github.com/lcollin/warehouse/helpers"
)

/*PlanI describes the requests about billing plans that
  can be handled*/
type EventI interface {
	Handle(*gin.Context)
	/*Time tracks the length of execution for each call in the handler*/
	Time() gin.HandlerFunc
}

/*Event implements EventI with coinage helpers*/
type Event struct {
	*handlers.BaseHandler
	Event helpers.Event
}

/*NewEvent initializes and returns a plan with the given gateways*/
func NewEvent(ctx *handlers.GatewayContext) EventI {
	stats := ctx.Stats.Clone(statsd.Prefix("api.event"))
	return &Event{
		Event:       helpers.NewEvent(ctx.Rabbit),
		BaseHandler: &handlers.BaseHandler{Stats: stats},
	}
}

func (e *Event) Handle(ctx *gin.Context) {
	var t shippo.Transaction
	err := ctx.BindJSON(&t)
	if err != nil {
		fmt.Println(err.Error())
		e.ServerError(ctx, err, nil)
		return
	}

	err = e.Event.Send(&t)
	if err != nil {
		e.ServerError(ctx, err, nil)
		return
	}

	e.Success(ctx, nil)
}
