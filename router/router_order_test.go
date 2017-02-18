package router

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/lcollin/warehouse/models"

	"github.com/pborman/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gopkg.in/gin-gonic/gin.v1"
)

func TestOrderViewSuccess(t *testing.T) {
	assert := assert.New(t)

	gin.SetMode(gin.TestMode)

	id := uuid.NewUUID()
	tc, orderMock := mockOrder()
	orderMock.On("GetByID", id.String()).Return(&models.Order{}, nil)

	recorder := httptest.NewRecorder()
	request, _ := http.NewRequest("GET", "/api/order/"+id.String(), nil)
	tc.router.ServeHTTP(recorder, request)

	assert.Equal(200, recorder.Code)
}

func TestOrderViewFail(t *testing.T) {
	assert := assert.New(t)

	gin.SetMode(gin.TestMode)

	id := uuid.NewUUID()
	tc, orderMock := mockOrder()
	orderMock.On("GetByID", id.String()).Return(nil, fmt.Errorf("This is an error"))

	recorder := httptest.NewRecorder()
	request, _ := http.NewRequest("GET", "/api/order/"+id.String(), nil)
	tc.router.ServeHTTP(recorder, request)

	assert.Equal(500, recorder.Code)
}

func TestOrderViewAllSuccess(t *testing.T) {
	assert := assert.New(t)

	gin.SetMode(gin.TestMode)

	tc, orderMock := mockOrder()
	orderMock.On("GetAll", 0, 20).Return(make([]*models.Order, 0), nil)

	recorder := httptest.NewRecorder()
	request, _ := http.NewRequest("GET", "/api/order", nil)
	tc.router.ServeHTTP(recorder, request)

	assert.Equal(200, recorder.Code)
}

func TestOrderViewAllFail(t *testing.T) {
	assert := assert.New(t)

	gin.SetMode(gin.TestMode)

	tc, orderMock := mockOrder()
	orderMock.On("GetAll", 0, 20).Return(make([]*models.Order, 0), fmt.Errorf("This is an error"))

	recorder := httptest.NewRecorder()
	request, _ := http.NewRequest("GET", "/api/order", nil)
	tc.router.ServeHTTP(recorder, request)

	assert.Equal(500, recorder.Code)
}

func TestOrderViewAllParams(t *testing.T) {
	assert := assert.New(t)

	gin.SetMode(gin.TestMode)

	tc, orderMock := mockOrder()
	orderMock.On("GetAll", 20, 40).Return(make([]*models.Order, 0), nil)

	recorder := httptest.NewRecorder()
	request, _ := http.NewRequest("GET", "/api/order?offset=20&limit=40", nil)
	tc.router.ServeHTTP(recorder, request)

	assert.Equal(200, recorder.Code)
}

func TestOrderNewSuccess(t *testing.T) {
	assert := assert.New(t)

	gin.SetMode(gin.TestMode)

	tc, orderMock := mockOrder()
	orderMock.On("Insert", mock.AnythingOfType("*models.Order")).Return(nil)

	order := getOrderString(models.NewOrder("", "", "", ""))
	recorder := httptest.NewRecorder()
	request, _ := http.NewRequest("POST", "/api/order", order)
	tc.router.ServeHTTP(recorder, request)

	assert.Equal(200, recorder.Code)
}

func TestOrderNewFail(t *testing.T) {
	assert := assert.New(t)

	gin.SetMode(gin.TestMode)

	tc, orderMock := mockOrder()
	orderMock.On("Insert", mock.AnythingOfType("*models.Order")).Return(fmt.Errorf("This is an error"))

	order := getOrderString(models.NewOrder("", "", "", ""))
	recorder := httptest.NewRecorder()
	request, _ := http.NewRequest("POST", "/api/order", order)
	tc.router.ServeHTTP(recorder, request)

	assert.Equal(500, recorder.Code)
}

func TestOrderNewInvalid(t *testing.T) {
	assert := assert.New(t)

	gin.SetMode(gin.TestMode)

	tc, orderMock := mockOrder()
	orderMock.On("Insert", mock.AnythingOfType("*models.Order")).Return(nil)

	order := bytes.NewReader([]byte("{\"id\": \"INVALID\"}"))
	recorder := httptest.NewRecorder()
	request, _ := http.NewRequest("POST", "/api/order", order)
	tc.router.ServeHTTP(recorder, request)

	assert.Equal(400, recorder.Code)
}

func TestOrderUpdateSuccess(t *testing.T) {
	assert := assert.New(t)

	gin.SetMode(gin.TestMode)

	order := models.NewOrder("", "", "", "")

	tc, orderMock := mockOrder()
	orderMock.On("Update", order, order.ID.String()).Return(nil)

	recorder := httptest.NewRecorder()
	request, _ := http.NewRequest("PUT", "/api/order/"+order.ID.String(), getOrderString(order))
	tc.router.ServeHTTP(recorder, request)

	assert.Equal(200, recorder.Code)
}

func TestOrderUpdateFail(t *testing.T) {
	assert := assert.New(t)

	gin.SetMode(gin.TestMode)

	order := models.NewOrder("", "", "", "")

	tc, orderMock := mockOrder()
	orderMock.On("Update", order, order.ID.String()).Return(fmt.Errorf("This is an error"))

	recorder := httptest.NewRecorder()
	request, _ := http.NewRequest("PUT", "/api/order/"+order.ID.String(), getOrderString(order))
	tc.router.ServeHTTP(recorder, request)

	assert.Equal(500, recorder.Code)
}

func TestOrderUpdateInvalid(t *testing.T) {
	assert := assert.New(t)

	gin.SetMode(gin.TestMode)

	tc, orderMock := mockOrder()
	orderMock.On("Update", mock.AnythingOfType("*models.Order"), "").Return(fmt.Errorf("some error"))

	recorder := httptest.NewRecorder()
	request, _ := http.NewRequest("PUT", "/api/order/INVALID", bytes.NewReader([]byte("{\"id\": \"INVALID\"}")))
	tc.router.ServeHTTP(recorder, request)

	assert.Equal(400, recorder.Code)
}

func TestOrderDeleteSuccess(t *testing.T) {
	assert := assert.New(t)

	gin.SetMode(gin.TestMode)

	id := uuid.NewUUID()
	tc, orderMock := mockOrder()
	orderMock.On("Delete", id.String()).Return(nil)

	recorder := httptest.NewRecorder()
	request, _ := http.NewRequest("DELETE", "/api/order/"+id.String(), nil)
	tc.router.ServeHTTP(recorder, request)

	assert.Equal(200, recorder.Code)
}

func TestOrderDeleteFail(t *testing.T) {
	assert := assert.New(t)

	gin.SetMode(gin.TestMode)

	id := uuid.NewUUID()
	tc, orderMock := mockOrder()
	orderMock.On("Delete", id.String()).Return(fmt.Errorf("This is an error"))

	recorder := httptest.NewRecorder()
	request, _ := http.NewRequest("DELETE", "/api/order/"+id.String(), nil)
	tc.router.ServeHTTP(recorder, request)

	assert.Equal(500, recorder.Code)
}

func getOrderString(m *models.Order) io.Reader {
	s, _ := json.Marshal(m)
	return bytes.NewReader(s)
}
