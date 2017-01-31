package models

import (
	"database/sql"
	"github.com/pborman/uuid"
)

type SubOrder struct {
	ID      uuid.UUID `json: "id"`
	OrderID uuid.UUID `json: "order_id"`
	ItemID  uuid.UUID `json: "item_id"`
}

func NewSubOrder(orderID, itemID uuid.UUID) *SubOrder {
	return &SubOrder{
		ID:      uuid.NewUUID(),
		OrderID: orderID,
		ItemID:  itemID,
	}
}

func SubOrderFromSQL(rows *sql.Rows) ([]*SubOrder, error) {
	suborder := make([]*SubOrder, 0)

	for rows.Next() {
		s := &SubOrder{}
		rows.Scan(&s.ID, &s.OrderID, &s.ItemID)
		suborder = append(suborder, s)
	}

	return suborder, nil
}
