package orders

import "github.com/google/uuid"

type OrderCreated struct {
	OrderId     string      `json:"orderId"`
	CustomerId  string      `json:"customerId"`
	Items       []OrderItem `json:"items"`
	TotalAmount float64     `json:"totalAmount"`
	Status      string      `json:"status"`
}

func NewOrder() OrderCreated {
	return OrderCreated{
		OrderId:    uuid.NewString(),
		CustomerId: "customer456",
		Items: []OrderItem{
			{ProductId: "product789", Quantity: 2, Price: 19.99},
			{ProductId: "product012", Quantity: 1, Price: 9.99},
		},
		TotalAmount: 49.97,
		Status:      "Created",
	}
}
