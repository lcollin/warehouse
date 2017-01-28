package helpers

import (
	"github.com/ghmeier/bloodlines/gateways"
	"github.com/lcollin/warehouse/containers"
)

type baseHelper struct {
	sql gateways.SQL
}

type OrderI interface {
	GetByID(string) (*containers.Order, error)
	GetByUserID(string) (*containers.Order, error)
	GetAll(int, int) ([]*containers.Order, error)
	Insert(*containers.Order) error
	Update(*containers.Order, string) error
	Delete(string) error
}

type Order struct {
	*baseHelper
}

func NewOrder(sql gateways.SQL) *Order {
	return &Order{baseHelper: &baseHelper{sql: sql}}
}

func (i *Order) GetByID(id string) (*containers.Order, error) {
	rows, err := i.sql.Select("SELECT * FROM order WHERE id=?", id)
	if err != nil {
		return nil, err
	}

	items, err := containers.OrderFromSQL(rows)
	if err != nil {
		return nil, err
	}

	return items[0], err
}

func (i *Order) GetByUserID(userID string) (*containers.Order, error) {
	rows, err := i.sql.Select("SELECT * FROM order WHERE user_id=?", userID)
	if err != nil {
		return nil, err
	}

	items, err := containers.OrderFromSQL(rows)
	if err != nil {
		return nil, err
	}

	return items[0], err
}

func (i *Order) GetAll(offset int, limit int) ([]*containers.Order, error) {
	rows, err := i.sql.Select("SELECT * FROM order ORDER BY id ASC LIMIT ?,?", offset, limit)
	if err != nil {
		return nil, err
	}

	items, err := containers.OrderFromSQL(rows)
	if err != nil {
		return nil, err
	}

	return items, err
}

func (i *Order) Insert(order *containers.Order) error {
	err := i.sql.Modify(
		"INSERT INTO order (id, user_id) VALUE (?,?)",
		order.ID,
		order.UserID,
	)

	return err
}

func (i *Order) Update(order *containers.Order, id string) error {
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
