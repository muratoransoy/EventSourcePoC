package cart

import "github.com/google/uuid"

type Cart struct {
	ID       uuid.UUID `json:"id"`
	UserID   uuid.UUID `json:"userId"`
	Currency Currency  `json:"currency"`
}

func NewCart(userID uuid.UUID, currency Currency) Cart {
	return Cart{
		ID:       uuid.New(),
		UserID:   userID,
		Currency: currency,
	}
}
