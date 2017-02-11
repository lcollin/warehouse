package models

import (
	"database/sql"
	"github.com/pborman/uuid"
)

type Item struct {
	ID            uuid.UUID `json: "id"`
	RoasterID     uuid.UUID `json: "roasterID"`
	Name          string    `json: "name"`
	Picture       string    `json: "pictureUrl"`
	Type          string    `json: "coffeeType"`
	InStockBags   int       `json: "inStock"`
	ProviderPrice float64   `json: "providerPrice"`
	ConsumerPrice float64   `json: "consumerPrice"`
	OzInBag       float64   `json: "ozInBag"`

	// // These can be utilized in a later version if desired
	// LeadTime      int `json: "lead_time"`
	// ReorderLevel  int `json: "reorder_level"`
	// PipelineStock int `json: "pipeline_stock"`
}

func NewItem(roasterID uuid.UUID, name string, pictureURL string, coffeeType string, inStock int, providerPrice float64, consumerPrice float64, ozInBag float64) *Item {
	return &Item{
		ID:            uuid.NewUUID(),
		RoasterID:     roasterID,
		Name:          name,
		Picture:       pictureURL,
		Type:          coffeeType,
		InStockBags:   inStock,
		ProviderPrice: providerPrice,
		ConsumerPrice: consumerPrice,
		OzInBag:       ozInBag,
	}
}

func ItemFromSQL(rows *sql.Rows) ([]*Item, error) {
	item := make([]*Item, 0)

	for rows.Next() {
		s := &Item{}
		rows.Scan(&s.ID, &s.RoasterID, &s.Name, &s.Picture, &s.Type, &s.InStockBags, &s.ProviderPrice, &s.ConsumerPrice, &s.OzInBag)
		item = append(item, s)
	}

	return item, nil
}
