package biz

import (
	"sync"

	"github.com/google/uuid"
)

type EventType string

const (
	EventTypeDashboardUpdate EventType = "dashboard.update"
)

type Event struct {
	Type   EventType
	UserID uuid.UUID
}

// EventBus is a simple in-memory event bus for broadcasting updates.
type EventBus struct {
	mu          sync.RWMutex
	subscribers map[uuid.UUID][]chan *Event
}

func NewEventBus() *EventBus {
	return &EventBus{
		subscribers: make(map[uuid.UUID][]chan *Event),
	}
}

// Subscribe returns a channel that receives events for a specific user.
func (b *EventBus) Subscribe(userID uuid.UUID) chan *Event {
	b.mu.Lock()
	defer b.mu.Unlock()

	ch := make(chan *Event, 10)
	b.subscribers[userID] = append(b.subscribers[userID], ch)
	return ch
}

// Unsubscribe removes a channel from the subscribers list.
func (b *EventBus) Unsubscribe(userID uuid.UUID, ch chan *Event) {
	b.mu.Lock()
	defer b.mu.Unlock()

	subs := b.subscribers[userID]
	for i, sub := range subs {
		if sub == ch {
			b.subscribers[userID] = append(subs[:i], subs[i+1:]...)
			close(ch)
			break
		}
	}
	if len(b.subscribers[userID]) == 0 {
		delete(b.subscribers, userID)
	}
}

// Publish broadcasts an event to all subscribers for a specific user.
func (b *EventBus) Publish(event *Event) {
	b.mu.RLock()
	defer b.mu.RUnlock()

	if subs, ok := b.subscribers[event.UserID]; ok {
		for _, ch := range subs {
			select {
			case ch <- event:
			default:
				// If channel is full, skip to avoid blocking the publisher
			}
		}
	}
}

// Global dashboard event trigger helper
func (b *EventBus) TriggerDashboardUpdate(userID uuid.UUID) {
	b.Publish(&Event{
		Type:   EventTypeDashboardUpdate,
		UserID: userID,
	})
}
