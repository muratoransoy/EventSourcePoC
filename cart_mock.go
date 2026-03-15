package main

import (
	"EvenSourcePOC/cart"
	"EvenSourcePOC/product"
	"context"
	"fmt"
	"log"

	"github.com/google/uuid"
	"github.com/kurrent-io/KurrentDB-Client-Go/kurrentdb"
)

var (
	cartID   = uuid.New()
	streamID = fmt.Sprintf("shoppingCart-%s", cartID)

	macbookPro     = product.NewProduct(190.000, "Macbook PRO")
	Iphone16       = product.NewProduct(60.000, "Iphone 16")
	KGAntepFistiği = product.NewProduct(799.99, "1 KG Antep Fıstık")

	events = []interface{}{
		cart.NewCartAdded(macbookPro, 1),
		cart.NewCartAdded(macbookPro, 2),
		cart.NewCartAdded(Iphone16, 1),
		cart.NewCartAdded(KGAntepFistiği, 3),
		cart.NewCartRemoved(macbookPro, 1),
	}

	cartMockOptions = kurrentdb.AppendToStreamOptions{
		StreamState: kurrentdb.Any{},
	}
)

func SetStreamMockCustomerId(ShoppingCartId string) {
	streamID = fmt.Sprintf("shoppingCart-%s", ShoppingCartId)
}

func StreamMockCart() {
	for i, event := range events {

		var eventData []byte
		var EventType string

		switch e := event.(type) {
		case cart.CartAdded:
			eventData, _ = e.ToJSON()
			EventType = string(e.CartEventType())
		case cart.CartRemoved:
			eventData, _ = e.ToJSON()
			EventType = string(e.CartEventType())
		default:
			panic("Unknown event type")
		}

		_, err := db.AppendToStream(
			context.Background(),
			streamID,
			cartMockOptions,
			kurrentdb.EventData{
				ContentType: kurrentdb.ContentTypeJson,
				EventType:   EventType,
				Data:        eventData,
			},
		)

		if err != nil {
			log.Printf("Error appending event %d: %v", i+1, err)
		}

		fmt.Printf("Appended event %d to stream %s\n", i+1, streamID)
	}
}
