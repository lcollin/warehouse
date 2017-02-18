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

func TestSubOrderViewSuccess(t *testing.T) {
	assert := assert.New(t)

	gin.SetMode(gin.TestMode)

	id := uuid.NewUUID()
	tc, suborderMock := mockSubOrder()
	suborderMock.On("GetByID", id.String()).Return(&models.SubOrder{}, nil)

	recsuborder := httptest.NewRecsuborder()
	request, _ := http.NewRequest("GET", "/api/suborder/"+id.String(), nil)
	tc.router.ServeHTTP(recsuborder, request)

	assert.Equal(200, recsuborder.Code)
}

func TestSubOrderViewFail(t *testing.T) {
	assert := assert.New(t)

	gin.SetMode(gin.TestMode)

	id := uuid.NewUUID()
	tc, suborderMock := mockSubOrder()
	suborderMock.On("GetByID", id.String()).Return(nil, fmt.Errorf("This is an error"))

	recsuborder := httptest.NewRecsuborder()
	request, _ := http.NewRequest("GET", "/api/suborder/"+id.String(), nil)
	tc.router.ServeHTTP(recsuborder, request)

	assert.Equal(500, recsuborder.Code)
}

func TestSubOrderViewAllSuccess(t *testing.T) {
	assert := assert.New(t)

	gin.SetMode(gin.TestMode)

	tc, suborderMock := mockSubOrder()
	suborderMock.On("GetAll", 0, 20).Return(make([]*models.SubOrder, 0), nil)

	recsuborder := httptest.NewRecsuborder()
	request, _ := http.NewRequest("GET", "/api/suborder", nil)
	tc.router.ServeHTTP(recsuborder, request)

	assert.Equal(200, recsuborder.Code)
}

func TestSubOrderViewAllFail(t *testing.T) {
	assert := assert.New(t)

	gin.SetMode(gin.TestMode)

	tc, suborderMock := mockSubOrder()
	suborderMock.On("GetAll", 0, 20).Return(make([]*models.SubOrder, 0), fmt.Errorf("This is an error"))

	recsuborder := httptest.NewRecsuborder()
	request, _ := http.NewRequest("GET", "/api/suborder", nil)
	tc.router.ServeHTTP(recsuborder, request)

	assert.Equal(500, recsuborder.Code)
}

func TestSubOrderViewAllParams(t *testing.T) {
	assert := assert.New(t)

	gin.SetMode(gin.TestMode)

	tc, suborderMock := mockSubOrder()
	suborderMock.On("GetAll", 20, 40).Return(make([]*models.SubOrder, 0), nil)

	recsuborder := httptest.NewRecsuborder()
	request, _ := http.NewRequest("GET", "/api/suborder?offset=20&limit=40", nil)
	tc.router.ServeHTTP(recsuborder, request)

	assert.Equal(200, recsuborder.Code)
}

func TestSubOrderNewSuccess(t *testing.T) {
	assert := assert.New(t)

	gin.SetMode(gin.TestMode)

	tc, suborderMock := mockSubOrder()
	suborderMock.On("Insert", mock.AnythingOfType("*models.SubOrder")).Return(nil)

	suborder := getSubOrderString(models.NewSubOrder("", "", ""))
	recsuborder := httptest.NewRecsuborder()
	request, _ := http.NewRequest("POST", "/api/suborder", suborder)
	tc.router.ServeHTTP(recsuborder, request)

	assert.Equal(200, recsuborder.Code)
}

func TestSubOrderNewFail(t *testing.T) {
	assert := assert.New(t)

	gin.SetMode(gin.TestMode)

	tc, suborderMock := mockSubOrder()
	suborderMock.On("Insert", mock.AnythingOfType("*models.SubOrder")).Return(fmt.Errorf("This is an error"))

	suborder := getSubOrderString(models.NewSubOrder("", "", ""))
	recsuborder := httptest.NewRecsuborder()
	request, _ := http.NewRequest("POST", "/api/suborder", suborder)
	tc.router.ServeHTTP(recsuborder, request)

	assert.Equal(500, recsuborder.Code)
}

func TestSubOrderNewInvalid(t *testing.T) {
	assert := assert.New(t)

	gin.SetMode(gin.TestMode)

	tc, suborderMock := mockSubOrder()
	suborderMock.On("Insert", mock.AnythingOfType("*models.SubOrder")).Return(nil)

	suborder := bytes.NewReader([]byte("{\"id\": \"INVALID\"}"))
	recsuborder := httptest.NewRecsuborder()
	request, _ := http.NewRequest("POST", "/api/suborder", suborder)
	tc.router.ServeHTTP(recsuborder, request)

	assert.Equal(400, recsuborder.Code)
}

func TestSubOrderUpdateSuccess(t *testing.T) {
	assert := assert.New(t)

	gin.SetMode(gin.TestMode)

	suborder := models.NewSubOrder("", "", "")

	tc, suborderMock := mockSubOrder()
	suborderMock.On("Update", suborder, suborder.ID.String()).Return(nil)

	recsuborder := httptest.NewRecsuborder()
	request, _ := http.NewRequest("PUT", "/api/suborder/"+suborder.ID.String(), getSubOrderString(suborder))
	tc.router.ServeHTTP(recsuborder, request)

	assert.Equal(200, recsuborder.Code)
}

func TestSubOrderUpdateFail(t *testing.T) {
	assert := assert.New(t)

	gin.SetMode(gin.TestMode)

	suborder := models.NewSubOrder("", "", "")

	tc, suborderMock := mockSubOrder()
	suborderMock.On("Update", suborder, suborder.ID.String()).Return(fmt.Errorf("This is an error"))

	recsuborder := httptest.NewRecsuborder()
	request, _ := http.NewRequest("PUT", "/api/suborder/"+suborder.ID.String(), getSubOrderString(suborder))
	tc.router.ServeHTTP(recsuborder, request)

	assert.Equal(500, recsuborder.Code)
}

func TestSubOrderUpdateInvalid(t *testing.T) {
	assert := assert.New(t)

	gin.SetMode(gin.TestMode)

	tc, suborderMock := mockSubOrder()
	suborderMock.On("Update", mock.AnythingOfType("*models.SubOrder"), "").Return(fmt.Errorf("some error"))

	recsuborder := httptest.NewRecsuborder()
	request, _ := http.NewRequest("PUT", "/api/suborder/INVALID", bytes.NewReader([]byte("{\"id\": \"INVALID\"}")))
	tc.router.ServeHTTP(recsuborder, request)

	assert.Equal(400, recsuborder.Code)
}

func TestSubOrderDeleteSuccess(t *testing.T) {
	assert := assert.New(t)

	gin.SetMode(gin.TestMode)

	id := uuid.NewUUID()
	tc, suborderMock := mockSubOrder()
	suborderMock.On("Delete", id.String()).Return(nil)

	recsuborder := httptest.NewRecsuborder()
	request, _ := http.NewRequest("DELETE", "/api/suborder/"+id.String(), nil)
	tc.router.ServeHTTP(recsuborder, request)

	assert.Equal(200, recsuborder.Code)
}

func TestSubOrderDeleteFail(t *testing.T) {
	assert := assert.New(t)

	gin.SetMode(gin.TestMode)

	id := uuid.NewUUID()
	tc, suborderMock := mockSubOrder()
	suborderMock.On("Delete", id.String()).Return(fmt.Errorf("This is an error"))

	recsuborder := httptest.NewRecsuborder()
	request, _ := http.NewRequest("DELETE", "/api/suborder/"+id.String(), nil)
	tc.router.ServeHTTP(recsuborder, request)

	assert.Equal(500, recsuborder.Code)
}

func getSubOrderString(m *models.SubOrder) io.Reader {
	s, _ := json.Marshal(m)
	return bytes.NewReader(s)
}
