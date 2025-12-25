package service

import (
	"time"

	"order-service/repository"
	"shared/events"

	"github.com/google/uuid"
)

// CreateOrderRequest represents the request to create an order
type CreateOrderRequest struct {
	CustomerID string   `json:"customer_id"`
	Items      []string `json:"items"`
	Total      float64  `json:"total"`
}

// OrderService handles business logic
// Demonstrates: Service layer (Business Logic Layer)
type OrderService struct {
	repo     *repository.OrderRepository
	eventBus *events.EventBus
}

// NewOrderService creates a new service instance
func NewOrderService(repo *repository.OrderRepository, eventBus *events.EventBus) *OrderService {
	return &OrderService{
		repo:     repo,
		eventBus: eventBus,
	}
}

// CreateOrder creates a new order and publishes an event
// Demonstrates: Business logic + Event-driven architecture
func (s *OrderService) CreateOrder(req CreateOrderRequest) (*repository.Order, error) {
	// Create order entity
	order := &repository.Order{
		ID:         uuid.New().String(),
		CustomerID: req.CustomerID,
		Items:      req.Items,
		Total:      req.Total,
		Status:     "pending",
	}
	
	// Persist to repository
	if err := s.repo.Save(order); err != nil {
		return nil, err
	}
	
	// Publish event for other services
	// Demonstrates: Loose coupling via events
	event := events.Event{
		ID:        uuid.New().String(),
		Type:      events.OrderCreatedEvent,
		Timestamp: time.Now(),
		Data: map[string]interface{}{
			"order_id":    order.ID,
			"customer_id": order.CustomerID,
			"items":       order.Items,
			"total":       order.Total,
		},
	}
	
	s.eventBus.Publish(event)
	
	return order, nil
}

// GetOrder retrieves an order by ID
func (s *OrderService) GetOrder(id string) (*repository.Order, error) {
	return s.repo.FindByID(id)
}

// GetAllOrders retrieves all orders
func (s *OrderService) GetAllOrders() []*repository.Order {
	return s.repo.FindAll()
}
