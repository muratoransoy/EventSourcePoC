package application

import (
	"context"
	"fmt"
	"time"

	"eventsourcepoc/internal/domain/order"

	"github.com/google/uuid"
)

type OrderService struct {
	store EventStore
	carts *CartService
	clock Clock
}

func NewOrderService(store EventStore, carts *CartService) *OrderService {
	return &OrderService{
		store: store,
		carts: carts,
		clock: time.Now().UTC,
	}
}

func (s *OrderService) CreateOrderFromCart(ctx context.Context, orderID, cartID uuid.UUID) (order.OrderCreated, error) {
	state, err := s.carts.LoadCart(ctx, cartID)
	if err != nil {
		return order.OrderCreated{}, err
	}
	if state.IsEmpty() {
		return order.OrderCreated{}, fmt.Errorf("cart %s has no items", cartID)
	}

	event := order.OrderCreated{
		OrderID:     orderID,
		CartID:      cartID,
		UserID:      state.UserID,
		Currency:    state.Currency,
		Items:       state.SnapshotItems(),
		TotalAmount: state.TotalAmount(),
		OccurredAt:  s.clock(),
	}

	record, err := MarshalEvent(event.EventType(), event)
	if err != nil {
		return order.OrderCreated{}, err
	}

	if err := s.store.Append(ctx, OrderStreamName(orderID), ExpectNew, record); err != nil {
		return order.OrderCreated{}, err
	}

	return event, nil
}
