package main

import (
	"context"
	"errors"
	"fmt"
	"io"

	"github.com/kurrent-io/KurrentDB-Client-Go/kurrentdb"
)

func Read() {
	stream, err := db.ReadStream(context.Background(), "orders", kurrentdb.ReadStreamOptions{}, 10)

	if err != nil {
		panic(err)
	}

	defer stream.Close()

	for {
		event, err := stream.Recv()

		if errors.Is(err, io.EOF) {
			break
		}

		if err != nil {
			panic(err)
		}

		// Process the order event
		fmt.Printf("Order event: %s - %s\n", event.Event.EventType, string(event.Event.Data))
	}
}
