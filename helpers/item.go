package helpers

import (
	s "database/sql"
	"fmt"
	"mime/multipart"
	"strings"
	"time"

	"github.com/ghmeier/bloodlines/gateways"
	"github.com/ghmeier/bloodlines/gateways/sql"
	"github.com/lcollin/warehouse/models"
)

type ItemI interface {
	GetByID(string) (*models.Item, error)
	GetByRoasterID(int, int, string) ([]*models.Item, error)
	GetAll(int, int, sql.Search) ([]*models.Item, error)
	GetAllInStock(int, int) ([]*models.Item, error)
	Insert(*models.Item) error
	Update(*models.Item, string) error
	Delete(string) error
	Upload(string, string, multipart.File) error
}

type Item struct {
	*baseHelper
	s3 gateways.S3
}

func NewItem(sql gateways.SQL, s3 gateways.S3) *Item {
	return &Item{
		baseHelper: &baseHelper{sql: sql},
		s3:         s3,
	}
}

func (i *Item) GetByID(id string) (*models.Item, error) {
	rows, err := i.sql.Select(models.SELECT_ALL+" FROM item WHERE id=?", id)
	if err != nil {
		return nil, err
	}

	items, err := models.ItemFromSQL(rows)
	if err != nil {
		return nil, err
	}

	return items[0], err
}

func (i *Item) GetByRoasterID(offset, limit int, roasterID string) ([]*models.Item, error) {
	rows, err := i.sql.Select(
		models.SELECT_ALL+" FROM item WHERE roasterID=? ORDER BY id ASC LIMIT ?,?",
		roasterID,
		offset,
		limit)
	if err != nil {
		return nil, err
	}

	return i.handleItemsQuery(rows)
}

func (i *Item) GetAll(offset int, limit int, search sql.Search) ([]*models.Item, error) {
	rows, err := i.sql.Select(search.ToQuery(), offset, limit)
	if err != nil {
		return nil, err
	}

	return i.handleItemsQuery(rows)
}

func (i *Item) GetAllInStock(offset int, limit int) ([]*models.Item, error) {
	rows, err := i.sql.Select(models.SELECT_ALL+" FROM item WHERE inStockBags>0 ORDER BY id ASC LIMIT ?,?", offset, limit)
	if err != nil {
		return nil, err
	}

	return i.handleItemsQuery(rows)
}

func (i *Item) handleItemsQuery(rows *s.Rows) ([]*models.Item, error) {
	items, err := models.ItemFromSQL(rows)
	if err != nil {
		return nil, err
	}

	return items, err
}

func (i *Item) Insert(item *models.Item) error {
	err := i.sql.Modify(
		"INSERT INTO item (id, roasterID, name, pictureURL, coffeeType, inStockBags, providerPrice, consumerPrice, ozInBag, description, isDecaf, isActive, tags, updatedAt) VALUE (?,?,?,?,?,?,?,?,?,?,?,?,?,?)",
		item.ID,
		item.RoasterID,
		item.Name,
		item.PictureURL,
		item.CoffeeType,
		item.InStockBags,
		item.ProviderPrice,
		item.ConsumerPrice,
		item.OzInBag,
		item.Description,
		item.Decaf,
		item.Active,
		strings.Join(item.Tags, ","),
		time.Now(),
	)

	return err
}

func (i *Item) Update(item *models.Item, id string) error {
	err := i.sql.Modify(
		"UPDATE item SET roasterID=?, name=?, pictureURL=?, coffeeType=?, inStockBags=?, providerPrice=?, consumerPrice=?, ozInBag=?, description=?, isDecaf=?, isActive=?, tags=?, updatedAt=? WHERE id=?",
		item.RoasterID,
		item.Name,
		item.PictureURL,
		item.CoffeeType,
		item.InStockBags,
		item.ProviderPrice,
		item.ConsumerPrice,
		item.OzInBag,
		item.Description,
		item.Decaf,
		item.Active,
		strings.Join(item.Tags, ","),
		time.Now(),
		id,
	)

	return err
}

func (i *Item) Delete(id string) error {
	err := i.sql.Modify("DELETE FROM item WHERE id=?", id)
	return err
}

func (i *Item) Upload(id string, name string, body multipart.File) error {
	filename := fmt.Sprintf("%s-%s", id, name)
	url, err := i.s3.Upload("item-photo", filename, body)
	if err != nil {
		return err
	}

	err = i.sql.Modify("UPDATE item SET pictureUrl=? WHERE id=?", url, id)
	return err
}
