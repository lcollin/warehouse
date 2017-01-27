package containers

import (
	"database/sql"
	"github.com/pborman/uuid"
)

type InventoryItem struct {
	ID            uuid.UUID `json: "id"`
	ShopID        uuid.UUID `json: "shop_id"`
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

func InventoryItemFromSql(rows *sql.Rows) ([]*InventoryItem, error) {
	inventoryItem := make([]*InventoryItem, 0)

	for rows.Next() {
		s := &InventoryItem{}
		rows.Scan(&s.ID, &s.ShopID, &s.Name, &s.Picture, &s.Type, &s.InStockBags, &s.ProviderPrice, &s.ConsumerPrice, &s.OzInBag)
		inventoryItem = append(inventoryItem, s)
	}

	return inventoryItem, nil
}
