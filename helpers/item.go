package helpers

import (
	"github.com/ghmeier/bloodlines/gateways"
	"github.com/lcollin/warehouse/containers"
)

type ItemI interface {
	GetByID(string) (*containers.Item, error)
	GetByShopID(string) (*containers.Item, error)
	GetAll(int, int) ([]*containers.Item, error)
	Insert(*containers.Item) error
	Update(*containers.Item, string) error
	Delete(string) error
}

type Item struct {
	*baseHelper
}

func NewItem(sql gateways.SQL) *Item {
	return &Item{baseHelper: &baseHelper{sql: sql}}
}

func (i *Item) GetByID(id string) (*containers.Item, error) {
	rows, err := i.sql.Select("SELECT * FROM item WHERE id=?", id)
	if err != nil {
		return nil, err
	}

	items, err := containers.ItemFromSQL(rows)
	if err != nil {
		return nil, err
	}

	return items[0], err
}

func (i *Item) GetByShopID(shopID string) (*containers.Item, error) {
	rows, err := i.sql.Select("SELECT * FROM item WHERE shop_id=?", shopID)
	if err != nil {
		return nil, err
	}

	items, err := containers.ItemFromSQL(rows)
	if err != nil {
		return nil, err
	}

	return items[0], err
}

func (i *Item) GetAll(offset int, limit int) ([]*containers.Item, error) {
	rows, err := i.sql.Select("SELECT * FROM item ORDER BY id ASC LIMIT ?,?", offset, limit)
	if err != nil {
		return nil, err
	}

	items, err := containers.ItemFromSQL(rows)
	if err != nil {
		return nil, err
	}

	return items, err
}

func (i *Item) Insert(item *containers.Item) error {
	err := i.sql.Modify(
		"INSERT INTO item (id, shop_id, name, picture_url, type, in_stock, provider_price, consumer_price, oz_in_bag) VALUE (?,?,?,?,?,?,?,?,?)",
		item.ID,
		item.ShopID,
		item.Name,
		item.Picture,
		item.Type,
		item.InStockBags,
		item.ProviderPrice,
		item.ConsumerPrice,
		item.OzInBag,
	)

	return err
}

func (i *Item) Update(item *containers.Item, id string) error {
	err := i.sql.Modify(
		"UPDATE item SET shop_id=?, name=?, picture_url=?, type=?, in_stock=?, provider_price=?, consumer_price=?, oz_in_bag=? WHERE id=?",
		item.ShopID,
		item.Name,
		item.Picture,
		item.Type,
		item.InStockBags,
		item.ProviderPrice,
		item.ConsumerPrice,
		item.OzInBag,
		id,
	)

	return err
}

func (i *Item) Delete(id string) error {
	err := i.sql.Modify("DELETE FROM item WHERE id=?", id)
	return err
}
