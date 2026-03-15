package product

import "github.com/google/uuid"

type Product struct {
	ProductId   uuid.UUID `json:"productId"`
	ProductName string    `json:"product_name"`
	Price       float64   `json:"price"`
}

func NewProduct(price float64, name string) Product {
	return Product{
		ProductId:   uuid.New(),
		Price:       price,
		ProductName: name,
	}
}
