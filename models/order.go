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
	Quantity       uint64      `json:"quantity"`
	Status         OrderStatus `json:"status"`
	LabelURL       string      `json:"labelUrl"`
	TrackingURL    string      `json:"trackingUrl"`
	TransactionID  string      `json:transactionId"`
}

func NewOrder(userID, subscriptionID uuid.UUID, quantity uint64) *Order {
	return &Order{
		ID:             uuid.NewUUID(),
		UserID:         userID,
		SubscriptionID: subscriptionID,
		RequestDate:    time.Now(),
		Quantity:       quantity,
		Status:         PENDING,
	}
}

/*Set status will set the order's status to the given status*/
func (o *Order) SetStatus(status string) error {
	s, ok := toOrderStatus(status)
	if !ok {
		return fmt.Errorf("Invalid status")
	}
	o.Status = s
	return nil
}

func OrderFromSQL(rows *sql.Rows) ([]*Order, error) {
	order := make([]*Order, 0)
	for rows.Next() {
		o := &Order{}
		var status string
		err := rows.Scan(&o.ID, &o.UserID, &o.SubscriptionID, &o.RequestDate, &o.ShipDate, &o.Quantity, &status, &o.LabelURL, &o.TrackingURL, &o.TransactionID)
		if err != nil {
			return nil, err
		}
		statusType, ok := toOrderStatus(status)
		if !ok {
			return nil, fmt.Errorf("Invalid Status")
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
	case MISSING:
		return MISSING, true
	case DELIVERED:
		return DELIVERED, true
	case SHIPPED:
		return SHIPPED, true
	case TRANSIT:
		return TRANSIT, true
	case FAILURE:
		return FAILURE, true
	case RETURNED:
		return RETURNED, true
	default:
		return "", false
	}
}

/*OrderStatus is an enum wrapper for valid content type*/
type OrderStatus string

/*valid OrderStatus*/
const (
	PENDING   = "PENDING" //equivalent to shippo's UNKNOWN
	MISSING   = "MISSING"
	DELIVERED = "DELIVERED"
	SHIPPED   = "SHIPPED"
	TRANSIT   = "TRANSIT"
	FAILURE   = "FAILURE"
	RETURNED  = "RETURNED"
)
