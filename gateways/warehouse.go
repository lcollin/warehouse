package gateways

import (
	"fmt"
	"net/http"

	"github.com/pborman/uuid"

	"github.com/ghmeier/bloodlines/config"
	g "github.com/ghmeier/bloodlines/gateways"
	"github.com/lcollin/warehouse/models"
)

/*Warehouse wraps all methods of the warehouse API*/
type Warehouse interface {
	GetAllItems(offset int, limit int) ([]*models.Item, error)
	NewItem(newItem *models.Item) (*models.Item, error)
	GetItemByID(id uuid.UUID) (*models.Item, error)
	UpdateItem(update *models.Item) (*models.Item, error)
	DeleteItem(id uuid.UUID) error
	GetAllOrders(offset int, limit int) ([]*models.Order, error)
	NewOrder(newOrder *models.Order) (*models.Order, error)
	GetOrderByID(id uuid.UUID) (*models.Order, error)
	GetOrderLabel(id uuid.UUID) (string, error)
	UpdateOrder(update *models.Order) (*models.Order, error)
	DeleteOrder(id uuid.UUID) error
	GetAllSubOrders(offset int, limit int) ([]*models.SubOrder, error)
	NewSubOrder(newSubOrder *models.SubOrder) (*models.SubOrder, error)
	GetSubOrderByID(id uuid.UUID) (*models.SubOrder, error)
	UpdateSubOrder(update *models.SubOrder) (*models.SubOrder, error)
	DeleteSubOrder(id uuid.UUID) error
}

type warehouse struct {
	*g.BaseService
	url    string
	client *http.Client
}

func NewWarehouse(config config.Warehouse) Warehouse {
	var url string
	if config.Port != "" {
		url = fmt.Sprintf("http://%s:%s/api/", config.Host, config.Port)
	} else {
		url = fmt.Sprintf("https://%s/api/", config.Host)
	}

	return &warehouse{
		BaseService: g.NewBaseService(),
		url:         url,
	}
}

func (w *warehouse) NewItem(newItem *models.Item) (*models.Item, error) {
	url := fmt.Sprintf("%sitem", w.url)

	var item models.Item
	err := w.ServiceSend(http.MethodPost, url, newItem, &item)
	if err != nil {
		return nil, err
	}

	return &item, nil
}

func (w *warehouse) GetAllItems(offset int, limit int) ([]*models.Item, error) {
	url := fmt.Sprintf("%sitem?offset=%d&limit=%d", w.url, offset, limit)

	item := make([]*models.Item, 0)
	err := w.ServiceSend(http.MethodGet, url, nil, &item)
	if err != nil {
		return nil, err
	}

	return item, nil
}

func (w *warehouse) GetItemByID(id uuid.UUID) (*models.Item, error) {
	url := fmt.Sprintf("%sitem/%s", w.url, id.String())

	var item models.Item
	err := w.ServiceSend(http.MethodGet, url, nil, &item)
	if err != nil {
		return nil, err
	}

	return &item, nil
}

func (w *warehouse) UpdateItem(update *models.Item) (*models.Item, error) {
	url := fmt.Sprintf("%sitem/%s", w.url, update.ID.String())

	var item models.Item
	err := w.ServiceSend(http.MethodPut, url, update, &item)
	if err != nil {
		return nil, err
	}

	return &item, nil
}

func (w *warehouse) DeleteItem(id uuid.UUID) error {
	url := fmt.Sprintf("%sitem/%s", w.url, id.String())

	err := w.ServiceSend(http.MethodDelete, url, nil, nil)
	if err != nil {
		return err
	}

	return nil
}

func (w *warehouse) NewOrder(newOrder *models.Order) (*models.Order, error) {
	url := fmt.Sprintf("%sorder", w.url)

	var order models.Order
	err := w.ServiceSend(http.MethodPost, url, newOrder, &order)
	if err != nil {
		return nil, err
	}

	return &order, nil
}

func (w *warehouse) GetAllOrders(offset int, limit int) ([]*models.Order, error) {
	url := fmt.Sprintf("%sorder?offset=%d&limit=%d", w.url, offset, limit)

	order := make([]*models.Order, 0)
	err := w.ServiceSend(http.MethodGet, url, nil, &order)
	if err != nil {
		return nil, err
	}

	return order, nil
}

func (w *warehouse) GetOrderByID(id uuid.UUID) (*models.Order, error) {
	url := fmt.Sprintf("%sorder/%s", w.url, id.String())

	var order models.Order
	err := w.ServiceSend(http.MethodGet, url, nil, &order)
	if err != nil {
		return nil, err
	}

	return &order, nil
}

func (w *warehouse) GetShippingLabel(newShipmentRequest *models.ShipmentRequest) (string, error) {
	url := fmt.Sprintf("%slabel", w.url)

	var shipmentRequest models.ShipmentRequest
	err := w.ServiceSend(http.MethodPost, url, newShipmentRequest, &shipmentRequest)
	if err != nil {
		return "", err
	}

	return label, nil
}

func (w *warehouse) UpdateOrder(update *models.Order) (*models.Order, error) {
	url := fmt.Sprintf("%sorder/%s", w.url, update.ID.String())

	var order models.Order
	err := w.ServiceSend(http.MethodPut, url, update, &order)
	if err != nil {
		return nil, err
	}

	return &order, nil
}

func (w *warehouse) DeleteOrder(id uuid.UUID) error {
	url := fmt.Sprintf("%sorder/%s", w.url, id.String())

	err := w.ServiceSend(http.MethodDelete, url, nil, nil)
	if err != nil {
		return err
	}

	return nil
}

func (w *warehouse) NewSubOrder(newSubOrder *models.SubOrder) (*models.SubOrder, error) {
	url := fmt.Sprintf("%ssuborder", w.url)

	var suborder models.SubOrder
	err := w.ServiceSend(http.MethodPost, url, newSubOrder, &suborder)
	if err != nil {
		return nil, err
	}

	return &suborder, nil
}

func (w *warehouse) GetAllSubOrders(offset int, limit int) ([]*models.SubOrder, error) {
	url := fmt.Sprintf("%ssuborder?offset=%d&limit=%d", w.url, offset, limit)

	suborder := make([]*models.SubOrder, 0)
	err := w.ServiceSend(http.MethodGet, url, nil, &suborder)
	if err != nil {
		return nil, err
	}

	return suborder, nil
}

func (w *warehouse) GetSubOrderByID(id uuid.UUID) (*models.SubOrder, error) {
	url := fmt.Sprintf("%ssuborder/%s", w.url, id.String())

	var suborder models.SubOrder
	err := w.ServiceSend(http.MethodGet, url, nil, &suborder)
	if err != nil {
		return nil, err
	}

	return &suborder, nil
}

func (w *warehouse) UpdateSubOrder(update *models.SubOrder) (*models.SubOrder, error) {
	url := fmt.Sprintf("%ssuborder/%s", w.url, update.ID.String())

	var suborder models.SubOrder
	err := w.ServiceSend(http.MethodPut, url, update, &suborder)
	if err != nil {
		return nil, err
	}

	return &suborder, nil
}

func (w *warehouse) DeleteSubOrder(id uuid.UUID) error {
	url := fmt.Sprintf("%ssuborder/%s", w.url, id.String())

	err := w.ServiceSend(http.MethodDelete, url, nil, nil)
	if err != nil {
		return err
	}

	return nil
}
