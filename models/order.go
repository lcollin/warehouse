package models

import (
	"database/sql"
	"fmt"
	"github.com/pborman/uuid"
	"time"
)

type Order struct {
	ID             uuid.UUID   `json:"id"`
	UserID         uuid.UUID   `json:"userId"`
	SubscriptionID uuid.UUID   `json:"subscriptionId"`
	RequestDate    time.Time   `json:"requestDate"`
	ShipDate       time.Time   `json:"shipDate"`
	Quantity       int         `json:"quantity"`
	Status         OrderStatus `json"status"`
	LabelURL       string      `json:"labelUrl"`
}

type ShipmentRequest struct {
	ID        uuid.UUID `json:"id"`
	UserID    uuid.UUID `json:"userID"`
	RoasterID uuid.UUID `json:"roasterID"`
	Quantity  float64   `json:"quantity"`
	OzInBag   float64   `json:"ozInBag"`
}

func NewOrder(userID, subscriptionID uuid.UUID, quantity int) *Order {
	return &Order{
		ID:             uuid.NewUUID(),
		UserID:         userID,
		SubscriptionID: subscriptionID,
		RequestDate:    time.Now(),
		Quantity:       quantity,
		Status:         PENDING,
	}
}

func OrderFromSQL(rows *sql.Rows) ([]*Order, error) {
	order := make([]*Order, 0)

	for rows.Next() {
		o := &Order{}
		var status string
		rows.Scan(&o.ID, &o.UserID, &o.SubscriptionID, &o.RequestDate, &o.ShipDate, &o.Quantity, &status, &o.LabelURL)
		statusType, ok := toOrderStatus(status)
		if !ok {
			return nil, fmt.Errorf("Invalid Error Status")
		}
		o.Status = statusType
		order = append(order, o)
	}

	return order, nil
}

func toOrderStatus(s string) (OrderStatus, bool) {
	switch s {
	case PENDING:
		return PENDING, true
	case SHIPPED:
		return SHIPPED, true
	case ARRIVED:
		return ARRIVED, true
	case FINISHED:
		return FINISHED, true
	default:
		return "", false
	}
}

/*OrderStatus is an enum wrapper for valid content type*/
type OrderStatus string

/*valid OrderStatus*/
const (
	PENDING  = "PENDING"
	SHIPPED  = "SHIPPED"
	ARRIVED  = "ARRIVED"
	FINISHED = "FINISHED"
)
