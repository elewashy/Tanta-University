package repository

import (
	"errors"
	"sync"
)

// Order represents an order entity
type Order struct {
	ID         string   `json:"id"`
	CustomerID string   `json:"customer_id"`
	Items      []string `json:"items"`
	Total      float64  `json:"total"`
	Status     string   `json:"status"`
}

// OrderRepository handles data persistence
// Demonstrates: Repository pattern (Data Access Layer)
type OrderRepository struct {
	orders map[string]*Order
	mu     sync.RWMutex // Concurrent safety for distributed access
}

// NewOrderRepository creates a new repository instance
func NewOrderRepository() *OrderRepository {
	return &OrderRepository{
		orders: make(map[string]*Order),
	}
}

// Save stores an order
// Demonstrates: Thread-safe write operation
func (r *OrderRepository) Save(order *Order) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	
	r.orders[order.ID] = order
	return nil
}

// FindByID retrieves an order by ID
// Demonstrates: Thread-safe read operation
func (r *OrderRepository) FindByID(id string) (*Order, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	
	order, exists := r.orders[id]
	if !exists {
		return nil, errors.New("order not found")
	}
	
	return order, nil
}

// FindAll retrieves all orders
func (r *OrderRepository) FindAll() []*Order {
	r.mu.RLock()
	defer r.mu.RUnlock()
	
	orders := make([]*Order, 0, len(r.orders))
	for _, order := range r.orders {
		orders = append(orders, order)
	}
	
	return orders
}
