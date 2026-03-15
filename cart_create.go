package main

import (
	"EvenSourcePOC/cart"
	"context"
	"encoding/json"
	"log"

	"github.com/google/uuid"
	"github.com/kurrent-io/KurrentDB-Client-Go/kurrentdb"
)

var (
	streamType = "Cart"
	eventType  = "CartCreated"

	eventCartItemAddedType = "CartItemAdded"
)

func CreateCarts() {
	TLuserID := uuid.MustParse("427933cf-d8b7-4381-8924-a703ccf00eb6")
	euruserID := uuid.MustParse("f11282d7-05f8-4282-b14e-a8387abf5632")
	tlCart := cart.NewCart(TLuserID, cart.TL)
	eurCart := cart.NewCart(euruserID, cart.EUR)

	createCartEvent(tlCart)
	createCartEvent(eurCart)
}

func createCartEvent(_cart cart.Cart) {
	data, err := json.Marshal(_cart)

	if err != nil {
		panic(err)
	}

	eventData := kurrentdb.EventData{
		ContentType: kurrentdb.ContentTypeJson,
		EventType:   eventType,
		Data:        data,
	}

	options := kurrentdb.AppendToStreamOptions{
		StreamState: kurrentdb.NoStream{},
	}

	streamName := streamType + "-" + _cart.UserID.String()

	if _, err := db.AppendToStream(context.Background(), streamName, options, eventData); err != nil {
		panic(err)
	}

	log.Printf("Created cart event: %s", streamName)
}

func CreateProducts(stream string) {
	product1 := cart.NewCartItem(2, 100.23)
	product2 := cart.NewCartItem(1, 200.54)
	CreateCartItem(stream, product1)
	CreateCartItem(stream, product2)
}

func CreateCartItem(streamName string, item cart.CartItem) {
	data, err := json.Marshal(item)

	if err != nil {
		panic(err)
	}

	eventData := kurrentdb.EventData{
		ContentType: kurrentdb.ContentTypeJson,
		EventType:   eventCartItemAddedType,
		Data:        data,
	}

	options := kurrentdb.AppendToStreamOptions{
		StreamState: kurrentdb.StreamExists{},
	}

	if _, err := db.AppendToStream(context.Background(), streamName, options, eventData); err != nil {
		panic(err)
	}
	log.Printf("Created cart item event: %s, item %v", streamName, item)
}
