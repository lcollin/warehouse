package containers

import (
	"database/sql"
	"github.com/pborman/uuid"
)

type Item struct {
	ID            uuid.UUID `json: "id"`
	ShopID        uuid.UUID `json: "shop_id"`
	Name          string    `json: "name"`
	Picture       string    `json: "picture_url"`
	Type          string    `json: "type"`
	InStockBags   int       `json: "in_stock"`
	ProviderPrice float64   `json: "provider_price"`
	ConsumerPrice float64   `json: "consumer_price"`
	OzInBag       float64   `json: "oz_in_bag"`

	// // These can be utilized in a later version if desired
	// LeadTime      int `json: "lead_time"`
	// ReorderLevel  int `json: "reorder_level"`
	// PipelineStock int `json: "pipeline_stock"`
}

func NewItem(shop_id, name, picture_url, coffee_type, in_stock, provider_price, consumer_price, oz_in_bag string) *Item {
	return &Item{
		ID:             uuid.NewUUID(),
		ShopID:         shop_id,
		Name:           name,
		Picture:        picture_url,
		Type:						coffee_type,
		InStockBags:		in_stock,
		ProviderPrice:	provider_price,
		ConsumerPrice:	consumer_price,
		OzInBag:				oz_in_bag,
	}
}

func ItemFromSQL(rows *sql.Rows) ([]*Item, error) {
	item := make([]*Item, 0)

	for rows.Next() {
		s := &Item{}
		rows.Scan(&s.ID, &s.ShopID, &s.Name, &s.Picture, &s.Type, &s.InStockBags, &s.ProviderPrice, &s.ConsumerPrice, &s.OzInBag)
		item = append(item, s)
	}

	return item, nil
}
