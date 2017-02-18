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

func TestSubOrderGetByID(t *testing.T) {
	assert := assert.New(t)

	id := uuid.NewUUID()
	s, mock, _ := sqlmock.New()
	r := getMockSubOrder(s)

	mock.ExpectQuery("SELECT id, suborderID, itemID, quantity FROM suborder").
		WithArgs(id.String()).
		WillReturnRows(getSubOrderMockRows().AddRow(id.String(), "OrderID", "ItemID", "Quantity"))

	suborder, err := r.GetByID(id.String())

	assert.Equal(mock.ExpectationsWereMet(), nil)
	assert.NoError(err)
	assert.Equal(suborder.ID, id)
	assert.Equal(suborder.RoasterID, "OrderID")
	assert.Equal(suborder.Name, "ItemID")
	assert.Equal(suborder.CoffeeType, "Quantity")
}

func TestSubOrderGetByIDError(t *testing.T) {
	assert := assert.New(t)

	id := uuid.NewUUID()
	s, mock, _ := sqlmock.New()
	r := getMockSubOrder(s)

	mock.ExpectQuery("SELECT id, suborderID, itemID, quantity FROM suborder").
		WithArgs(id.String()).
		WillReturnError(fmt.Errorf("This is an error"))

	_, err := r.GetByID(id.String())

	assert.Equal(mock.ExpectationsWereMet(), nil)
	assert.Error(err)
}

func TestSubOrderGetAll(t *testing.T) {
	assert := assert.New(t)

	offset, limit := 0, 20
	s, mock, _ := sqlmock.New()
	r := getMockSubOrder(s)

	mock.ExpectQuery("SELECT id, orderID, itemID, quantity FROM suborder").
		WithArgs(offset, limit).
		WillReturnRows(getSubOrderMockRows().
			AddRow(uuid.New(), "OrderID", "ItemID", "Quantity").
			AddRow(uuid.New(), "OrderID", "ItemID", "Quantity"))

	suborders, err := r.GetAll(offset, limit)

	assert.Equal(mock.ExpectationsWereMet(), nil)
	assert.NoError(err)
	assert.Equal(2, len(suborders))
}

func TestSubOrderGetAllError(t *testing.T) {
	assert := assert.New(t)

	offset, limit := 0, 20
	s, mock, _ := sqlmock.New()
	r := getMockSubOrder(s)

	mock.ExpectQuery("SELECT id, orderID, itemID, quantity FROM suborder").
		WithArgs(offset, limit).
		WillReturnError(fmt.Errorf("This is an error"))

	_, err := r.GetAll(offset, limit)

	assert.Equal(mock.ExpectationsWereMet(), nil)
	assert.Error(err)
}

func TestSubOrderInsert(t *testing.T) {
	assert := assert.New(t)

	suborder := getDefaultSubOrder()
	s, mock, _ := sqlmock.New()
	r := getMockSubOrder(s)

	mock.ExpectPrepare("INSERT INTO suborder").
		ExpectExec().
		WithArgs(suborder.ID.String(), suborder.OrderID, suborder.ItemID, suborder.Quantity).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err := r.Insert(suborder)

	assert.Equal(mock.ExpectationsWereMet(), nil)
	assert.NoError(err)
}

func TestSubOrderInsertError(t *testing.T) {
	assert := assert.New(t)

	suborder := getDefaultSubOrder()
	s, mock, _ := sqlmock.New()
	r := getMockSubOrder(s)

	mock.ExpectPrepare("INSERT INTO suborder").
		ExpectExec().
		WithArgs(suborder.ID.String(), suborder.OrderID, suborder.ItemID, suborder.Quantity).
		WillReturnError(fmt.Errorf("This is an error"))

	err := r.Insert(suborder)

	assert.Equal(mock.ExpectationsWereMet(), nil)
	assert.Error(err)
}

func TestSubOrderUpdate(t *testing.T) {
	assert := assert.New(t)

	suborder := getDefaultSubOrder()
	s, mock, _ := sqlmock.New()
	r := getMockSubOrder(s)

	mock.ExpectPrepare("UPDATE suborder").
		ExpectExec().
		WithArgs(suborder.OrderID, suborder.ItemID, suborder.Quantity, suborder.ID.String()).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err := r.Update(suborder, suborder.ID.String())

	assert.Equal(mock.ExpectationsWereMet(), nil)
	assert.NoError(err)
}

func TestSubOrderUpdateError(t *testing.T) {
	assert := assert.New(t)

	suborder := getDefaultSubOrder()
	s, mock, _ := sqlmock.New()
	r := getMockSubOrder(s)

	mock.ExpectPrepare("UPDATE suborder").
		ExpectExec().
		WithArgs(suborder.OrderID, suborder.ItemID, suborder.Quantity, suborder.ID.String()).
		WillReturnError(fmt.Errorf("This is an error"))

	err := r.Update(suborder, suborder.ID.String())

	assert.Equal(mock.ExpectationsWereMet(), nil)
	assert.Error(err)
}

func TestDeleteSubOrder(t *testing.T) {
	assert := assert.New(t)

	id := uuid.NewUUID()
	s, mock, _ := sqlmock.New()
	r := getMockSubOrder(s)

	mock.ExpectPrepare("DELETE FROM suborder").
		ExpectExec().
		WithArgs(id.String()).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err := r.Delete(id.String())

	assert.Equal(mock.ExpectationsWereMet(), nil)
	assert.NoError(err)
}

func TestDeleteSubOrderError(t *testing.T) {
	assert := assert.New(t)

	id := uuid.NewUUID()
	s, mock, _ := sqlmock.New()
	r := getMockSubOrder(s)

	mock.ExpectPrepare("DELETE FROM suborder").
		ExpectExec().
		WithArgs(id.String()).
		WillReturnError(fmt.Errorf("This is an error"))

	err := r.Delete(id.String())

	assert.Equal(mock.ExpectationsWereMet(), nil)
	assert.Error(err)
}

func getDefaultSubOrder() *models.SubOrder {
	return models.NewSubOrder("OrderID", "ItemID", "Quantity")
}

func getSubOrderMockRows() sqlmock.Rows {
	return sqlmock.NewRows([]string{"id", "orderID", "itemID", "quantity"})
}

func getMockSubOrder(s *sql.DB) *SubOrder {
	return NewSubOrder(&gateways.MySQL{DB: s})
}
