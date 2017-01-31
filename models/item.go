package models

import (
	"database/sql"
	"github.com/pborman/uuid"
)

type Item struct {
	ID            uuid.UUID `json: "id"`
	ShopID        uuid.UUID `json: "shopId"`
	Name          string    `json: "name"`
	Picture       string    `json: "pictureUrl"`
	Type          string    `json: "type"`
	InStockBags   int       `json: "inStock"`
	ProviderPrice float64   `json: "providerPrice"`
	ConsumerPrice float64   `json: "consumerPrice"`
	OZInBag       float64   `json: "ozInBag"`

	// // These can be utilized in a later version if desired
	// LeadTime      int `json: "lead_time"`
	// ReorderLevel  int `json: "reorder_level"`
	// PipelineStock int `json: "pipeline_stock"`
}

func NewItem(shopID uuid.UUID, name, pictureURL, coffeeType string, inStock int, providerPrice, consumerPrice, ozInBag float64) *Item {
	return &Item{
		ID:            uuid.NewUUID(),
		ShopID:        shopID,
		Name:          name,
		Picture:       pictureURL,
		Type:          coffeeType,
		InStockBags:   inStock,
		ProviderPrice: providerPrice,
		ConsumerPrice: consumerPrice,
		OZInBag:       ozInBag,
	}
}

func ItemFromSQL(rows *sql.Rows) ([]*Item, error) {
	item := make([]*Item, 0)

	for rows.Next() {
		s := &Item{}
		rows.Scan(&s.ID, &s.ShopID, &s.Name, &s.Picture, &s.Type, &s.InStockBags, &s.ProviderPrice, &s.ConsumerPrice, &s.OZInBag)
		item = append(item, s)
	}

	return item, nil
}
