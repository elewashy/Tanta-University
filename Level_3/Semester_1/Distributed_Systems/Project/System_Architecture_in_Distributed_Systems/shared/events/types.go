package events

import "time"

// EventType represents the type of event
type EventType string

const (
	OrderCreatedEvent      EventType = "order.created"
	PaymentProcessedEvent  EventType = "payment.processed"
)

// Event is the base event structure
type Event struct {
	ID        string                 `json:"id"`
	Type      EventType              `json:"type"`
	Timestamp time.Time              `json:"timestamp"`
	Data      map[string]interface{} `json:"data"`
}

// OrderCreatedData represents the data for an order created event
type OrderCreatedData struct {
	OrderID    string   `json:"order_id"`
	CustomerID string   `json:"customer_id"`
	Items      []string `json:"items"`
	Total      float64  `json:"total"`
}

// PaymentProcessedData represents the data for a payment processed event
type PaymentProcessedData struct {
	OrderID       string  `json:"order_id"`
	PaymentID     string  `json:"payment_id"`
	Amount        float64 `json:"amount"`
	Status        string  `json:"status"`
}
