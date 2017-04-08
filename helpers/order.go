package helpers

import (
	"github.com/ghmeier/bloodlines/gateways"
	"github.com/lcollin/warehouse/models"

	"github.com/pborman/uuid"
)

const SELECT_ALL = "SELECT id, userID, subscriptionID, requestDate, shipDate, quantity, status "

type baseHelper struct {
	sql gateways.SQL
}

type OrderI interface {
	GetByID(string) (*models.Order, error)
	GetByUserID(string) ([]*models.Order, error)
	GetAll(int, int) ([]*models.Order, error)
	Insert(*models.Order) error
	Update(*models.Order) error
	SetStatus(id uuid.UUID, status models.OrderStatus) error
	Delete(string) error
}

type Order struct {
	*baseHelper
}

func NewOrder(sql gateways.SQL) *Order {
	return &Order{baseHelper: &baseHelper{sql: sql}}
}

func (i *Order) GetByID(id string) (*models.Order, error) {
	rows, err := i.sql.Select(SELECT_ALL+" FROM orderT WHERE id=?", id)
	if err != nil {
		return nil, err
	}

	items, err := models.OrderFromSQL(rows)
	if err != nil {
		return nil, err
	}

	if len(items) <= 0 {
		return nil, nil
	}

	return items[0], err
}

func (i *Order) GetByUserID(userID string) ([]*models.Order, error) {
	rows, err := i.sql.Select(SELECT_ALL+" FROM orderT WHERE userID=? ORDER BY status ASC, id ASC", userID)
	if err != nil {
		return nil, err
	}

	items, err := models.OrderFromSQL(rows)
	if err != nil {
		return nil, err
	}

	return items, err
}

func (i *Order) GetAll(offset int, limit int) ([]*models.Order, error) {
	rows, err := i.sql.Select(SELECT_ALL+" FROM orderT ORDER BY status ASC, id ASC LIMIT ?,?", offset, limit)
	if err != nil {
		return nil, err
	}

	items, err := models.OrderFromSQL(rows)
	if err != nil {
		return nil, err
	}

	return items, err
}

func (i *Order) Insert(order *models.Order) error {
	err := i.sql.Modify(
		"INSERT INTO orderT (id, userID, subscriptionID, requestDate, shipDate, quantity, status) VALUE (?,?,?,?,?,?,?)",
		order.ID,
		order.UserID,
		order.SubscriptionID,
		order.RequestDate,
		order.ShipDate,
		order.Quantity,
		order.Status,
	)

	return err
}

func (i *Order) Update(order *models.Order) error {
	err := i.sql.Modify(
		"UPDATE orderT SET userID=?, subscriptionID=?, requestDate=?, shipDate=?, quantity=?, status=? WHERE id=?",
		order.UserID,
		order.SubscriptionID,
		order.RequestDate,
		order.ShipDate,
		order.Quantity,
		order.Status,
		order.ID.String(),
	)

	return err
}

func (i *Order) SetStatus(id uuid.UUID, status models.OrderStatus) error {
	err := i.sql.Modify("UPDATE orderT SET status=? WHERE id=?", string(status), id.String())
	return err
}

func (i *Order) Delete(id string) error {
	err := i.sql.Modify("DELETE FROM orderT WHERE id=?", id)
	return err
}
