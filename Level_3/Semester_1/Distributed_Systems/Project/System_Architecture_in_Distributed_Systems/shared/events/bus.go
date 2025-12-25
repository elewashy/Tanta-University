package events

import (
	"log"
	"sync"
)

// Handler is a function that processes events
type Handler func(event Event)

// EventBus is an in-memory event bus for demonstration
// In production, use Kafka, RabbitMQ, or NATS
type EventBus struct {
	subscribers map[EventType][]Handler
	mu          sync.RWMutex
}

var (
	globalBus *EventBus
	once      sync.Once
)

// GetEventBus returns the singleton event bus instance
func GetEventBus() *EventBus {
	once.Do(func() {
		globalBus = &EventBus{
			subscribers: make(map[EventType][]Handler),
		}
	})
	return globalBus
}

// Subscribe registers a handler for a specific event type
// Demonstrates: Loose coupling - subscribers don't know about publishers
func (eb *EventBus) Subscribe(eventType EventType, handler Handler) {
	eb.mu.Lock()
	defer eb.mu.Unlock()
	
	eb.subscribers[eventType] = append(eb.subscribers[eventType], handler)
	log.Printf("[EventBus] Subscribed to event type: %s", eventType)
}

// Publish sends an event to all registered handlers
// Demonstrates: Asynchronous communication - fire and forget
func (eb *EventBus) Publish(event Event) {
	eb.mu.RLock()
	handlers := eb.subscribers[event.Type]
	eb.mu.RUnlock()
	
	log.Printf("[EventBus] Publishing event: %s (ID: %s)", event.Type, event.ID)
	
	// Process handlers asynchronously to avoid blocking
	// Demonstrates: Non-blocking event processing
	for _, handler := range handlers {
		go func(h Handler) {
			defer func() {
				if r := recover(); r != nil {
					log.Printf("[EventBus] Handler panic recovered: %v", r)
				}
			}()
			h(event)
		}(handler)
	}
}
