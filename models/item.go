package models

import (
	"database/sql"
	"github.com/pborman/uuid"
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
