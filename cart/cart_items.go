package cart

import "github.com/google/uuid"

type CartItem struct {
	ProductId uuid.UUID `json:"productId"`
	Quantity  int       `json:"quantity"`
	Price     float64   `json:"price"`
}

func NewCartItem(quantity int, price float64) CartItem {
	return CartItem{
		ProductId: uuid.New(),
		Quantity:  quantity,
		Price:     price,
	}
}
