package application

import "context"

type StreamExpectation int

const (
	ExpectAny StreamExpectation = iota
	ExpectNew
	ExpectExisting
)

type EventRecord struct {
	StreamID    string
	EventNumber uint64
	EventType   string
	Data        []byte
}

type EventStore interface {
	Append(ctx context.Context, streamID string, expectation StreamExpectation, events ...EventRecord) error
	ReadStream(ctx context.Context, streamID string, count uint64) ([]EventRecord, error)
	Subscribe(ctx context.Context, streamID string, fromEnd bool, handler func(EventRecord) error) error
}
