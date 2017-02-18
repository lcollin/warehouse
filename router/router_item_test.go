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

func TestItemViewSuccess(t *testing.T) {
	assert := assert.New(t)

	gin.SetMode(gin.TestMode)

	id := uuid.NewUUID()
	tc, itemMock := mockItem()
	itemMock.On("GetByID", id.String()).Return(&models.Item{}, nil)

	recorder := httptest.NewRecorder()
	request, _ := http.NewRequest("GET", "/api/item/"+id.String(), nil)
	tc.router.ServeHTTP(recorder, request)

	assert.Equal(200, recorder.Code)
}

func TestItemViewFail(t *testing.T) {
	assert := assert.New(t)

	gin.SetMode(gin.TestMode)

	id := uuid.NewUUID()
	tc, itemMock := mockItem()
	itemMock.On("GetByID", id.String()).Return(nil, fmt.Errorf("This is an error"))

	recorder := httptest.NewRecorder()
	request, _ := http.NewRequest("GET", "/api/item/"+id.String(), nil)
	tc.router.ServeHTTP(recorder, request)

	assert.Equal(500, recorder.Code)
}

func TestItemViewAllSuccess(t *testing.T) {
	assert := assert.New(t)

	gin.SetMode(gin.TestMode)

	tc, itemMock := mockItem()
	itemMock.On("GetAll", 0, 20).Return(make([]*models.Item, 0), nil)

	recorder := httptest.NewRecorder()
	request, _ := http.NewRequest("GET", "/api/item", nil)
	tc.router.ServeHTTP(recorder, request)

	assert.Equal(200, recorder.Code)
}

func TestItemViewAllFail(t *testing.T) {
	assert := assert.New(t)

	gin.SetMode(gin.TestMode)

	tc, itemMock := mockItem()
	itemMock.On("GetAll", 0, 20).Return(make([]*models.Item, 0), fmt.Errorf("This is an error"))

	recorder := httptest.NewRecorder()
	request, _ := http.NewRequest("GET", "/api/item", nil)
	tc.router.ServeHTTP(recorder, request)

	assert.Equal(500, recorder.Code)
}

func TestItemViewAllParams(t *testing.T) {
	assert := assert.New(t)

	gin.SetMode(gin.TestMode)

	tc, itemMock := mockItem()
	itemMock.On("GetAll", 20, 40).Return(make([]*models.Item, 0), nil)

	recorder := httptest.NewRecorder()
	request, _ := http.NewRequest("GET", "/api/item?offset=20&limit=40", nil)
	tc.router.ServeHTTP(recorder, request)

	assert.Equal(200, recorder.Code)
}

func TestItemNewSuccess(t *testing.T) {
	assert := assert.New(t)

	gin.SetMode(gin.TestMode)

	tc, itemMock := mockItem()
	itemMock.On("Insert", mock.AnythingOfType("*models.Item")).Return(nil)

	item := getItemString(models.NewItem("", "", "", "", "", "", "", ""))
	recorder := httptest.NewRecorder()
	request, _ := http.NewRequest("POST", "/api/item", item)
	tc.router.ServeHTTP(recorder, request)

	assert.Equal(200, recorder.Code)
}

func TestItemNewFail(t *testing.T) {
	assert := assert.New(t)

	gin.SetMode(gin.TestMode)

	tc, itemMock := mockItem()
	itemMock.On("Insert", mock.AnythingOfType("*models.Item")).Return(fmt.Errorf("This is an error"))

	item := getItemString(models.NewItem("", "", "", "", "", "", "", ""))
	recorder := httptest.NewRecorder()
	request, _ := http.NewRequest("POST", "/api/item", item)
	tc.router.ServeHTTP(recorder, request)

	assert.Equal(500, recorder.Code)
}

func TestItemNewInvalid(t *testing.T) {
	assert := assert.New(t)

	gin.SetMode(gin.TestMode)

	tc, itemMock := mockItem()
	itemMock.On("Insert", mock.AnythingOfType("*models.Item")).Return(nil)

	item := bytes.NewReader([]byte("{\"id\": \"INVALID\"}"))
	recorder := httptest.NewRecorder()
	request, _ := http.NewRequest("POST", "/api/item", item)
	tc.router.ServeHTTP(recorder, request)

	assert.Equal(400, recorder.Code)
}

func TestItemUpdateSuccess(t *testing.T) {
	assert := assert.New(t)

	gin.SetMode(gin.TestMode)

	item := models.NewItem("", "", "", "", "", "", "", "")

	tc, itemMock := mockItem()
	itemMock.On("Update", item, item.ID.String()).Return(nil)

	recorder := httptest.NewRecorder()
	request, _ := http.NewRequest("PUT", "/api/item/"+item.ID.String(), getItemString(item))
	tc.router.ServeHTTP(recorder, request)

	assert.Equal(200, recorder.Code)
}

func TestItemUpdateFail(t *testing.T) {
	assert := assert.New(t)

	gin.SetMode(gin.TestMode)

	item := models.NewItem("", "", "", "", "", "", "", "")

	tc, itemMock := mockItem()
	itemMock.On("Update", item, item.ID.String()).Return(fmt.Errorf("This is an error"))

	recorder := httptest.NewRecorder()
	request, _ := http.NewRequest("PUT", "/api/item/"+item.ID.String(), getItemString(item))
	tc.router.ServeHTTP(recorder, request)

	assert.Equal(500, recorder.Code)
}

func TestItemUpdateInvalid(t *testing.T) {
	assert := assert.New(t)

	gin.SetMode(gin.TestMode)

	tc, itemMock := mockItem()
	itemMock.On("Update", mock.AnythingOfType("*models.Item"), "").Return(fmt.Errorf("some error"))

	recorder := httptest.NewRecorder()
	request, _ := http.NewRequest("PUT", "/api/item/INVALID", bytes.NewReader([]byte("{\"id\": \"INVALID\"}")))
	tc.router.ServeHTTP(recorder, request)

	assert.Equal(400, recorder.Code)
}

func TestItemDeleteSuccess(t *testing.T) {
	assert := assert.New(t)

	gin.SetMode(gin.TestMode)

	id := uuid.NewUUID()
	tc, itemMock := mockItem()
	itemMock.On("Delete", id.String()).Return(nil)

	recorder := httptest.NewRecorder()
	request, _ := http.NewRequest("DELETE", "/api/item/"+id.String(), nil)
	tc.router.ServeHTTP(recorder, request)

	assert.Equal(200, recorder.Code)
}

func TestItemDeleteFail(t *testing.T) {
	assert := assert.New(t)

	gin.SetMode(gin.TestMode)

	id := uuid.NewUUID()
	tc, itemMock := mockItem()
	itemMock.On("Delete", id.String()).Return(fmt.Errorf("This is an error"))

	recorder := httptest.NewRecorder()
	request, _ := http.NewRequest("DELETE", "/api/item/"+id.String(), nil)
	tc.router.ServeHTTP(recorder, request)

	assert.Equal(500, recorder.Code)
}

func getItemString(m *models.Item) io.Reader {
	s, _ := json.Marshal(m)
	return bytes.NewReader(s)
}
