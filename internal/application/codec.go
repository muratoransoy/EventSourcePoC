package application

import (
	"encoding/json"
	"fmt"

	"eventsourcepoc/internal/domain/cart"
	"eventsourcepoc/internal/domain/order"
)

func MarshalEvent(eventType string, payload any) (EventRecord, error) {
	data, err := json.Marshal(payload)
	if err != nil {
		return EventRecord{}, fmt.Errorf("marshal %s: %w", eventType, err)
	}

	return EventRecord{
		EventType: eventType,
		Data:      data,
	}, nil
}

func DecodePayload(record EventRecord) (any, error) {
	switch record.EventType {
	case cart.EventTypeCartCreated:
		var event cart.CartCreated
		return event, json.Unmarshal(record.Data, &event)
	case cart.EventTypeItemAddedToCart:
		var event cart.ItemAddedToCart
		return event, json.Unmarshal(record.Data, &event)
	case cart.EventTypeItemRemovedFromCart:
		var event cart.ItemRemovedFromCart
		return event, json.Unmarshal(record.Data, &event)
	case order.EventTypeOrderCreated:
		var event order.OrderCreated
		return event, json.Unmarshal(record.Data, &event)
	default:
		var generic map[string]any
		return generic, json.Unmarshal(record.Data, &generic)
	}
}
