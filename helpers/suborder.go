package helpers

import (
	"github.com/ghmeier/bloodlines/gateways"
	"github.com/lcollin/warehouse/models"
)

type SubOrderI interface {
	GetByID(string) (*models.SubOrder, error)
	GetBySubOrderID(string) (*models.SubOrder, error)
	GetAll(int, int) ([]*models.SubOrder, error)
	Insert(*models.SubOrder) error
	Update(*models.SubOrder, string) error
	Delete(string) error
}

type SubOrder struct {
	*baseHelper
}

func NewSubOrder(sql gateways.SQL) *SubOrder {
	return &SubOrder{baseHelper: &baseHelper{sql: sql}}
}

func (i *SubOrder) GetByID(id string) (*models.SubOrder, error) {
	rows, err := i.sql.Select("SELECT * FROM suborder WHERE id=?", id)
	if err != nil {
		return nil, err
	}

	suborders, err := models.SubOrderFromSQL(rows)
	if err != nil {
		return nil, err
	}

	return suborders[0], err
}

func (i *SubOrder) GetBySubOrderID(orderID string) (*models.SubOrder, error) {
	rows, err := i.sql.Select("SELECT * FROM suborder WHERE order_id=?", orderID)
	if err != nil {
		return nil, err
	}

	suborders, err := models.SubOrderFromSQL(rows)
	if err != nil {
		return nil, err
	}

	return suborders[0], err
}

func (i *SubOrder) GetAll(offset int, limit int) ([]*models.SubOrder, error) {
	rows, err := i.sql.Select("SELECT * FROM suborder ORDER BY id ASC LIMIT ?,?", offset, limit)
	if err != nil {
		return nil, err
	}

	suborders, err := models.SubOrderFromSQL(rows)
	if err != nil {
		return nil, err
	}

	return suborders, err
}

func (i *SubOrder) Insert(suborder *models.SubOrder) error {
	err := i.sql.Modify(
		"INSERT INTO suborder (id, order_id, item_id) VALUE (?,?,?)",
		suborder.ID,
		suborder.OrderID,
		suborder.ItemID,
	)

	return err
}

func (i *SubOrder) Update(suborder *models.SubOrder, id string) error {
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
