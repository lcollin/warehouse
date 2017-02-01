package models

import (
	"time"
	"database/sql"
	"github.com/pborman/uuid"
)

type Order struct {
	ID     uuid.UUID `json: "id"`
	UserID uuid.UUID `json: "userID"`
	SubscriptionID uuid.UUID `json: "subscriptionID"`
	RequestDate time.Time `json: "requestDate"`
	ShipDate time.Time `json: "shipDate"`
}

func NewOrder(userID uuid.UUID) *Order {
	return &Order{
		ID:     uuid.NewUUID(),
		UserID: userID,
		SubscriptionID: subscriptionID,
		RequestDate: time.Now(),
		ShipDate: time.Date(2020, time.January, 1, 1, 0, 0, 0, time.UTC),
	}
}

func OrderFromSQL(rows *sql.Rows) ([]*Order, error) {
	order := make([]*Order, 0)

	for rows.Next() {
		o := &Order{}
		rows.Scan(&o.ID, &o.UserID, &o.SubscriptionID, &o.RequestDate, &o.ShipDate)
		order = append(order, o)
	}

	return order, nil
}
