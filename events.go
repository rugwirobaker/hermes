package helmes

import (
	"context"
	"errors"
	"sync"
)

var _ (Pubsub) = (*pubsub)(nil)

// Status defines delivery status
type Status string

// Possible delivery statuses
const (
	Unknown   Status = "unknown"
	Delivered Status = "delivered"
	Failed    Status = "failed"
	Pending   Status = "pending"
	Submitted Status = "submitted"
	Rejected  Status = "rejected"
)

// St converts int to Status
func St(in int) Status {
	switch in {
	case 1:
		return Delivered
	case 2:
		return Failed
	case 4:
		return Pending
	case 8:
		return Submitted
	case 16:
		return Rejected
	default:
		return Unknown
	}
}

// Event defines delivery event
type Event struct {
	ID        string `json:"id"`
	Status    Status `json:"status,string"`
	Recipient string `json:"recipient"`
}

// Pubsub recieves message delivery events and generates events
type Pubsub interface {
	//Subcribe to receive message delivery event
	Subscribe(context.Context, string) (<-chan Event, error)

	//Publish a new message delivery event
	Publish(ctx context.Context, event Event)

	// Done instructs pubsub to close the event channel when we are done reading
	Done(context.Context, string) error

	//Close the Pubsub
	Close()
}

type pubsub struct {
	mu     sync.RWMutex
	sink   map[string]chan Event
	closed bool
}

// NewPubsub creates a new in-memory pubsub instance
func NewPubsub() Pubsub {
	ps := &pubsub{}
	ps.sink = make(map[string]chan Event)
	ps.closed = false
	return ps
}

// the channel is closed of by the on-off subscriber
func (ps *pubsub) Subscribe(ctx context.Context, topic string) (<-chan Event, error) {
	if ps.exists(topic) {
		return nil, errors.New("already subscribed to this events")
	}
	ps.mu.Lock()
	event := make(chan Event, 1)
	ps.sink[topic] = event
	ps.mu.Unlock()

	go func() {
		select {
		case <-ctx.Done():
			ps.mu.Lock()
			delete(ps.sink, topic)
			ps.mu.Unlock()
		}
	}()
	return event, nil
}

func (ps *pubsub) Publish(ctx context.Context, event Event) {
	ps.mu.RLock()
	defer ps.mu.RUnlock()

	sub := ps.sink[event.ID]
	sub <- event
}

func (ps *pubsub) Done(ctx context.Context, topic string) error {
	ps.mu.RLock()
	defer ps.mu.RUnlock()

	ch, ok := ps.sink[topic]
	if !ok {
		return errors.New("topic doesn't exists")
	}
	close(ch)
	delete(ps.sink, topic)
	return nil
}

func (ps *pubsub) Close() {
	ps.mu.Lock()
	defer ps.mu.Unlock()

	if !ps.closed {
		ps.closed = true
		for _, ch := range ps.sink {
			close(ch)
		}

	}
}

func (ps *pubsub) exists(topic string) bool {
	ps.mu.Lock()
	defer ps.mu.Unlock()
	_, ok := ps.sink[topic]
	return ok
}
