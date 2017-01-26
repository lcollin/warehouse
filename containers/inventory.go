package containers

import (
	"database/sql"

	"github.com/pborman/uuid"
)

type Inventory struct {
	ID            uuid.UUID `json: "id"`
	ShopId        uuid.UUID `json: "shop_id"`
	Name          string    `json: "name"`
	Picture       string    `json: "picture_url"`
	Type          string    `json: "type"`
	InStockBags   int       `json: "in_stock"`
	ProviderPrice float64   `json: "provider_price"`
	ConsumerPrice float64   `json: "consumer_price"`
	OzInBag       float64   `json: "oz_in_bag"`

	// These can be utilized in a later version if desired
	LeadTime      int `json: "lead_time"`
	ReorderLevel  int `json: "reorder_level"`
	PipelineStock int `json: "pipeline_stock"`
}

func FromSql(rows *sql.Rows) ([]*Inventory, error) {
	inventory := make([]*Inventory, 0)

	for rows.Next() {
		s := &Inventory{}
		rows.Scan(&s.Name, &s.Picture, &s.Type, &s.InStockBags, &s.ProviderPrice, &s.ConsumerPrice, &s.OzInBag, &s.Id, &s.ShopId)
		inventory = append(inventory, s)
	}

	return inventory, nil
}
