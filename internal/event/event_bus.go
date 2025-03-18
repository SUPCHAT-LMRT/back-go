package event

import "sync"

type EventType string

type Event interface {
	Type() EventType
}

type Listener func(event Event)

type EventBus struct {
	listeners map[EventType][]Listener
	mu        sync.RWMutex
}

func NewEventBus() *EventBus {
	return &EventBus{
		listeners: make(map[EventType][]Listener),
	}
}

func (eb *EventBus) Subscribe(eventType EventType, listener Listener) {
	eb.mu.Lock()
	defer eb.mu.Unlock()
	eb.listeners[eventType] = append(eb.listeners[eventType], listener)
}

func (eb *EventBus) Publish(event Event) {
	eb.mu.RLock()
	defer eb.mu.RUnlock()
	if listeners, ok := eb.listeners[event.Type()]; ok {
		for _, listener := range listeners {
			go listener(event) // Run listeners asynchronously
		}
	}
}
