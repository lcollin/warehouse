package helpers

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/pborman/uuid"

	"github.com/ghmeier/bloodlines/gateways"
	bmodels "github.com/ghmeier/bloodlines/models"
	tcg "github.com/jakelong95/TownCenter/gateways"
	w "github.com/lcollin/warehouse/gateways"
	"github.com/lcollin/warehouse/models"
)

const SELECT_ALL = "SELECT id, userID, subscriptionID, requestDate, shipDate, quantity, status, labelUrl, trackingUrl, transactionId"

type baseHelper struct {
	sql gateways.SQL
}

type OrderI interface {
	GetByID(string) (*models.Order, error)
	GetByUserID(uuid.UUID, int, int) ([]*models.Order, error)
	GetByRoasterID(uuid.UUID, int, int) ([]*models.Order, error)
	GetByTransactionID(id string) (*models.Order, error)
	GetAll(int, int) ([]*models.Order, error)
	Insert(*models.Order) error
	Update(*models.Order) error
	SetStatus(uuid.UUID, models.OrderStatus) error
	Delete(string) error
	GetShipmentInfo(*models.Order, *models.Item, *models.ShipmentRequest) (*models.Order, error)
}

type Order struct {
	*baseHelper
	TC tcg.TownCenterI
	B  gateways.Bloodlines
	S  w.Shippo
}

func NewOrder(sql gateways.SQL, tc tcg.TownCenterI, b gateways.Bloodlines, s w.Shippo) *Order {
	return &Order{
		baseHelper: &baseHelper{sql: sql},
		TC:         tc,
		B:          b,
		S:          s,
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

func (i *Order) GetByTransactionID(id string) (*models.Order, error) {
	rows, err := i.sql.Select(SELECT_ALL+" FROM orderT WHERE transactionId=?", id)
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

func (i *Order) GetByRoasterID(id uuid.UUID, offset, limit int) ([]*models.Order, error) {
	rows, err := i.sql.Select(
		"SELECT o.id, o.userID, o.subscriptionID, o.requestDate, o.shipDate, o.quantity, o.status, o.labelUrl, o.trackingUrl, o.transactionId FROM orderT o "+
			"INNER JOIN covenant.subscription as s ON s.id=o.subscriptionId AND s.roasterId=? "+
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
func (i *Order) GetShipmentInfo(order *models.Order, item *models.Item, req *models.ShipmentRequest) (*models.Order, error) {
	user, err := i.TC.GetUser(order.UserID)
	if err != nil {
		return nil, err
	}

	roaster, err := i.TC.GetRoaster(item.RoasterID)
	if err != nil {
		return nil, err
	}

	dimensions, err := models.NewDimensions(order.Quantity, item.OzInBag, req)
	if err != nil {
		return nil, err
	}

	fmt.Println(user, roaster, dimensions)
	shipment, err := i.S.CreateShipment(user, roaster, dimensions)
	if err != nil {
		return nil, err
	}

	transaction, err := i.S.PurchaseShippingLabel(shipment)
	if err != nil {
		return nil, err
	}

	order.LabelURL = transaction.LabelURL
	order.TrackingURL = transaction.TrackingURLProvider
	order.TransactionID = transaction.ObjectID

	//insert urls into database
	err = i.Update(order)
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
			"request_date": order.RequestDate.Format("Mon May 2, 2017"),
			"address":      addr,
		},
	})
	return nil
}

func (i *Order) Update(order *models.Order) error {
	err := i.sql.Modify(
		"UPDATE orderT SET userID=?, subscriptionID=?, requestDate=?, shipDate=?, quantity=?, status=?, labelUrl=?, trackingUrl=?, transactionId=? WHERE id=?",
		order.UserID,
		order.SubscriptionID,
		order.RequestDate,
		order.ShipDate,
		order.Quantity,
		string(order.Status),
		order.LabelURL,
		order.TrackingURL,
		order.TransactionID,
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
