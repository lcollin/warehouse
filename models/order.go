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
	Status         OrderStatus `json"status"`
	LabelURL       string      `json:"labelUrl"`
	ItemID         uuid.UUID   `json:"itemId"`
}

/*ShipmentRequest represents the data needed to create a shipping label using Shippo API*/
type ShipmentRequest struct {
	OrderID      uuid.UUID `json:"orderId"`
	UserID       uuid.UUID `json:"userId"`
	RoasterID    uuid.UUID `json:"roasterId"`
	Quantity     uint64    `json:"quantity"`
	OzInBag      float64   `json:"ozInBag"`
	Length       float64   `json:"length"`
	Width        float64   `json:"width"`
	Height       float64   `json:"height"`
	DistanceUnit string    `json:"distanceUnit"`
	MassUnit     string    `json:"massUnit"`
}

/*Dimensions represent the data needed to get the correct sized shipment*/
type Dimensions struct {
	Quantity     uint64  `json:"quantity" binding:"required"`
	OzInBag      float64 `json:"ozInBag" binding:"required"`
	Length       float64 `json:"length" binding:"required"`
	Width        float64 `json:"width" binding:"required"`
	Height       float64 `json:"height" binding:"required"`
	DistanceUnit string  `json:"distanceUnit" binding:"required"`
	MassUnit     string  `json:"massUnit" binding:"required"`
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

/*NewDimensions creates a new Dimensions struct with the specified size and weight*/
func NewDimensions(quantity uint64, ozInBag float64, length float64, width float64, height float64, distanceUnit string, massUnit string) *Dimensions {
	return &Dimensions{
		Quantity:     quantity,
		OzInBag:      ozInBag,
		Length:       length,
		Width:        width,
		Height:       height,
		DistanceUnit: distanceUnit,
		MassUnit:     massUnit,
	}
}

func OrderFromSQL(rows *sql.Rows) ([]*Order, error) {
	order := make([]*Order, 0)

	for rows.Next() {
		o := &Order{}
		var status string
		rows.Scan(&o.ID, &o.UserID, &o.SubscriptionID, &o.RequestDate, &o.ShipDate, &o.Quantity, &status, &o.LabelURL, &o.ItemID)
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

func toDistanceUnit(s string) (DistanceUnit, bool) {
	switch s {
	case CENTIMETER:
		return CENTIMETER, true
	case MILLIMETER:
		return MILLIMETER, true
	case METER:
		return METER, true
	case INCH:
		return INCH, true
	case FOOT:
		return FOOT, true
	case YARD:
		return YARD, true
	default:
		return "", false
	}
}

/*DistanceUnit is an enum wrapper for valid content type*/
type DistanceUnit string

/*valid DistanceUnits that Shippo accepts*/
const (
	CENTIMETER = "cm"
	MILLIMETER = "mm"
	METER      = "m"
	INCH       = "in"
	FOOT       = "ft"
	YARD       = "yd"
)

func toMassUnit(s string) (MassUnit, bool) {
	switch s {
	case GRAM:
		return GRAM, true
	case OUNCE:
		return OUNCE, true
	case POUND:
		return POUND, true
	case KILOGRAM:
		return KILOGRAM, true
	default:
		return "", false
	}
}

/*DistanceUnit is an enum wrapper for valid content type*/
type MassUnit string

/*valid MassUnits that Shippo accepts*/
const (
	GRAM     = "g"
	OUNCE    = "oz"
	POUND    = "lb"
	KILOGRAM = "kg"
)
