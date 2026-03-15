package main

import (
	orders "EvenSourcePOC/Orders"
	"context"
	"encoding/json"
	"log"

	"github.com/kurrent-io/KurrentDB-Client-Go/kurrentdb"
)

func CreateOrders() {
	for i := 0; i < 5; i++ {
		order := orders.NewOrder()
		data, err := json.Marshal(order)

		if err != nil {
			log.Fatalf("Failed to marshal event data: %v", err)
		}
		eventData := kurrentdb.EventData{
			ContentType: kurrentdb.ContentTypeJson,
			EventType:   "OrderCreated",
			Data:        data,
		}
		if _, err := db.AppendToStream(context.Background(), "orders", kurrentdb.AppendToStreamOptions{}, eventData); err != nil {
			panic(err)
		}
		log.Printf("Appended OrderCreated event for Order ID: %s", order.OrderId)
	}
	Read()
}
