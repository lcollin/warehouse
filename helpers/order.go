package helpers

import (
	"gopkg.in/alexcesaro/statsd.v2"

	"github.com/ghmeier/bloodlines/gateways"
	"github.com/lcollin/warehouse/containers"
)

type baseHelper struct {
	sql   gateways.SQL
	stats *statsd.Client
}

type InventoryI interface {
	GetByID(string) (*containers.Inventory, error)
	GetByShopID(string) (*containers.Inventory, error)
	GetAll(int, int) ([]*containers.Inventory, error)
	Insert(*containers.Inventory) error
	Update(*containers.Inventory, string) error	
	Delete(string) error
}

type Inventory struct {
	*baseHelper
}

func NewInventory(sql gateways.SQL) *Inventory {
	return &Inventory{baseHelper: &baseHelper{sql: sql}}
}

func (i *Inventory) GetByID(id string) (*containers.Inventory, error) {
	rows, err := i.sql.Select("SELECT * FROM inventory WHERE id=?", id)
	if err != nil {
		return nil, err
	}

	items, err := containers.InventoryFromSQL(rows)
	if err != nil {
		return nil, err
	}

	return items[0], err
}

func (i *Inventory) GetByShopID(shop_id string) (*containers.Inventory, error) {
	rows, err := i.sql.Select("SELECT * FROM inventory WHERE shop_id=?", shop_id)
	if err != nil {
		return nil, err
	}

	items, err := containers.InventoryFromSQL(rows)
	if err != nil {
		return nil, err
	}

	return items, err
}

func (i *Inventory) GetAll(offset int, limit int) ([]*containers.Inventory, error) {
	rows, err := i.sql.Select("SELECT * FROM inventory ORDER BY id ASC LIMIT ?,?", offset, limit)
	if err != nil {
		return nil, err
	}

	items, err := containers.InventoryFromSQL(rows)
	if err != nil {
		return nil, err
	}

	return items, err
}

func (i *Inventory) Insert(inventory *containers.Inventory) error {
	err := i.sql.Modify(
		"INSERT INTO inventory (id, shop_id, name, picture_url, type, in_stock, provider_price, consumer_price, oz_in_bag) VALUE (?,?,?,?,?,?,?,?,?)",
		inventory.ID,
		inventory.ShopId,
		inventory.Name,
		inventory.Picture,
		inventory.Type,
		inventory.InStockBags,
		inventory.ProviderPrice,
		inventory.ConsumerPrice,
		inventory.OzInBag,     
	)

	return err
}

func (i *Inventory) Update(inventory *containers.Inventory, id string) error {
	err := i.sql.Modify(
		"UPDATE inventory SET shop_id=?, name=?, picture_url=?, type=?, in_stock=?, provider_price=?, consumer_price=?, oz_in_bag=? WHERE id=?",
		inventory.ShopId,
		inventory.Name,
		inventory.Picture,
		inventory.Type,
		inventory.InStockBags,
		inventory.ProviderPrice,
		inventory.ConsumerPrice,
		inventory.OzInBag,  
		id,
	)

	return err
}

func (i *Inventory) Delete(id string) error {
	err := i.sql.Modify("DELETE FROM inventory WHERE id=?", id)
	return err
}