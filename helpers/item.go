package helpers

import (
	s "database/sql"
	"fmt"
	"mime/multipart"
	"strconv"
	"strings"
	"time"

	"github.com/pborman/uuid"

	"github.com/ghmeier/bloodlines/gateways"
	"github.com/ghmeier/bloodlines/gateways/sql"
	bmodels "github.com/ghmeier/bloodlines/models"
	c "github.com/ghmeier/coinage/gateways"
	cmodel "github.com/ghmeier/coinage/models"
	tcg "github.com/jakelong95/TownCenter/gateways"
	tcmodels "github.com/jakelong95/TownCenter/models"
	"github.com/lcollin/warehouse/models"
	cg "github.com/yuderekyu/covenant/gateways"
)

type ItemI interface {
	GetByID(string) (*models.Item, error)
	GetBySubscription(uuid.UUID) (*models.Item, error)
	GetByRoasterID(int, int, string) ([]*models.Item, error)
	GetAll(int, int, sql.Search) ([]*models.Item, error)
	GetAllInStock(int, int) ([]*models.Item, error)
	RemoveStock(*models.Item, int) error
	Insert(*models.Item) error
	Update(*models.Item) error
	Delete(string) error
	Upload(string, string, multipart.File) error
}

type Item struct {
	*baseHelper
	s3         gateways.S3
	Coinage    c.Coinage
	Covenant   cg.Covenant
	Bloodlines gateways.Bloodlines
	TownCenter tcg.TownCenterI
}

func NewItem(sql gateways.SQL, s3 gateways.S3, c c.Coinage, co cg.Covenant, b gateways.Bloodlines, t tcg.TownCenterI) *Item {
	return &Item{
		baseHelper: &baseHelper{sql: sql},
		s3:         s3,
		Bloodlines: b,
		Coinage:    c,
		Covenant:   co,
		TownCenter: t,
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

func (i *Item) GetBySubscription(id uuid.UUID) (*models.Item, error) {
	subscription, err := i.Covenant.GetSubscriptionById(id)
	if err != nil {
		return nil, err
	}

	return i.GetByID(subscription.ItemID.String())
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

	if err != nil {
		return err
	}

	request := &cmodel.PlanRequest{
		ItemID: item.ID,
	}
	_, err = i.Coinage.NewPlan(item.RoasterID, request)
	return err
}

func (i *Item) Update(item *models.Item) error {
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
		item.ID.String(),
	)

	return err
}

func (i *Item) RemoveStock(item *models.Item, quantity int) error {
	item.InStockBags -= quantity
	retErr := i.Update(item)

	if item.InStockBags > 5 {
		return retErr
	}

	user, roaster, err := i.owner(item)
	if err == nil || user == nil || roaster == nil {
		fmt.Println(err.Error())
		return retErr
	}

	if item.InStockBags < 0 {
		// send insufficient stock email
		i.Bloodlines.ActivateTrigger("no_stock", &bmodels.Receipt{
			UserID: user.ID,
			Values: map[string]string{
				"first_name":   user.FirstName,
				"roaster_name": roaster.Name,
				"quantity":     strconv.Itoa(quantity),
				"bean_name":    item.Name,
			},
		})
		return retErr
	}

	// send low stock email
	i.Bloodlines.ActivateTrigger("low_stock", &bmodels.Receipt{
		UserID: user.ID,
		Values: map[string]string{
			"first_name":   user.FirstName,
			"roaster_name": roaster.Name,
			"quantity":     strconv.Itoa(quantity),
			"bean_name":    item.Name,
			"remaining":    strconv.Itoa(item.InStockBags - quantity),
		},
	})
	return retErr

}

func (i *Item) owner(item *models.Item) (*tcmodels.User, *tcmodels.Roaster, error) {
	roaster, err := i.TownCenter.GetRoaster(item.RoasterID)
	if err != nil {
		return nil, nil, err
	}

	user, err := i.TownCenter.GetUserByRoaster(item.RoasterID)
	if err != nil {
		return nil, nil, err
	}

	return user, roaster, nil

}

func (i *Item) Delete(id string) error {
	err := i.sql.Modify("DELETE FROM item WHERE id=?", id)
	return err
}

func (i *Item) Upload(id string, name string, body multipart.File) error {
	filename := fmt.Sprintf("%s-%s", id, name)
	url, err := i.s3.Upload("profile", filename, body)
	if err != nil {
		return err
	}

	err = i.sql.Modify("UPDATE item SET pictureUrl=? WHERE id=?", url, id)
	return err
}
