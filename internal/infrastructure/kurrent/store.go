package kurrent

import (
	"context"
	"fmt"
	"io"

	"eventsourcepoc/internal/application"

	"github.com/kurrent-io/KurrentDB-Client-Go/kurrentdb"
)

type Store struct {
	client *kurrentdb.Client
}

func NewStore(connectionString string) (*Store, error) {
	cfg, err := kurrentdb.ParseConnectionString(connectionString)
	if err != nil {
		return nil, fmt.Errorf("parse connection string: %w", err)
	}

	client, err := kurrentdb.NewClient(cfg)
	if err != nil {
		return nil, fmt.Errorf("create KurrentDB client: %w", err)
	}

	return &Store{client: client}, nil
}

func (s *Store) Close() error {
	return s.client.Close()
}

func (s *Store) Append(ctx context.Context, streamID string, expectation application.StreamExpectation, events ...application.EventRecord) error {
	kurrentEvents := make([]kurrentdb.EventData, 0, len(events))
	for _, event := range events {
		kurrentEvents = append(kurrentEvents, kurrentdb.EventData{
			ContentType: kurrentdb.ContentTypeJson,
			EventType:   event.EventType,
			Data:        event.Data,
		})
	}

	_, err := s.client.AppendToStream(ctx, streamID, kurrentdb.AppendToStreamOptions{
		StreamState: s.toExpectedRevision(expectation),
	}, kurrentEvents...)
	if err != nil {
		return fmt.Errorf("append to stream %s: %w", streamID, err)
	}

	return nil
}

func (s *Store) ReadStream(ctx context.Context, streamID string, count uint64) ([]application.EventRecord, error) {
	stream, err := s.client.ReadStream(ctx, streamID, kurrentdb.ReadStreamOptions{From: kurrentdb.Start{}}, count)
	if err != nil {
		return nil, fmt.Errorf("read stream %s: %w", streamID, err)
	}
	defer stream.Close()

	records := make([]application.EventRecord, 0)
	for {
		event, err := stream.Recv()
		if err == io.EOF {
			return records, nil
		}
		if err != nil {
			return nil, fmt.Errorf("receive stream event from %s: %w", streamID, err)
		}
		if event == nil || event.Event == nil {
			continue
		}

		records = append(records, toEventRecord(streamID, *event))
	}
}

func (s *Store) Subscribe(ctx context.Context, streamID string, fromEnd bool, handler func(application.EventRecord) error) error {
	options := kurrentdb.SubscribeToStreamOptions{From: kurrentdb.Start{}}
	if fromEnd {
		options.From = kurrentdb.End{}
	}

	subscription, err := s.client.SubscribeToStream(ctx, streamID, options)
	if err != nil {
		return fmt.Errorf("subscribe to stream %s: %w", streamID, err)
	}
	defer subscription.Close()

	for {
		message := subscription.Recv()
		if message == nil {
			return nil
		}
		if message.SubscriptionDropped != nil {
			if ctx.Err() != nil {
				return ctx.Err()
			}
			return fmt.Errorf("subscription dropped for stream %s: %v", streamID, message.SubscriptionDropped)
		}
		if message.EventAppeared == nil {
			continue
		}

		record := toEventRecord(streamID, *message.EventAppeared)
		if err := handler(record); err != nil {
			return err
		}
	}
}

func (s *Store) toExpectedRevision(expectation application.StreamExpectation) kurrentdb.StreamState {
	switch expectation {
	case application.ExpectNew:
		return kurrentdb.NoStream{}
	case application.ExpectExisting:
		return kurrentdb.StreamExists{}
	default:
		return kurrentdb.Any{}
	}
}

func toEventRecord(streamID string, event kurrentdb.ResolvedEvent) application.EventRecord {
	original := event.OriginalEvent()
	return application.EventRecord{
		StreamID:    streamID,
		EventNumber: original.EventNumber,
		EventType:   original.EventType,
		Data:        original.Data,
	}
}
