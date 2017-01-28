package helpers

import (
	"github.com/ghmeier/bloodlines/gateways"
	"github.com/lcollin/warehouse/containers"
)

type baseHelper struct {
	sql gateways.SQL
}

type SubOrderI interface {
	GetByID(string) (*containers.SubOrder, error)
	GetBySubOrderID(string) (*containers.SubOrder, error)
	GetAll(int, int) ([]*containers.SubOrder, error)
	Insert(*containers.SubOrder) error
	Update(*containers.SubOrder, string) error
	Delete(string) error
}

type SubOrder struct {
	*baseHelper
}

func NewSubOrder(sql gateways.SQL) *SubOrder {
	return &SubOrder{baseHelper: &baseHelper{sql: sql}}
}

func (i *SubOrder) GetByID(id string) (*containers.SubOrder, error) {
	rows, err := i.sql.Select("SELECT * FROM suborder WHERE id=?", id)
	if err != nil {
		return nil, err
	}

	suborders, err := containers.SubOrderFromSQL(rows)
	if err != nil {
		return nil, err
	}

	return suborders[0], err
}

func (i *SubOrder) GetByOrderID(orderID string) (*containers.SubOrder, error) {
	rows, err := i.sql.Select("SELECT * FROM suborder WHERE order_id=?", orderID)
	if err != nil {
		return nil, err
	}

	suborders, err := containers.SubOrderFromSQL(rows)
	if err != nil {
		return nil, err
	}

	return suborders[0], err
}

func (i *SubOrder) GetAll(offset int, limit int) ([]*containers.SubOrder, error) {
	rows, err := i.sql.Select("SELECT * FROM suborder ORDER BY id ASC LIMIT ?,?", offset, limit)
	if err != nil {
		return nil, err
	}

	suborders, err := containers.SubOrderFromSQL(rows)
	if err != nil {
		return nil, err
	}

	return suborders, err
}

func (i *SubOrder) Insert(suborder *containers.SubOrder) error {
	err := i.sql.Modify(
		"INSERT INTO suborder (id, order_id, item_id) VALUE (?,?,?)",
		suborder.ID,
		suborder.OrderID,
		suborder.ItemID,
	)

	return err
}

func (i *SubOrder) Update(suborder *containers.SubOrder, id string) error {
	err := i.sql.Modify(
		"UPDATE suborder SET user_id=?, item_id=? WHERE id=?",
		suborder.OrderID,
		suborder.ItemID,
		id,
	)

	return err
}

func (i *SubOrder) Delete(id string) error {
	err := i.sql.Modify("DELETE FROM suborder WHERE id=?", id)
	return err
}
