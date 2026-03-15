package cart

import (
	"encoding/json"
	"strings"
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestItemAddedToCartJSON(t *testing.T) {
	event := ItemAddedToCart{
		CartID:      uuid.MustParse("11111111-1111-1111-1111-111111111111"),
		ProductID:   uuid.MustParse("22222222-2222-2222-2222-222222222222"),
		ProductName: "Coffee",
		Quantity:    3,
		UnitPrice:   7.5,
		OccurredAt:  time.Date(2026, 3, 15, 12, 30, 0, 0, time.UTC),
	}

	data, err := json.Marshal(event)
	if err != nil {
		t.Fatalf("marshal event: %v", err)
	}

	got := string(data)
	for _, expected := range []string{
		`"cartId":"11111111-1111-1111-1111-111111111111"`,
		`"productId":"22222222-2222-2222-2222-222222222222"`,
		`"productName":"Coffee"`,
		`"quantity":3`,
		`"unitPrice":7.5`,
		`"occurredAt":"2026-03-15T12:30:00Z"`,
	} {
		if !strings.Contains(got, expected) {
			t.Fatalf("expected JSON to contain %s, got %s", expected, got)
		}
	}
}
