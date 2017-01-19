package containers

import(
	"database/sql"

	"github.com/pborman/uuid"
)

type Subscription struct {
	Name string `json: "name"`
	Picture string `json: "picture_url"`
	Type string `json: "type"`
	InStockBags int `json: "in_stock"`
	ProviderPrice float64 `json: "provider_price"`
	ConsumerPrice float64 `json: "consumer_price"`
	OzInBag float64 `json: "oz_in_bag"`
	Id uuid.UUID `json: "id"`
	ShopId uuid.UUID `json: "shop_id"`

	// These can be utilized in a later version if desired
	LeadTime int `json: "lead_time"`
	ReorderLevel int `json: "reorder_level"`
	PipelineStock int `json: "pipeline_stock"`
}

// type Inventory struct {
// 	Id uuid.UUID `json: "id"`
// 	OrderId int `json: "order_id"`
// 	Type string `json:"type"` 
// 	UserId int `json: "user_id"`
// 	Status string `json:"status"`
// 	CreatedAt string `json:"created_at"` //change to time.Date
// 	StartAt string `json:"start_at"` //change to time.Date
// 	TotalPrice string `json:"total_price"`
// }

// func FromSql(rows *sql.Rows) ([]*Subscription, error) {
// 	subscription := make([]*Subscription,0)

// 	for rows.Next() {
// 		s := &Subscription{}
// 		rows.Scan(&s.Id, &s.OrderId, &s.Type, &s.UserId, &s.Status, &s.CreatedAt, &s.StartAt, &s.TotalPrice)
// 		subscription = append(subscription, s)
// 	}

// 	return subscription, nil
// }

func FromSql(rows *sql.Rows) ([]*Inventory, error) {
	inventory := make([]*Inventory,0)

	for rows.Next() {
		s := &Inventory{}
		rows.Scan(&s.Name, &s.Picture, &s.Type, &s.InStockBags, &s.ProviderPrice, &s.ConsumerPrice, &s.OzInBag, &s.Id, &s.ShopId)
		inventory = append(inventory, s)
	}

	return inventory, nil
}