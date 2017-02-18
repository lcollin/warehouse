package router

import (
	"testing"

	mockg "github.com/ghmeier/bloodlines/_mocks/gateways"
	"github.com/ghmeier/bloodlines/config"
	h "github.com/ghmeier/bloodlines/handlers"
	mocks "github.com/lcollin/warehouse/_mocks"
	"github.com/lcollin/warehouse/handlers"

	"github.com/stretchr/testify/assert"
	"gopkg.in/alexcesaro/statsd.v2"
)

func TestNewSuccess(t *testing.T) {
	assert := assert.New(t)

	r, err := New(&config.Root{SQL: config.MySQL{}})

	assert.NoError(err)
	assert.NotNil(r)
}

func getMockTownCenter() *TownCenter {
	sql := new(mockg.SQL)
	stats, _ := statsd.New()
	ctx := &h.GatewayContext{
		Sql:   sql,
		Stats: stats,
	}

	return &TownCenter{
		item:     handlers.NewItem(ctx),
		order:    handlers.NewOrder(ctx),
		suborder: handlers.NewSubOrder(ctx),
	}
}

func mockItem() (*TownCenter, *mocks.ItemI) {
	t := getMockTownCenter()
	mock := new(mocks.ItemI)
	t.item = &handlers.Item{Helper: mock, BaseHandler: &h.BaseHandler{Stats: nil}}
	InitRouter(t)

	return t, mock
}

func mockOrder() (*TownCenter, *mocks.OrderI) {
	t := getMockTownCenter()
	mock := new(mocks.OrderI)
	t.order = &handlers.Order{Helper: mock, BaseHandler: &h.BaseHandler{Stats: nil}}
	InitRouter(t)

	return t, mock
}

func mockSubOrder() (*TownCenter, *mocks.SubOrderI) {
	t := getMockTownCenter()
	mock := new(mocks.SubOrderI)
	t.suborder = &handlers.SubOrder{Helper: mock, BaseHandler: &h.BaseHandler{Stats: nil}}
	InitRouter(t)

	return t, mock
}
