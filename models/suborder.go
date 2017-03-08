package models

import (
	"database/sql"
	"errors"
	"github.com/pborman/uuid"
)

type SubOrder struct {
	ID       uuid.UUID `json:"id"`
	OrderID  uuid.UUID `json:"orderId"`
	ItemID   uuid.UUID `json:"itemId"`
	Quantity int       `json:"quantity"`
}

func NewSubOrder(orderID uuid.UUID, itemID uuid.UUID, quantity int) *SubOrder {
	return &SubOrder{
		ID:       uuid.NewUUID(),
		OrderID:  orderID,
		ItemID:   itemID,
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

	if len(suborder) <= 0 {
		return nil, errors.New("no entries returned")
	}

	return suborder, nil
}
