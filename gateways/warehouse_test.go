package gateways

import (
	//"encoding/json"
	//"fmt"
	"testing"

	"github.com/jarcoal/httpmock"
	//"github.com/pborman/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"

	"github.com/ghmeier/bloodlines/config"
	"github.com/ghmeier/bloodlines/gateways"
	"github.com/lcollin/warehouse/models"
)

type WarehouseSuite struct {
	suite.Suite
	warehouse Warehouse
	url       string
}

func (b *WarehouseSuite) SetupSuite() {
	httpmock.Activate()
	b.warehouse = NewWarehouse(config.Warehouse{
		Host: "warehouse",
		Port: "8080",
	})
	b.url = "http://warehouse:8080/api/"
}

func (b *WarehouseSuite) BeforeTest() {
	httpmock.Reset()
}

func (b *WarehouseSuite) AfterTest() {

}

func (b *WarehouseSuite) TearDownSuite() {
	httpmock.DeactivateAndReset()
}

func TestRunBloodlinesSuite(t *testing.T) {
	s := new(WarehouseSuite)
	suite.Run(t, s)
}

func (b *WarehouseSuite) TestGetAllContentSuccess() {
	assert := assert.New(b.T())

	data := b.SuccessResponse()
	data.Data = make([]*models.Item, 0)
	res, _ := httpmock.NewJsonResponder(200, data)

	httpmock.RegisterResponder("GET", b.url+"item?offset=0&limit=20", res)

	items, err := b.warehouse.GetAllItems(0, 20)

	assert.NoError(err)
	assert.NotNil(items)
}

func (b *WarehouseSuite) EmptyResponse() *gateways.ServiceResponse {
	return &gateways.ServiceResponse{}
}

func (b *WarehouseSuite) SuccessResponse() *gateways.ServiceResponse {
	r := b.EmptyResponse()
	r.Success = true
	return r
}

func (b *WarehouseSuite) ErrorResponse(msg string) *gateways.ServiceResponse {
	r := b.EmptyResponse()
	r.Success = false
	r.Msg = msg
	return r
}
