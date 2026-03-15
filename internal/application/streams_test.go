package application

import (
	"testing"

	"github.com/google/uuid"
)

func TestCartStreamName(t *testing.T) {
	cartID := uuid.MustParse("11111111-1111-1111-1111-111111111111")
	if got := CartStreamName(cartID); got != "cart-11111111-1111-1111-1111-111111111111" {
		t.Fatalf("unexpected cart stream name: %s", got)
	}
}

func TestOrderStreamName(t *testing.T) {
	orderID := uuid.MustParse("22222222-2222-2222-2222-222222222222")
	if got := OrderStreamName(orderID); got != "order-22222222-2222-2222-2222-222222222222" {
		t.Fatalf("unexpected order stream name: %s", got)
	}
}
