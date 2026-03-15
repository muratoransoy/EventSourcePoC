package cart

import (
	"EvenSourcePOC/product"
	"encoding/json"

	"github.com/google/uuid"
)

type CartAdded struct {
	ProductId   uuid.UUID `json:"productId"`
	ProductName string    `json:"product_name"`
	Quantity    int       `json:"quantity"`
	Price       float64   `json:"price"`
}

func (*CartAdded) CartEventType() CartEventType {
	return CartItemAdded
}

func (item CartAdded) ToJSON() ([]byte, error) {
	return json.Marshal(item)
}

func NewCartAdded(product product.Product, quantity int) CartAdded {
	return CartAdded{
		ProductId:   product.ProductId,
		Quantity:    quantity,
		ProductName: product.ProductName,
		Price:       product.Price,
	}
}
