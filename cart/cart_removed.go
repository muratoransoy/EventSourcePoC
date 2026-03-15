package cart

import (
	"EvenSourcePOC/product"
	"encoding/json"

	"github.com/google/uuid"
)

type CartRemoved struct {
	ProductId   uuid.UUID `json:"productId"`
	Quantity    int       `json:"quantity"`
	ProductName string    `json:"product_name"`
}

func (*CartRemoved) CartEventType() CartEventType {
	return CartItemRemoved
}

func (item CartRemoved) ToJSON() ([]byte, error) {
	return json.Marshal(item)
}

func NewCartRemoved(product product.Product, quantity int) CartRemoved {
	return CartRemoved{
		ProductId:   product.ProductId,
		Quantity:    quantity,
		ProductName: product.ProductName,
	}
}
