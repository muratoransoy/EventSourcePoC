package cart

import (
	"encoding/json"
	"fmt"

	"github.com/google/uuid"
)

type State struct {
	CartID   uuid.UUID
	UserID   uuid.UUID
	Currency Currency
	Items    map[uuid.UUID]Item
	Created  bool
}

func NewState() State {
	return State{
		Items: make(map[uuid.UUID]Item),
	}
}

func (s *State) ApplyEvent(eventType string, data []byte) error {
	switch eventType {
	case EventTypeCartCreated:
		var event CartCreated
		if err := json.Unmarshal(data, &event); err != nil {
			return fmt.Errorf("decode %s: %w", eventType, err)
		}
		s.CartID = event.CartID
		s.UserID = event.UserID
		s.Currency = event.Currency
		s.Created = true
		return nil
	case EventTypeItemAddedToCart:
		var event ItemAddedToCart
		if err := json.Unmarshal(data, &event); err != nil {
			return fmt.Errorf("decode %s: %w", eventType, err)
		}
		item := s.Items[event.ProductID]
		item.ProductID = event.ProductID
		item.ProductName = event.ProductName
		item.UnitPrice = event.UnitPrice
		item.Quantity += event.Quantity
		s.Items[event.ProductID] = item
		return nil
	case EventTypeItemRemovedFromCart:
		var event ItemRemovedFromCart
		if err := json.Unmarshal(data, &event); err != nil {
			return fmt.Errorf("decode %s: %w", eventType, err)
		}
		item, ok := s.Items[event.ProductID]
		if !ok {
			return fmt.Errorf("product %s does not exist in cart", event.ProductID)
		}
		item.Quantity -= event.Quantity
		if item.Quantity <= 0 {
			delete(s.Items, event.ProductID)
			return nil
		}
		s.Items[event.ProductID] = item
		return nil
	default:
		return fmt.Errorf("unsupported cart event type %q", eventType)
	}
}

func (s State) IsEmpty() bool {
	return len(s.Items) == 0
}

func (s State) TotalAmount() float64 {
	total := 0.0
	for _, item := range s.Items {
		total += item.UnitPrice * float64(item.Quantity)
	}
	return total
}

func (s State) SnapshotItems() []Item {
	items := make([]Item, 0, len(s.Items))
	for _, item := range s.Items {
		items = append(items, item)
	}
	return items
}
