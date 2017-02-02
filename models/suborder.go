package models

import (
	"database/sql"
	"github.com/pborman/uuid"
)

type SubOrder struct {
	ID      uuid.UUID `json: "id"`
	OrderID uuid.UUID `json: "orderID"`
	ItemID  uuid.UUID `json: "itemID"`
	Quantity int `json: "quantity"`
}

func NewSubOrder(orderID uuid.UUID, itemID uuid.UUID, quantity int) *SubOrder {
	return &SubOrder{
		ID:      uuid.NewUUID(),
		OrderID: orderID,
		ItemID:  itemID,
		Quantity: quantity,
	}
}

func SubOrderFromSQL(rows *sql.Rows) ([]*SubOrder, error) {
	suborder := make([]*SubOrder, 0)

	for rows.Next() {
		s := &SubOrder{}
		rows.Scan(&s.ID, &s.OrderID, &s.ItemID, &s.Quantity)
		suborder = append(suborder, s)
	}

	return suborder, nil
}
