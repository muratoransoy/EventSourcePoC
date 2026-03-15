package application

import (
	"fmt"

	"github.com/google/uuid"
)

func CartStreamName(cartID uuid.UUID) string {
	return fmt.Sprintf("cart-%s", cartID)
}

func OrderStreamName(orderID uuid.UUID) string {
	return fmt.Sprintf("order-%s", orderID)
}
