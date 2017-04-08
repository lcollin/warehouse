package helpers

import (
	"database/sql"
	"fmt"
	"testing"

	"github.com/ghmeier/bloodlines/gateways"
	query "github.com/ghmeier/bloodlines/gateways/sql"
	"github.com/lcollin/warehouse/models"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/pborman/uuid"
	"github.com/stretchr/testify/assert"
)

func TestItemGetByID(t *testing.T) {
	assert := assert.New(t)

	id := uuid.NewUUID()
	s, mock, _ := sqlmock.New()
	r := getMockItem(s)

	mock.ExpectQuery("SELECT id, roasterID, name, pictureURL, coffeeType, inStockBags, providerPrice, consumerPrice, ozInBag FROM item").
		WithArgs(id.String()).
		WillReturnRows(getItemMockRows().AddRow(id.String(), "RoasterID", "Name", "PictureURL", "CoffeeType", "InStockBags", "ProviderPrice", "ConsumerPrice", "OzInBag"))

	item, err := r.GetByID(id.String())

	assert.Equal(mock.ExpectationsWereMet(), nil)
	assert.NoError(err)
	assert.Equal(item.ID, id)
	assert.Equal(item.RoasterID, "RoasterID")
	assert.Equal(item.Name, "Name")
	assert.Equal(item.CoffeeType, "CoffeeType")
	assert.Equal(item.InStockBags, "InStockBags")
	assert.Equal(item.ProviderPrice, "ProviderPrice")
	assert.Equal(item.ConsumerPrice, "ConsumerPrice")
	assert.Equal(item.OzInBag, "OzInBag")
}

func TestItemGetByIDError(t *testing.T) {
	assert := assert.New(t)

	id := uuid.NewUUID()
	s, mock, _ := sqlmock.New()
	r := getMockItem(s)

	mock.ExpectQuery("SELECT id, roasterID, name, pictureURL, coffeeType, inStockBags, providerPrice, consumerPrice, ozInBag FROM item").
		WithArgs(id.String()).
		WillReturnError(fmt.Errorf("This is an error"))

	_, err := r.GetByID(id.String())

	assert.Equal(mock.ExpectationsWereMet(), nil)
	assert.Error(err)
}

func TestItemGetAll(t *testing.T) {
	assert := assert.New(t)

	offset, limit := 0, 20
	s, mock, _ := sqlmock.New()
	r := getMockItem(s)

	mock.ExpectQuery("SELECT id, roasterID, name, pictureURL, coffeeType, inStockBags, providerPrice, consumerPrice, ozInBag FROM item").
		WithArgs(offset, limit).
		WillReturnRows(getItemMockRows().
			AddRow(uuid.New(), "RoasterID", "Name", "PictureURL", "CoffeeType", "InStockBags", "ProviderPrice", "ConsumerPrice", "OzInBag").
			AddRow(uuid.New(), "RoasterID", "Name", "PictureURL", "CoffeeType", "InStockBags", "ProviderPrice", "ConsumerPrice", "OzInBag"))

	items, err := r.GetAll(offset, limit, mockSelectQuery())

	assert.Equal(mock.ExpectationsWereMet(), nil)
	assert.NoError(err)
	assert.Equal(2, len(items))
}

func TestItemGetAllError(t *testing.T) {
	assert := assert.New(t)

	offset, limit := 0, 20
	s, mock, _ := sqlmock.New()
	r := getMockItem(s)

	mock.ExpectQuery("SELECT id, roasterID, name, pictureURL, coffeeType, inStockBags, providerPrice, consumerPrice, ozInBag FROM item").
		WithArgs(offset, limit).
		WillReturnError(fmt.Errorf("This is an error"))

	_, err := r.GetAll(offset, limit, &mockQuery{query: "SELECT id, roasterID, name, pictureURL, coffeeType, inStockBags, providerPrice, consumerPrice, ozInBag FROM item"})

	assert.Equal(mock.ExpectationsWereMet(), nil)
	assert.Error(err)
}

func TestItemInsert(t *testing.T) {
	assert := assert.New(t)

	item := getDefaultItem()
	s, mock, _ := sqlmock.New()
	r := getMockItem(s)

	mock.ExpectPrepare("INSERT INTO item").
		ExpectExec().
		WithArgs(item.ID.String(), item.RoasterID, item.Name, item.PictureURL, item.CoffeeType, item.InStockBags, item.ProviderPrice, item.ConsumerPrice, item.OzInBag).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err := r.Insert(item)

	assert.Equal(mock.ExpectationsWereMet(), nil)
	assert.NoError(err)
}

func TestItemInsertError(t *testing.T) {
	assert := assert.New(t)

	item := getDefaultItem()
	s, mock, _ := sqlmock.New()
	r := getMockItem(s)

	mock.ExpectPrepare("INSERT INTO item").
		ExpectExec().
		WithArgs(item.ID.String(), item.RoasterID, item.Name, item.PictureURL, item.CoffeeType, item.InStockBags, item.ProviderPrice, item.ConsumerPrice, item.OzInBag).
		WillReturnError(fmt.Errorf("This is an error"))

	err := r.Insert(item)

	assert.Equal(mock.ExpectationsWereMet(), nil)
	assert.Error(err)
}

func TestItemUpdate(t *testing.T) {
	assert := assert.New(t)

	item := getDefaultItem()
	s, mock, _ := sqlmock.New()
	r := getMockItem(s)

	mock.ExpectPrepare("UPDATE item").
		ExpectExec().
		WithArgs(item.RoasterID, item.Name, item.PictureURL, item.CoffeeType, item.InStockBags, item.ProviderPrice, item.ConsumerPrice, item.OzInBag, item.ID.String()).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err := r.Update(item, item.ID.String())

	assert.Equal(mock.ExpectationsWereMet(), nil)
	assert.NoError(err)
}

func TestItemUpdateError(t *testing.T) {
	assert := assert.New(t)

	item := getDefaultItem()
	s, mock, _ := sqlmock.New()
	r := getMockItem(s)

	mock.ExpectPrepare("UPDATE item").
		ExpectExec().
		WithArgs(item.RoasterID, item.Name, item.PictureURL, item.CoffeeType, item.InStockBags, item.ProviderPrice, item.ConsumerPrice, item.OzInBag, item.ID.String()).
		WillReturnError(fmt.Errorf("This is an error"))

	err := r.Update(item, item.ID.String())

	assert.Equal(mock.ExpectationsWereMet(), nil)
	assert.Error(err)
}

func TestDeleteItem(t *testing.T) {
	assert := assert.New(t)

	id := uuid.NewUUID()
	s, mock, _ := sqlmock.New()
	r := getMockItem(s)

	mock.ExpectPrepare("DELETE FROM item").
		ExpectExec().
		WithArgs(id.String()).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err := r.Delete(id.String())

	assert.Equal(mock.ExpectationsWereMet(), nil)
	assert.NoError(err)
}

func TestDeleteItemError(t *testing.T) {
	assert := assert.New(t)

	id := uuid.NewUUID()
	s, mock, _ := sqlmock.New()
	r := getMockItem(s)

	mock.ExpectPrepare("DELETE FROM item").
		ExpectExec().
		WithArgs(id.String()).
		WillReturnError(fmt.Errorf("This is an error"))

	err := r.Delete(id.String())

	assert.Equal(mock.ExpectationsWereMet(), nil)
	assert.Error(err)
}

func getDefaultItem() *models.Item {
	return models.NewItem(uuid.NewUUID(), "Name", "CoffeeType", "description", make([]string, 0), 0, 0.0, 0.0, 0.0, false, false)
}

func getItemMockRows() sqlmock.Rows {
	return sqlmock.NewRows([]string{"id", "roasterID", "name", "pictureURL", "coffeeType", "inStockBags", "providerPrice", "consumerPrice", "ozInBag", "description", "isDecaf", "isActive", "tags", "createdAt", "updatedAt"})
}

func getMockItem(s *sql.DB) *Item {
	return NewItem(&gateways.MySQL{DB: s}, nil)
}

func mockSelectQuery() *mockQuery {
	return &mockQuery{
		BaseSearch: &query.BaseSearch{},
		query:      "SELECT id, roasterID, name, pictureURL, coffeeType, inStockBags, providerPrice, consumerPrice, ozInBag FROM item",
	}
}

type mockQuery struct {
	*query.BaseSearch
	query string
}

func (i *mockQuery) ToQuery() string {
	return i.query
}
