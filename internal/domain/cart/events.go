package cart

import (
	"errors"
	"strings"
	"time"

	"github.com/google/uuid"
)

const (
	EventTypeCartCreated         = "CartCreated"
	EventTypeItemAddedToCart     = "ItemAddedToCart"
	EventTypeItemRemovedFromCart = "ItemRemovedFromCart"
)

type Currency string

const (
	CurrencyTRY Currency = "TRY"
	CurrencyUSD Currency = "USD"
	CurrencyEUR Currency = "EUR"
)

func ParseCurrency(value string) (Currency, error) {
	normalized := Currency(strings.ToUpper(strings.TrimSpace(value)))
	switch normalized {
	case CurrencyTRY, CurrencyUSD, CurrencyEUR:
		return normalized, nil
	default:
		return "", errors.New("currency must be one of TRY, USD, EUR")
	}
}

type Item struct {
	ProductID   uuid.UUID `json:"productId"`
	ProductName string    `json:"productName"`
	Quantity    int       `json:"quantity"`
	UnitPrice   float64   `json:"unitPrice"`
}

type CartCreated struct {
	CartID     uuid.UUID `json:"cartId"`
	UserID     uuid.UUID `json:"userId"`
	Currency   Currency  `json:"currency"`
	OccurredAt time.Time `json:"occurredAt"`
}

func (CartCreated) EventType() string {
	return EventTypeCartCreated
}

type ItemAddedToCart struct {
	CartID      uuid.UUID `json:"cartId"`
	ProductID   uuid.UUID `json:"productId"`
	ProductName string    `json:"productName"`
	Quantity    int       `json:"quantity"`
	UnitPrice   float64   `json:"unitPrice"`
	OccurredAt  time.Time `json:"occurredAt"`
}

func (ItemAddedToCart) EventType() string {
	return EventTypeItemAddedToCart
}

type ItemRemovedFromCart struct {
	CartID      uuid.UUID `json:"cartId"`
	ProductID   uuid.UUID `json:"productId"`
	ProductName string    `json:"productName"`
	Quantity    int       `json:"quantity"`
	OccurredAt  time.Time `json:"occurredAt"`
}

func (ItemRemovedFromCart) EventType() string {
	return EventTypeItemRemovedFromCart
}
