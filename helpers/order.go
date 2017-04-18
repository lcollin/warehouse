package helpers

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/ghmeier/bloodlines/gateways"
	bmodels "github.com/ghmeier/bloodlines/models"
	tcg "github.com/jakelong95/TownCenter/gateways"
	"github.com/lcollin/warehouse/models"
	"github.com/pborman/uuid"
	"github.com/yuderekyu/go-shippo"
)

const SELECT_ALL = "SELECT id, userID, subscriptionID, requestDate, shipDate, quantity, status, labelUrl"

type baseHelper struct {
	sql gateways.SQL
}

type OrderI interface {
	GetByID(string) (*models.Order, error)
	GetByUserID(uuid.UUID, int, int) ([]*models.Order, error)
	GetByRoasterID(uuid.UUID, int, int) ([]*models.Order, error)
	GetAll(int, int) ([]*models.Order, error)
	Insert(*models.Order) error
	Update(*models.Order) error
	SetStatus(id uuid.UUID, status models.OrderStatus) error
	Delete(string) error
	GetShipmentLabel(shipmentRequest *models.ShipmentRequest) (*models.Order, error)
}

type Order struct {
	*baseHelper
	TC tcg.TownCenterI
	B  gateways.Bloodlines
}

func NewOrder(sql gateways.SQL, tc tcg.TownCenterI, b gateways.Bloodlines) *Order {
	return &Order{
		baseHelper: &baseHelper{sql: sql},
		TC:         tc,
		B:          b,
	}
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
		return nil, errors.New("Order does not exist")
	}

	return items[0], err
}

func (i *Order) GetByUserID(id uuid.UUID, offset, limit int) ([]*models.Order, error) {
	rows, err := i.sql.Select(SELECT_ALL+" FROM orderT WHERE userID=? ORDER BY status ASC, id ASC LIMIT ?,?", id.String(), offset, limit)
	if err != nil {
		return nil, err
	}

	return i.getList(rows)
}

func (i *Order) GetByRoasterID(id uuid.UUID, offset, limit int) ([]*models.Order, error) {
	rows, err := i.sql.Select(
		"SELECT o.id, o.userID, o.subscriptionID, o.requestDate, o.shipDate, o.quantity, o.status, o.labelUrl FROM orderT o "+
			"INNER JOIN covenant.subscription as s ON s.id=o.subscriptionId AND s.roasterId='?' "+
			"ORDER BY o.status ASC, id ASC LIMIT ?,?",
		id.String(),
		offset,
		limit)
	if err != nil {
		return nil, err
	}

	return i.getList(rows)
}

func (i *Order) GetAll(offset int, limit int) ([]*models.Order, error) {
	rows, err := i.sql.Select(SELECT_ALL+" FROM orderT ORDER BY status ASC, id ASC LIMIT ?,?", offset, limit)
	if err != nil {
		return nil, err
	}

	return i.getList(rows)
}

func (i *Order) getList(rows *sql.Rows) ([]*models.Order, error) {
	items, err := models.OrderFromSQL(rows)
	if err != nil {
		return nil, err
	}

	return items, nil
}

/*
GetShipmentLabel retrieves the specified order, user, and roaster information,
then creates a shippo shipment and transaction object and updates the order's labelURL
*/
func (i *Order) GetShipmentLabel(shipmentRequest *models.ShipmentRequest) (*models.Order, error) {
	order, err := i.GetByID(shipmentRequest.OrderID.String())
	if err != nil {
		return nil, err
	}
	user, err := i.TC.GetUser(shipmentRequest.UserID)
	if err != nil {
		return nil, err
	}
	roaster, err := i.TC.GetRoaster(shipmentRequest.RoasterID)
	if err != nil {
		return nil, err
	}
	dimensions, err := models.NewDimensions(shipmentRequest.Quantity, shipmentRequest.OzInBag, shipmentRequest.Length,
		shipmentRequest.Width, shipmentRequest.Height, shipmentRequest.DistanceUnit, shipmentRequest.MassUnit)
	if err != nil {
		return nil, err
	}
	//TODO: insert token within config
	var privateToken = "shippo_test_c235414aacd89a1597122e88e28476c624b8f106" //os.Getenv("PRIVATE_TOKEN")
	c := shippo.NewClient(privateToken)

	shipment, err := CreateShipment(c, user, roaster, dimensions)
	if err != nil {
		return nil, err
	}

	transaction, err := PurchaseShippingLabel(c, shipment)
	if err != nil {
		return nil, err
	}

	order.SetLabelURL(transaction.LabelURL)
	err = i.SetLabelURL(order.ID, transaction.LabelURL)
	if err != nil {
		return nil, err
	}
	return order, nil

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
		string(order.Status),
	)

	if err != nil {
		return err
	}

	user, err := i.TC.GetUser(order.UserID)
	if err != nil || user == nil {
		return err
	}

	addr := fmt.Sprintf(
		"%s\n%s\n%s, %s, %s, %s",
		user.AddressLine1,
		user.AddressLine2,
		user.AddressCity,
		user.AddressState,
		user.AddressCountry,
		user.AddressZip)

	i.B.ActivateTrigger("roaster_create_order", &bmodels.Receipt{
		UserID: user.ID,
		Values: map[string]string{
			"first_name":   user.FirstName,
			"last_name":    user.LastName,
			"quantity":     fmt.Sprintf("%d", order.Quantity),
			"request_date": order.RequestDate.String(),
			"address":      addr,
		},
	})
	return nil
}

func (i *Order) Update(order *models.Order) error {
	err := i.sql.Modify(
		"UPDATE orderT SET userID=?, subscriptionID=?, requestDate=?, shipDate=?, quantity=?, status=?, labelUrl=? WHERE id=?",
		order.UserID,
		order.SubscriptionID,
		order.RequestDate,
		order.ShipDate,
		order.Quantity,
		string(order.Status),
		order.LabelURL,
		order.ID.String(),
	)

	return err
}

func (i *Order) SetLabelURL(id uuid.UUID, labelURL string) error {
	err := i.sql.Modify("UPDATE orderT SET labelURL=? WHERE id=?", labelURL, id.String())
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
