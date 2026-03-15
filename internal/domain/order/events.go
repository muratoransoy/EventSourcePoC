package order

import (
	"time"

	"eventsourcepoc/internal/domain/cart"

	"github.com/google/uuid"
)

const EventTypeOrderCreated = "OrderCreated"

type OrderCreated struct {
	OrderID     uuid.UUID     `json:"orderId"`
	CartID      uuid.UUID     `json:"cartId"`
	UserID      uuid.UUID     `json:"userId"`
	Currency    cart.Currency `json:"currency"`
	Items       []cart.Item   `json:"items"`
	TotalAmount float64       `json:"totalAmount"`
	OccurredAt  time.Time     `json:"occurredAt"`
}

func (OrderCreated) EventType() string {
	return EventTypeOrderCreated
}
