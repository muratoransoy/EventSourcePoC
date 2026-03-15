# Go Event Sourcing POC with KurrentDB

This repository is a small, runnable Go proof of concept that demonstrates event sourcing and stream-based thinking with [KurrentDB](https://www.kurrent.io/) using a shopping cart and order flow.

It focuses on a few core ideas:

- append-only events instead of mutable rows
- one stream per aggregate
- rebuilding state by replaying events
- subscribing to a stream to observe changes as they happen

The project is intentionally small. It is meant to be educational, practical, and easy to publish as a backend architecture POC.

## What This Project Demonstrates

- Modeling a shopping cart as an event stream
- Appending events such as `CartCreated`, `ItemAddedToCart`, `ItemRemovedFromCart`, and `OrderCreated`
- Reading streams back from KurrentDB
- Rebuilding current cart state from historical events
- Creating an order by projecting cart events into a current snapshot
- Subscribing to a stream to observe event flow

## Why Event Sourcing

Event sourcing stores facts that happened to a domain model instead of only persisting the latest state.

In this POC that means:

- a cart is represented by the events written to `cart-<cart-id>`
- an order is represented by the events written to `order-<order-id>`
- the current cart is derived by replaying the cart stream

Benefits:

- full audit trail
- replayability
- easier temporal debugging
- a natural fit for event-driven integration

Trade-offs:

- reads require projection or replay
- schema evolution must be handled intentionally
- event naming and stream boundaries matter more than in CRUD-first systems

More detail is in [docs/architecture.md](/C:/Users/murat/Desktop/EventSourcePoC/docs/architecture.md).

## Architecture Overview

The code uses a light layered structure:

- `internal/domain/`: event types and projection logic
- `internal/application/`: use cases and stream naming
- `internal/infrastructure/`: KurrentDB adapter and config loading
- `internal/interfaces/cli/`: CLI command handling
- `cmd/eventsource-poc/`: entrypoint

There is no heavy DDD or CQRS framework here. The goal is to keep the event sourcing mechanics visible.

### Event sourcing in this example

1. `create-cart` appends `CartCreated` to a cart stream.
2. `add-item` appends `ItemAddedToCart`.
3. `remove-item` appends `ItemRemovedFromCart`.
4. `create-order` reads the cart stream, rebuilds current state, and appends `OrderCreated` to a new order stream.
5. `subscribe` listens to a stream and prints decoded events as they appear.

### Projection logic

Projection logic is deliberately simple:

- cart state is rebuilt by replaying the cart stream in [internal/domain/cart/state.go](/C:/Users/murat/Desktop/EventSourcePoC/internal/domain/cart/state.go)
- `OrderCreated` is generated from that projected state

That keeps projection visible without adding a separate read-model service.

## Project Structure

```text
cmd/
  eventsource-poc/
    main.go
internal/
  application/
    cart_service.go
    codec.go
    event_store.go
    order_service.go
    streams.go
  domain/
    cart/
      events.go
      state.go
    order/
      events.go
  infrastructure/
    config/
      config.go
    kurrent/
      store.go
  interfaces/
    cli/
      app.go
docs/
  architecture.md
```

Responsibilities:

- `cmd/eventsource-poc`: wires configuration, infrastructure, and application services
- `internal/domain/cart`: defines cart events and how a cart is projected from a stream
- `internal/domain/order`: defines the `OrderCreated` event
- `internal/application`: contains the use cases and stream naming rules
- `internal/infrastructure/kurrent`: wraps the KurrentDB Go client behind a small event store interface
- `internal/interfaces/cli`: exposes the POC as explicit commands

## Event Flow

The shopping cart flow is modeled with append-only events:

```mermaid
flowchart LR
    A["create-cart"] --> B["CartCreated -> cart-<cart-id>"]
    C["add-item"] --> D["ItemAddedToCart -> cart-<cart-id>"]
    E["remove-item"] --> F["ItemRemovedFromCart -> cart-<cart-id>"]
    B --> G["Replay cart stream"]
    D --> G
    F --> G
    G --> H["create-order"]
    H --> I["OrderCreated -> order-<order-id>"]
    D --> J["subscribe / read-stream"]
    F --> J
    I --> J
```

Event types used in the POC:

- `CartCreated`
- `ItemAddedToCart`
- `ItemRemovedFromCart`
- `OrderCreated`

Stream naming convention:

- carts: `cart-<uuid>`
- orders: `order-<uuid>`

## Requirements

- Go `1.24+`
- A running KurrentDB instance
- Docker only if you want the quickest local KurrentDB setup

## Setup

The CLI reads the connection string from `KURRENTDB_CONNECTION_STRING`.

Default value:

```bash
kurrentdb://admin:changeit@localhost:2113?tls=false
```

### Start KurrentDB with Docker

One simple local option:

```bash
docker run --rm --name kurrentdb ^
  -p 2113:2113 ^
  -e KURRENTDB_INSECURE=true ^
  docker.kurrent.io/kurrent-latest/kurrentdb:latest
```

If you already run EventStoreDB/KurrentDB another way, just point the connection string at it.

## Running the Project

### Run the CLI

```bash
go run ./cmd/eventsource-poc --help
```

### Use an environment variable

```bash
set KURRENTDB_CONNECTION_STRING=kurrentdb://admin:changeit@localhost:2113?tls=false
go run ./cmd/eventsource-poc sample-flow
```

### Or pass the connection string explicitly

```bash
go run ./cmd/eventsource-poc --connection "kurrentdb://admin:changeit@localhost:2113?tls=false" sample-flow
```

## Example Commands

Create a cart:

```bash
go run ./cmd/eventsource-poc create-cart --cart-id 2df23953-f319-4f04-b09c-2c39577987df --user-id 2454b0c6-95d4-4d2e-84d7-4f4b7487168b --currency TRY
```

Add an item:

```bash
go run ./cmd/eventsource-poc add-item --cart-id 2df23953-f319-4f04-b09c-2c39577987df --product-id 763c646a-bb66-422f-b23c-f3c34c62ee86 --name "MacBook Pro" --quantity 1 --price 190000
```

Remove an item:

```bash
go run ./cmd/eventsource-poc remove-item --cart-id 2df23953-f319-4f04-b09c-2c39577987df --product-id 763c646a-bb66-422f-b23c-f3c34c62ee86 --name "MacBook Pro" --quantity 1
```

Create an order from the current cart state:

```bash
go run ./cmd/eventsource-poc create-order --cart-id 2df23953-f319-4f04-b09c-2c39577987df --order-id 6703921b-373c-4e85-beb1-5926d53d1637
```

Read a stream:

```bash
go run ./cmd/eventsource-poc read-stream --stream cart-2df23953-f319-4f04-b09c-2c39577987df
```

Subscribe to a stream:

```bash
go run ./cmd/eventsource-poc subscribe --stream cart-2df23953-f319-4f04-b09c-2c39577987df
```

Run the full walkthrough:

```bash
go run ./cmd/eventsource-poc sample-flow
```

## Example Workflow

1. Create a cart.
2. Add one or more items.
3. Optionally remove an item.
4. Read the cart stream to see the append-only history.
5. Create an order from the cart.
6. Read the order stream or subscribe to either stream.

Internally:

- commands append JSON events to KurrentDB
- cart state is rebuilt by replaying the cart stream
- order creation uses that projected state instead of a mutable cart table

## Event Model Example

`CartCreated`

```json
{
  "cartId": "2df23953-f319-4f04-b09c-2c39577987df",
  "userId": "2454b0c6-95d4-4d2e-84d7-4f4b7487168b",
  "currency": "TRY",
  "occurredAt": "2026-03-15T12:00:00Z"
}
```

`ItemAddedToCart`

```json
{
  "cartId": "2df23953-f319-4f04-b09c-2c39577987df",
  "productId": "763c646a-bb66-422f-b23c-f3c34c62ee86",
  "productName": "MacBook Pro",
  "quantity": 1,
  "unitPrice": 190000,
  "occurredAt": "2026-03-15T12:01:00Z"
}
```

`OrderCreated`

```json
{
  "orderId": "6703921b-373c-4e85-beb1-5926d53d1637",
  "cartId": "2df23953-f319-4f04-b09c-2c39577987df",
  "userId": "2454b0c6-95d4-4d2e-84d7-4f4b7487168b",
  "currency": "TRY",
  "items": [
    {
      "productId": "763c646a-bb66-422f-b23c-f3c34c62ee86",
      "productName": "MacBook Pro",
      "quantity": 1,
      "unitPrice": 190000
    }
  ],
  "totalAmount": 190000,
  "occurredAt": "2026-03-15T12:05:00Z"
}
```

## Testing

The repo includes a few focused tests:

- event payload JSON serialization
- stream naming rules
- order creation from replayed cart state

Run them with:

```bash
go test ./...
```

## Future Improvements

- add versioned event upcasters
- add snapshots for long-running streams
- introduce persistent subscriptions
- add a small read model or projection worker
- expose the same use cases through HTTP or gRPC
- improve integration tests against a real KurrentDB instance

## Learning Goals

This repository is aimed at backend engineers who want to see:

- how event sourcing looks in idiomatic Go
- how stream naming and event naming shape the model
- how to reconstruct state from append-only history
- how KurrentDB fits into a lightweight backend architecture

## License

This repository includes an MIT license in [LICENSE](/C:/Users/murat/Desktop/EventSourcePoC/LICENSE).
