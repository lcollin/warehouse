package containers

import (
	"database/sql"
	"github.com/pborman/uuid"
)

type Order struct {
	ID 	          uuid.UUID `json: "id"`
	UserID        uuid.UUID `json: "user_id"`
}

func NewSubOrder(userID string) *Order {
	return &Order{
		ID:             uuid.NewUUID(),
		UserID: 				userID,
	}
}

func OrderFromSql(rows *sql.Rows) ([]*Order, error) {
	order := make([]*Order, 0)

	for rows.Next() {
		o := &Order{}
		rows.Scan(&o.ID, &o.UserID)
		order = append(order, o)
	}

	return order, nil
}
