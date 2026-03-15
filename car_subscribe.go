package main

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/google/uuid"
	"github.com/kurrent-io/KurrentDB-Client-Go/kurrentdb"
)

func SubscribeToCartEvents(guid string) {
	EnsureStreamExist(streamID)

	log.Printf("Subscribing to cart events for cart %s", guid)
	subscription, err := db.SubscribeToStream(
		context.Background(),
		streamID,
		kurrentdb.SubscribeToStreamOptions{},
	)

	if subscription == nil {
		log.Fatalf("Subscription is nil")
		return
	}

	if err != nil {
		log.Fatalf("Failed to subscribe to stream: %v", err)
	}

	for {
		event := subscription.Recv()

		if event == nil {
			log.Println("Subscription ended")
			break
		}
		if event.EventAppeared != nil {
			log.Printf("Received event: %s", event.EventAppeared.Event.EventType)
			log.Printf("Event data: %s", string(event.EventAppeared.Event.Data))
		}
		if event.SubscriptionDropped != nil {
			log.Printf("Subscription dropped: %v", event.SubscriptionDropped)
			break
		}
	}
}

func EnsureStreamExist(streamID string) {
	// Ensure the stream exists first
	data, err := db.ReadStream(context.Background(), streamID, kurrentdb.ReadStreamOptions{}, 1)
	if err == nil {
		evt, err := data.Recv()
		if err == nil {
			log.Printf("evt is exist %s", evt.Event.StreamID)
			return
		}
	}

	log.Printf("Stream %s does not exist, creating it...", streamID)
	// Append a dummy event just to create the stream
	_, appendErr := db.AppendToStream(
		context.Background(),
		streamID,
		kurrentdb.AppendToStreamOptions{StreamState: kurrentdb.Any{}},
		kurrentdb.EventData{
			EventType:   "stream-created",
			ContentType: kurrentdb.ContentTypeJson,
			Data:        []byte(`{"created": true}`),
		},
	)
	if appendErr != nil {
		log.Fatalf("Failed to create stream: %v", appendErr)
	}
}

func CreateShoppingCartProjection() {
	script := `
fromAll()
.when({
	$init:function(){
		return {
			count: 0
		}
	},
	myEventUpdatedType: function(state, event){
		state.count += 1;
	}
})
.transformBy(function(state){
	state.count = 10;
})
.outputState()
`

	name := fmt.Sprintf("shoppingCartItemsEvent_Create_%s", uuid.New())
	client, err := kurrentdb.NewProjectionClient(settings)
	if err != nil {
		panic(err)
	}
	err = client.Create(context.Background(), name, script, kurrentdb.CreateProjectionOptions{})

	if esdbErr, ok := kurrentdb.FromError(err); !ok {
		if esdbErr.IsErrorCode(kurrentdb.ErrorCodeUnknown) && strings.Contains(esdbErr.Err().Error(), "Conflict") {
			log.Printf("projection %s already exists", name)
			return
		}
	}
}
