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
}

func NewOrder(userID, subscriptionID uuid.UUID, quantity uint64) *Order {
	return &Order{
		ID:             uuid.NewUUID(),
		UserID:         userID,
		SubscriptionID: subscriptionID,
		RequestDate:    time.Now(),
		Quantity:       quantity,
		Status:         WAITING,
	}
}

/*SetURL sets the label and tracking url for the specified order*/
func (o *Order) SetURL(labelURL string, trackingURL string) error {
	if labelURL == "" {
		return fmt.Errorf("Invalid labelURL")
	}
	if trackingURL == "" {
		return fmt.Errorf("Invalid trackingURL")
	}
	o.LabelURL = labelURL
	o.TrackingURL = trackingURL
	return nil
}

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
		err := rows.Scan(&o.ID, &o.UserID, &o.SubscriptionID, &o.RequestDate, &o.ShipDate, &o.Quantity, &status, &o.LabelURL, &o.TrackingURL)
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
	case WAITING:
		return WAITING, true
	case QUEUED:
		return QUEUED, true
	case SUCCESS:
		return SUCCESS, true
	case ERROR:
		return ERROR, true
	case REFUNDED:
		return REFUNDED, true
	case REFUNDPENDING:
		return REFUNDPENDING, true
	case REFUNDREJECTED:
		return REFUNDREJECTED, true
	default:
		return "", false
	}
}

/*OrderStatus is an enum wrapper for valid content type*/
type OrderStatus string

/*valid OrderStatus*/
const (
	WAITING        = "WAITING"
	QUEUED         = "QUEUED"
	SUCCESS        = "SUCCESS"
	ERROR          = "ERROR"
	REFUNDED       = "REFUNDED"
	REFUNDPENDING  = "REFUNDPENDING"
	REFUNDREJECTED = "REFUNDREJECTED"
)
