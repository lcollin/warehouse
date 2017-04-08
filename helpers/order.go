package helpers

import (
	"fmt"

	shipm "github.com/coldbrewcloud/go-shippo/models"
	"github.com/ghmeier/bloodlines/gateways"
	tcg "github.com/jakelong95/TownCenter/gateways"
	tcm "github.com/jakelong95/TownCenter/models"
	"github.com/lcollin/warehouse/models"

	"github.com/pborman/uuid"
)

const SELECT_ALL = "SELECT id, userID, subscriptionID, requestDate, shipDate, quantity, status, labelUrl "

type baseHelper struct {
	sql gateways.SQL
}

type OrderI interface {
	GetByID(string) (*models.Order, error)
	GetByUserID(string) ([]*models.Order, error)
	GetAll(int, int) ([]*models.Order, error)
	Insert(*models.Order) error
	Update(*models.Order) error
	SetStatus(id uuid.UUID, status models.OrderStatus) error
	Delete(string) error
	GetShippingLabel(id uuid.UUID, userID uuid.UUID, roasterID uuid.UUID) (string, error)
}

type Order struct {
	*baseHelper
	TC tcg.TownCenterI
}

func NewOrder(sql gateways.SQL, tc tcg.TownCenterI) *Order {
	return &Order{baseHelper: &baseHelper{sql: sql}, TC: tc}
}

func (i *Order) GetByID(id string) (*models.Order, error) {
	rows, err := i.sql.Select(SELECT_ALL+" FROM orderT WHERE id=?", id)
	if err != nil {
		return nil, err
	}

	items, err := models.OrderFromSQL(rows)
	if err != nil {
		return nil, err
	}

	if len(items) <= 0 {
		return nil, nil
	}

	return items[0], err
}

func (i *Order) GetByUserID(userID string) ([]*models.Order, error) {
	rows, err := i.sql.Select(SELECT_ALL+" FROM orderT WHERE userID=? ORDER BY status ASC, id ASC", userID)
	if err != nil {
		return nil, err
	}

	items, err := models.OrderFromSQL(rows)
	if err != nil {
		return nil, err
	}

	return items, err
}

func (i *Order) GetAll(offset int, limit int) ([]*models.Order, error) {
	rows, err := i.sql.Select(SELECT_ALL+" FROM orderT ORDER BY status ASC, id ASC LIMIT ?,?", offset, limit)
	if err != nil {
		return nil, err
	}

	items, err := models.OrderFromSQL(rows)
	if err != nil {
		return nil, err
	}

	return items, err
}

/* GetShippingLabel for an order with the given ID */
func (i *Order) GetShippingLabel(id uuid.UUID, userID uuid.UUID, roasterID uuid.UUID) (string, error) {
	user, err := i.TC.GetUser(userID)
	if err != nil {
		return "", err
	}
	roaster, err := i.TC.GetRoaster(roasterID)
	if err != nil {
		return "", err
	}
	order, err := i.GetByID(id.String())
	if err != nil {
		return "", err
	}
	if order == nil {
		return "", fmt.Errorf("No order found.")
	}

	if order.LabelURL != "" {
		return order.LabelURL, nil
	}

	var privateToken = "shippo_test_c235414aacd89a1597122e88e28476c624b8f106" //os.Getenv("PRIVATE_TOKEN")
	//create Shippo Client instance
	c := shippo.NewClient(privateToken)
	//create shipment using carrier account
	shipment := createShipmentUsingCarrierAccount(c, user, roaster)
	//choose and purchase shipping label
	label := purchasingShippingLabel(c, shipment)
	//extract url from transaction object
	url := label.LabelURL

	return url, nil

}

func (i *Order) Insert(order *models.Order) error {
	err := i.sql.Modify(
		"INSERT INTO orderT (id, userID, subscriptionID, requestDate, shipDate, quantity, status) VALUE (?,?,?,?,?,?,?)",
		order.ID,
		order.UserID,
		order.SubscriptionID,
		order.RequestDate,
		order.ShipDate,
		order.Quantity,
		order.Status,
	)

	return err
}

func (i *Order) Update(order *models.Order) error {
	err := i.sql.Modify(
		"UPDATE orderT SET userID=?, subscriptionID=?, requestDate=?, shipDate=?, quantity=?, status=?, labelUrl WHERE id=?",
		order.UserID,
		order.SubscriptionID,
		order.RequestDate,
		order.ShipDate,
		order.Quantity,
		order.Status,
		order.LabelURL,
		order.ID.String(),
	)

	return err
}

func (i *Order) SetStatus(id uuid.UUID, status models.OrderStatus) error {
	err := i.sql.Modify("UPDATE orderT SET status=? WHERE id=?", string(status), id.String())
	return err
}

func (i *Order) Delete(id string) error {
	err := i.sql.Modify("DELETE FROM orderT WHERE id=?", id)
	return err
}
