package helpers

import (
	"errors"
	"github.com/coldbrewcloud/go-shippo"
	"github.com/ghmeier/bloodlines/gateways"
	tcg "github.com/jakelong95/TownCenter/gateways"
	"github.com/lcollin/warehouse/models"
	"github.com/pborman/uuid"
)

const SELECT_ALL = "SELECT id, userID, subscriptionID, requestDate, shipDate, quantity, status, labelUrl, itemId "

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
	GetShippingLabel(shipmentRequest *models.ShipmentRequest) (string, error)
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

/*GetShippingLabel for an order with the given ID */
func (i *Order) GetShippingLabel(shipmentRequest *models.ShipmentRequest) (string, error) {
	//check if order exists
	order, err := i.GetByID(shipmentRequest.OrderID.String())
	if err != nil {
		return "", err
	}
	//get user information
	user, err := i.TC.GetUser(shipmentRequest.UserID)
	if err != nil {
		return "", err
	}
	//get roaster information
	roaster, err := i.TC.GetRoaster(shipmentRequest.RoasterID)
	if err != nil {
		return "", err
	}

	dimensions := models.NewDimensions(shipmentRequest.Quantity, shipmentRequest.OzInBag, shipmentRequest.Length,
		shipmentRequest.Width, shipmentRequest.Height, shipmentRequest.DistanceUnit, shipmentRequest.MassUnit)

	//change this so shippo is defined within the gateway
	var privateToken = "shippo_test_c235414aacd89a1597122e88e28476c624b8f106" //os.Getenv("PRIVATE_TOKEN")
	//create Shippo Client instance
	c := shippo.NewClient(privateToken)
	//create shipment using carrier account
	shipment := CreateShipment(c, user, roaster, dimensions)
	//choose and purchase shipping label
	label := PurchaseShippingLabel(c, shipment)
	//extract url from transaction object
	url := label.LabelURL

	if url != "" {
		return "", errors.New("Shipping label failed to create")
	}

	// On success, insert url into database and return
	// var order models.Order
	order.LabelURL = url
	i.Insert(order)
	return url, nil

}

func (i *Order) Insert(order *models.Order) error {
	err := i.sql.Modify(
		"INSERT INTO orderT (id, userID, subscriptionID, requestDate, shipDate, quantity, status, labelUrl, itemId) VALUE (?,?,?,?,?,?,?,?,?)",
		order.ID,
		order.UserID,
		order.SubscriptionID,
		order.RequestDate,
		order.ShipDate,
		order.Quantity,
		order.Status,
		order.LabelURL,
		order.ItemID,
	)

	return err
}

func (i *Order) Update(order *models.Order) error {
	err := i.sql.Modify(
		"UPDATE orderT SET userID=?, subscriptionID=?, requestDate=?, shipDate=?, quantity=?, status=?, labelUrl=?, itemId=? WHERE id=?",
		order.UserID,
		order.SubscriptionID,
		order.RequestDate,
		order.ShipDate,
		order.Quantity,
		order.Status,
		order.LabelURL,
		order.ItemID,
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
