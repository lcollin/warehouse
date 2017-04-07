package models

import (
	"database/sql"
	"fmt"
	"strconv"
	"time"

	query "github.com/ghmeier/bloodlines/gateways/sql"
	"github.com/ghmeier/bloodlines/models"

	"github.com/pborman/uuid"
	"gopkg.in/gin-gonic/gin.v1"
)

const SELECT_ALL = "SELECT id, roasterID, name, pictureURL, coffeeType, inStockBags, providerPrice, consumerPrice, ozInBag, description, isDecaf, isActive, tags, createdAt, updatedAt "

type Item struct {
	ID            uuid.UUID `json:"id"`
	RoasterID     uuid.UUID `json:"roasterId"`
	Name          string    `json:"name"`
	PictureURL    string    `json:"pictureURL"`
	CoffeeType    string    `json:"coffeeType"`
	InStockBags   int       `json:"inStockBags"`
	ProviderPrice float64   `json:"providerPrice"`
	ConsumerPrice float64   `json:"consumerPrice"`
	OzInBag       float64   `json:"ozInBag"`
	Description   string    `json:"description"`
	Decaf         bool      `json:"isDecaf"`
	Active        bool      `json:"isActive"`
	Tags          []string  `json:"tags"`
	CreatedAt     time.Time `json:"createdAt"`
	UpdatedAt     time.Time `json:"updatedAt"`

	// // These can be utilized in a later version if desired
	// LeadTime      int `json: "lead_time"`
	// ReorderLevel  int `json: "reorder_level"`
	// PipelineStock int `json: "pipeline_stock"`
}

type itemSearch struct {
	*query.BaseSearch
	base       string
	cost       *query.SortTerm
	name       *query.SortTerm
	id         *query.SortTerm
	coffeeType *query.SortTerm
}

func NewItem(roasterID uuid.UUID, name, coffeeType, description string, tags []string, inStockBags int, providerPrice, consumerPrice, ozInBag float64, decaf, active bool) *Item {
	return &Item{
		ID:            uuid.NewUUID(),
		RoasterID:     roasterID,
		Name:          name,
		CoffeeType:    coffeeType,
		InStockBags:   inStockBags,
		ProviderPrice: providerPrice,
		ConsumerPrice: consumerPrice,
		OzInBag:       ozInBag,
		Description:   description,
		Decaf:         decaf,
		Active:        active,
		Tags:          tags,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}
}

func ItemFromSQL(rows *sql.Rows) ([]*Item, error) {
	item := make([]*Item, 0)

	for rows.Next() {
		s := &Item{}
		var tagList string
		rows.Scan(
			&s.ID,
			&s.RoasterID,
			&s.Name,
			&s.PictureURL,
			&s.CoffeeType,
			&s.InStockBags,
			&s.ProviderPrice,
			&s.ConsumerPrice,
			&s.OzInBag,
			&s.Description,
			&s.Decaf,
			&s.Active,
			&tagList,
			&s.CreatedAt,
			&s.UpdatedAt)
		s.Tags = models.ToList(tagList)
		item = append(item, s)
	}

	return item, nil
}

func ItemSearch(ctx *gin.Context) query.Search {
	cost := queryInt(ctx, "cost")
	name := queryInt(ctx, "name")
	q, _ := ctx.GetQuery("q")

	return &itemSearch{
		BaseSearch: &query.BaseSearch{},
		base:       SELECT_ALL + " FROM item",
		cost:       query.NewSortTerm("consumerPrice", "", cost, false),
		name:       query.NewSortTerm("name", q, name, true),
		coffeeType: query.NewSortTerm("coffeeType", q, name, true),
		id:         query.NewSortTerm("id", "", 1, false),
	}
}

func (i *itemSearch) ToQuery() string {
	query := i.Limit(i.Order(i.Where(i.base, i.coffeeType, i.name), i.cost, i.name, i.coffeeType, i.id))
	fmt.Println(query)
	fmt.Println(query)
	return query
}

func queryInt(ctx *gin.Context, key string) int {
	str, _ := ctx.GetQuery(key)
	i, err := strconv.Atoi(str)
	if err != nil {
		i = -1
	}
	return i
}
