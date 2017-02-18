package helpers

import (
	"database/sql"
	"fmt"
	"testing"

	"github.com/ghmeier/bloodlines/gateways"
	"github.com/lcollin/warehouse/models"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/pborman/uuid"
	"github.com/stretchr/testify/assert"
)

func TestOrderGetByID(t *testing.T) {
	assert := assert.New(t)

	id := uuid.NewUUID()
	s, mock, _ := sqlmock.New()
	r := getMockOrder(s)

	mock.ExpectQuery("SELECT id, userID, subscriptionID, requestDate, shipDate FROM order").
		WithArgs(id.String()).
		WillReturnRows(getOrderMockRows().AddRow(id.String(), "UserID", "SubscriptionID", "RequestDate", "ShipDate"))

	order, err := r.GetByID(id.String())

	assert.Equal(mock.ExpectationsWereMet(), nil)
	assert.NoError(err)
	assert.Equal(order.ID, id)
	assert.Equal(order.RoasterID, "UserID")
	assert.Equal(order.Name, "SubscriptionID")
	assert.Equal(order.CoffeeType, "RequestDate")
	assert.Equal(order.InStockBags, "ShipDate")
}

func TestOrderGetByIDError(t *testing.T) {
	assert := assert.New(t)

	id := uuid.NewUUID()
	s, mock, _ := sqlmock.New()
	r := getMockOrder(s)

	mock.ExpectQuery("SELECT id, userID, subscriptionID, requestDate, shipDate FROM order").
		WithArgs(id.String()).
		WillReturnError(fmt.Errorf("This is an error"))

	_, err := r.GetByID(id.String())

	assert.Equal(mock.ExpectationsWereMet(), nil)
	assert.Error(err)
}

func TestOrderGetAll(t *testing.T) {
	assert := assert.New(t)

	offset, limit := 0, 20
	s, mock, _ := sqlmock.New()
	r := getMockOrder(s)

	mock.ExpectQuery("SELECT id, userID, subscriptionID, requestDate, shipDate FROM order").
		WithArgs(offset, limit).
		WillReturnRows(getOrderMockRows().
			AddRow(uuid.New(), "UserID", "SubscriptionID", "RequestDate", "ShipDate").
			AddRow(uuid.New(), "UserID", "SubscriptionID", "RequestDate", "ShipDate"))

	orders, err := r.GetAll(offset, limit)

	assert.Equal(mock.ExpectationsWereMet(), nil)
	assert.NoError(err)
	assert.Equal(2, len(orders))
}

func TestOrderGetAllError(t *testing.T) {
	assert := assert.New(t)

	offset, limit := 0, 20
	s, mock, _ := sqlmock.New()
	r := getMockOrder(s)

	mock.ExpectQuery("SELECT userID, subscriptionID, requestDate, shipDate FROM order").
		WithArgs(offset, limit).
		WillReturnError(fmt.Errorf("This is an error"))

	_, err := r.GetAll(offset, limit)

	assert.Equal(mock.ExpectationsWereMet(), nil)
	assert.Error(err)
}

func TestOrderInsert(t *testing.T) {
	assert := assert.New(t)

	order := getDefaultOrder()
	s, mock, _ := sqlmock.New()
	r := getMockOrder(s)

	mock.ExpectPrepare("INSERT INTO order").
		ExpectExec().
		WithArgs(order.ID.String(), order.UserID, order.SubscriptionID, order.RequestDate, order.ShipDate).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err := r.Insert(order)

	assert.Equal(mock.ExpectationsWereMet(), nil)
	assert.NoError(err)
}

func TestOrderInsertError(t *testing.T) {
	assert := assert.New(t)

	order := getDefaultOrder()
	s, mock, _ := sqlmock.New()
	r := getMockOrder(s)

	mock.ExpectPrepare("INSERT INTO order").
		ExpectExec().
		WithArgs(order.ID.String(), order.UserID, order.SubscriptionID, order.RequestDate, order.ShipDate).
		WillReturnError(fmt.Errorf("This is an error"))

	err := r.Insert(order)

	assert.Equal(mock.ExpectationsWereMet(), nil)
	assert.Error(err)
}

func TestOrderUpdate(t *testing.T) {
	assert := assert.New(t)

	order := getDefaultOrder()
	s, mock, _ := sqlmock.New()
	r := getMockOrder(s)

	mock.ExpectPrepare("UPDATE order").
		ExpectExec().
		WithArgs(order.UserID, order.SubscriptionID, order.RequestDate, order.ShipDate, order.ID.String()).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err := r.Update(order, order.ID.String())

	assert.Equal(mock.ExpectationsWereMet(), nil)
	assert.NoError(err)
}

func TestOrderUpdateError(t *testing.T) {
	assert := assert.New(t)

	order := getDefaultOrder()
	s, mock, _ := sqlmock.New()
	r := getMockOrder(s)

	mock.ExpectPrepare("UPDATE order").
		ExpectExec().
		WithArgs(order.UserID, order.SubscriptionID, order.RequestDate, order.ShipDate, order.ID.String()).
		WillReturnError(fmt.Errorf("This is an error"))

	err := r.Update(order, order.ID.String())

	assert.Equal(mock.ExpectationsWereMet(), nil)
	assert.Error(err)
}

func TestDeleteOrder(t *testing.T) {
	assert := assert.New(t)

	id := uuid.NewUUID()
	s, mock, _ := sqlmock.New()
	r := getMockOrder(s)

	mock.ExpectPrepare("DELETE FROM order").
		ExpectExec().
		WithArgs(id.String()).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err := r.Delete(id.String())

	assert.Equal(mock.ExpectationsWereMet(), nil)
	assert.NoError(err)
}

func TestDeleteOrderError(t *testing.T) {
	assert := assert.New(t)

	id := uuid.NewUUID()
	s, mock, _ := sqlmock.New()
	r := getMockOrder(s)

	mock.ExpectPrepare("DELETE FROM order").
		ExpectExec().
		WithArgs(id.String()).
		WillReturnError(fmt.Errorf("This is an error"))

	err := r.Delete(id.String())

	assert.Equal(mock.ExpectationsWereMet(), nil)
	assert.Error(err)
}

func getDefaultOrder() *models.Order {
	return models.NewOrder("UserID", "SubscriptionID", "RequestDate", "ShipDate")
}

func getOrderMockRows() sqlmock.Rows {
	return sqlmock.NewRows([]string{"id", "userID", "subscriptionID", "requestDate", "shipDate"})
}

func getMockOrder(s *sql.DB) *Order {
	return NewOrder(&gateways.MySQL{DB: s})
}
