package models

import (
	"database/sql"
	"fmt"
	"strconv"

	query "github.com/ghmeier/bloodlines/gateways/sql"

	"github.com/pborman/uuid"
	"gopkg.in/gin-gonic/gin.v1"
)

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
	PhotoURL      string    `json:"photoUrl"`

	// // These can be utilized in a later version if desired
	// LeadTime      int `json: "lead_time"`
	// ReorderLevel  int `json: "reorder_level"`
	// PipelineStock int `json: "pipeline_stock"`
}

type itemSearch struct {
	*query.BaseSearch
	base string
	cost *query.SortTerm
	name *query.SortTerm
	id   *query.SortTerm
}

func NewItem(roasterID uuid.UUID, name string, pictureURL string, coffeeType string, inStockBags int, providerPrice float64, consumerPrice float64, ozInBag float64) *Item {
	return &Item{
		ID:            uuid.NewUUID(),
		RoasterID:     roasterID,
		Name:          name,
		PictureURL:    pictureURL,
		CoffeeType:    coffeeType,
		InStockBags:   inStockBags,
		ProviderPrice: providerPrice,
		ConsumerPrice: consumerPrice,
		OzInBag:       ozInBag,
	}
}

func ItemFromSQL(rows *sql.Rows) ([]*Item, error) {
	item := make([]*Item, 0)

	for rows.Next() {
		s := &Item{}
		rows.Scan(&s.ID, &s.RoasterID, &s.Name, &s.PictureURL, &s.CoffeeType, &s.InStockBags, &s.ProviderPrice, &s.ConsumerPrice, &s.OzInBag, &s.PhotoURL)
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
		base:       "SELECT id, roasterID, name, pictureURL, coffeeType, inStockBags, providerPrice, consumerPrice, ozInBag, photoUrl FROM item",
		cost:       query.NewSortTerm("consumerPrice", q, cost, false),
		name:       query.NewSortTerm("name", q, name, true),
		id:         query.NewSortTerm("id", "", 1, false),
	}
}

func (i *itemSearch) ToQuery() string {
	query := i.Limit(i.Order(i.Where(i.base, i.cost, i.name), i.id, i.cost, i.name))
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
