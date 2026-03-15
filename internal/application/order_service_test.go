package application

import (
	"context"
	"testing"
	"time"

	"eventsourcepoc/internal/domain/cart"
	"eventsourcepoc/internal/domain/order"

	"github.com/google/uuid"
)

func TestCreateOrderFromCartProjectsCurrentCartState(t *testing.T) {
	store := &memoryEventStore{}
	cartService := NewCartService(store)
	orderService := NewOrderService(store, cartService)

	fixedTime := func() time.Time { return time.Date(2026, 3, 15, 12, 0, 0, 0, time.UTC) }
	cartService.clock = fixedTime
	orderService.clock = fixedTime

	cartID := uuid.New()
	userID := uuid.New()
	productID := uuid.New()
	orderID := uuid.New()

	if err := cartService.CreateCart(context.Background(), cartID, userID, cart.CurrencyEUR); err != nil {
		t.Fatalf("create cart: %v", err)
	}
	if err := cartService.AddItem(context.Background(), cartID, productID, "Book", 2, 12.5); err != nil {
		t.Fatalf("add item: %v", err)
	}

	event, err := orderService.CreateOrderFromCart(context.Background(), orderID, cartID)
	if err != nil {
		t.Fatalf("create order: %v", err)
	}

	if event.EventType() != order.EventTypeOrderCreated {
		t.Fatalf("unexpected event type: %s", event.EventType())
	}
	if event.TotalAmount != 25 {
		t.Fatalf("unexpected total amount: %.2f", event.TotalAmount)
	}
	if len(event.Items) != 1 || event.Items[0].Quantity != 2 {
		t.Fatalf("unexpected order items: %+v", event.Items)
	}
}

type memoryEventStore struct {
	streams map[string][]EventRecord
}

func (m *memoryEventStore) Append(_ context.Context, streamID string, _ StreamExpectation, events ...EventRecord) error {
	if m.streams == nil {
		m.streams = make(map[string][]EventRecord)
	}

	for _, event := range events {
		event.StreamID = streamID
		event.EventNumber = uint64(len(m.streams[streamID]))
		m.streams[streamID] = append(m.streams[streamID], event)
	}
	return nil
}

func (m *memoryEventStore) ReadStream(_ context.Context, streamID string, _ uint64) ([]EventRecord, error) {
	return append([]EventRecord(nil), m.streams[streamID]...), nil
}

func (m *memoryEventStore) Subscribe(_ context.Context, _ string, _ bool, _ func(EventRecord) error) error {
	return nil
}
