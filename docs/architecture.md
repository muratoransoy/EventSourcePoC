# Architecture Note

This POC uses event sourcing because the domain is naturally expressed as a sequence of facts:

- a cart was created
- an item was added
- an item was removed
- an order was created from the cart

That makes a stream-based model a better teaching tool than a CRUD-first table design.

## Why it fits this POC

- The shopping cart history is meaningful on its own.
- Replaying events is simple enough to understand in one file.
- Subscriptions are easy to demonstrate with real domain events.

## Trade-offs

- Reads are more expensive unless projections or snapshots are introduced.
- Event contracts need more care than table schemas.
- Aggregate boundaries and stream naming become first-class design choices.

## Scope decisions

This repository intentionally does not add:

- an HTTP API
- authentication
- distributed consumers
- persistent projections
- snapshotting

Those are useful next steps, but they would hide the core mechanics this POC is meant to teach.
