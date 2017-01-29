package helpers

import (
	"github.com/ghmeier/bloodlines/gateways"
	"github.com/lcollin/warehouse/models"
)

type baseHelper struct {
	sql gateways.SQL
}

type OrderI interface {
	GetByID(string) (*models.Order, error)
	GetByUserID(string) (*models.Order, error)
	GetAll(int, int) ([]*models.Order, error)
	Insert(*models.Order) error
	Update(*models.Order, string) error
	Delete(string) error
}

type Order struct {
	*baseHelper
}

func NewOrder(sql gateways.SQL) *Order {
	return &Order{baseHelper: &baseHelper{sql: sql}}
}

func (i *Order) GetByID(id string) (*models.Order, error) {
	rows, err := i.sql.Select("SELECT * FROM order WHERE id=?", id)
	if err != nil {
		return nil, err
	}

	items, err := models.OrderFromSQL(rows)
	if err != nil {
		return nil, err
	}

	return items[0], err
}

func (i *Order) GetByUserID(userID string) (*models.Order, error) {
	rows, err := i.sql.Select("SELECT * FROM order WHERE user_id=?", userID)
	if err != nil {
		return nil, err
	}

	items, err := models.OrderFromSQL(rows)
	if err != nil {
		return nil, err
	}

	return items[0], err
}

func (i *Order) GetAll(offset int, limit int) ([]*models.Order, error) {
	rows, err := i.sql.Select("SELECT * FROM order ORDER BY id ASC LIMIT ?,?", offset, limit)
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
		"INSERT INTO order (id, user_id) VALUE (?,?)",
		order.ID,
		order.UserID,
	)

	return err
}

func (i *Order) Update(order *models.Order, id string) error {
	err := i.sql.Modify(
		"UPDATE order SET user_id=? WHERE id=?",
		order.UserID,
		id,
	)

	return err
}

func (i *Order) Delete(id string) error {
	err := i.sql.Modify("DELETE FROM order WHERE id=?", id)
	return err
}
