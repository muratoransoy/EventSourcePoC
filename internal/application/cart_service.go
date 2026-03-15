package application

import (
	"context"
	"fmt"
	"time"

	"eventsourcepoc/internal/domain/cart"

	"github.com/google/uuid"
)

type Clock func() time.Time

type CartService struct {
	store EventStore
	clock Clock
}

func NewCartService(store EventStore) *CartService {
	return &CartService{
		store: store,
		clock: time.Now().UTC,
	}
}

func (s *CartService) CreateCart(ctx context.Context, cartID, userID uuid.UUID, currency cart.Currency) error {
	event := cart.CartCreated{
		CartID:     cartID,
		UserID:     userID,
		Currency:   currency,
		OccurredAt: s.clock(),
	}

	record, err := MarshalEvent(event.EventType(), event)
	if err != nil {
		return err
	}

	return s.store.Append(ctx, CartStreamName(cartID), ExpectNew, record)
}

func (s *CartService) AddItem(ctx context.Context, cartID, productID uuid.UUID, productName string, quantity int, unitPrice float64) error {
	if quantity <= 0 {
		return fmt.Errorf("quantity must be greater than zero")
	}
	if unitPrice < 0 {
		return fmt.Errorf("unit price must be non-negative")
	}

	event := cart.ItemAddedToCart{
		CartID:      cartID,
		ProductID:   productID,
		ProductName: productName,
		Quantity:    quantity,
		UnitPrice:   unitPrice,
		OccurredAt:  s.clock(),
	}

	record, err := MarshalEvent(event.EventType(), event)
	if err != nil {
		return err
	}

	return s.store.Append(ctx, CartStreamName(cartID), ExpectExisting, record)
}

func (s *CartService) RemoveItem(ctx context.Context, cartID, productID uuid.UUID, productName string, quantity int) error {
	if quantity <= 0 {
		return fmt.Errorf("quantity must be greater than zero")
	}

	state, err := s.LoadCart(ctx, cartID)
	if err != nil {
		return err
	}

	item, ok := state.Items[productID]
	if !ok {
		return fmt.Errorf("product %s is not in cart %s", productID, cartID)
	}
	if item.Quantity < quantity {
		return fmt.Errorf("cannot remove %d items; cart only contains %d", quantity, item.Quantity)
	}

	event := cart.ItemRemovedFromCart{
		CartID:      cartID,
		ProductID:   productID,
		ProductName: productName,
		Quantity:    quantity,
		OccurredAt:  s.clock(),
	}

	record, err := MarshalEvent(event.EventType(), event)
	if err != nil {
		return err
	}

	return s.store.Append(ctx, CartStreamName(cartID), ExpectExisting, record)
}

func (s *CartService) LoadCart(ctx context.Context, cartID uuid.UUID) (cart.State, error) {
	records, err := s.store.ReadStream(ctx, CartStreamName(cartID), 4_096)
	if err != nil {
		return cart.State{}, err
	}

	state := cart.NewState()
	for _, record := range records {
		if err := state.ApplyEvent(record.EventType, record.Data); err != nil {
			return cart.State{}, err
		}
	}

	if !state.Created {
		return cart.State{}, fmt.Errorf("cart %s does not exist", cartID)
	}

	return state, nil
}
