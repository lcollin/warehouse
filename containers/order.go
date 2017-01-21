package containers

import (
	"database/sql"
	"github.com/pborman/uuid"
)

type Order struct {
	Id            uuid.UUID `json: "id"`
	ShopId 				uuid.UUID `json: "shop_id"`
	OrderContents
}

// OrderContents' ItemId(s) will correspond to item ids
// in the corresponding shop's inventory
type OrderContents struct {
	ItemId 				[]uuid.UUID `json: "item_id"`
	/**
	Alternatively, use a slice of OrderContents, with OrderContents containing:
		ItemID
		OzInBag
		Type
		Name
		...
	*/
}

func FromSql(rows *sql.Rows) ([]*Order, error) {
	order := make([]*Order, 0)

	for rows.Next() {
		s := &Order{}
		rows.Scan(&s.Name, &s.Picture, &s.Type, &s.InStockBags, &s.ProviderPrice, &s.ConsumerPrice, &s.OzInBag, &s.Id, &s.ShopId)
		order = append(order, s)
	}

	return order, nil
}
