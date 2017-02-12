package helpers

import (
	"github.com/ghmeier/bloodlines/gateways"
	"github.com/lcollin/warehouse/models"
)

type ItemI interface {
	GetByID(string) (*models.Item, error)
	GetByRoasterID(string) (*models.Item, error)
	GetAll(int, int) ([]*models.Item, error)
	GetAllInStock(int, int) ([]*models.Item, error)
	Insert(*models.Item) error
	Update(*models.Item, string) error
	Delete(string) error
}

type Item struct {
	*baseHelper
}

func NewItem(sql gateways.SQL) *Item {
	return &Item{baseHelper: &baseHelper{sql: sql}}
}

func (i *Item) GetByID(id string) (*models.Item, error) {
	rows, err := i.sql.Select("SELECT * FROM item WHERE id=?", id)
	if err != nil {
		return nil, err
	}

	items, err := models.ItemFromSQL(rows)
	if err != nil {
		return nil, err
	}

	return items[0], err
}

func (i *Item) GetByRoasterID(roasterID string) (*models.Item, error) {
	rows, err := i.sql.Select("SELECT * FROM item WHERE shopID=?", roasterID)
	if err != nil {
		return nil, err
	}

	items, err := models.ItemFromSQL(rows)
	if err != nil {
		return nil, err
	}

	return items[0], err
}

func (i *Item) GetAll(offset int, limit int) ([]*models.Item, error) {
	rows, err := i.sql.Select("SELECT * FROM item ORDER BY id ASC LIMIT ?,?", offset, limit)
	if err != nil {
		return nil, err
	}

	items, err := models.ItemFromSQL(rows)
	if err != nil {
		return nil, err
	}

	return items, err
}

func (i *Item) GetAllInStock(offset int, limit int) ([]*models.Item, error) {
	rows, err := i.sql.Select("SELECT * FROM item WHERE inStockBags>0 ORDER BY id ASC LIMIT ?,?", offset, limit)
	if err != nil {
		return nil, err
	}

	items, err := models.ItemFromSQL(rows)
	if err != nil {
		return nil, err
	}

	return items, err
}

func (i *Item) Insert(item *models.Item) error {
	err := i.sql.Modify(
		"INSERT INTO item (id, roasterID, name, pictureURL, coffeeType, inStockBags, providerPrice, consumerPrice, ozInBag) VALUE (?,?,?,?,?,?,?,?,?)",
		item.ID,
		item.RoasterID,
		item.Name,
		item.PictureURL,
		item.CoffeeType,
		item.InStockBags,
		item.ProviderPrice,
		item.ConsumerPrice,
		item.OzInBag,
	)

	return err
}

func (i *Item) Update(item *models.Item, id string) error {
	err := i.sql.Modify(
		"UPDATE item SET roasterID=?, name=?, pictureURL=?, coffeeType=?, inStockBags=?, providerPrice=?, consumerPrice=?, ozInBag=? WHERE id=?",
		item.RoasterID,
		item.Name,
		item.PictureURL,
		item.CoffeeType,
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
