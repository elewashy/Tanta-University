package api

import (
	"encoding/json"
	"log"
	"net/http"

	"order-service/service"

	"github.com/gorilla/mux"
)

// OrderHandler handles HTTP requests
// Demonstrates: API/Presentation layer
type OrderHandler struct {
	service *service.OrderService
}

// NewOrderHandler creates a new handler instance
func NewOrderHandler(service *service.OrderService) *OrderHandler {
	return &OrderHandler{
		service: service,
	}
}

// CreateOrder handles POST /orders
// Demonstrates: Client-Server pattern (HTTP endpoint)
func (h *OrderHandler) CreateOrder(w http.ResponseWriter, r *http.Request) {
	var req service.CreateOrderRequest
	
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	
	// Validate request
	if req.CustomerID == "" || len(req.Items) == 0 || req.Total <= 0 {
		http.Error(w, "Invalid order data", http.StatusBadRequest)
		return
	}
	
	order, err := h.service.CreateOrder(req)
	if err != nil {
		log.Printf("Error creating order: %v", err)
		http.Error(w, "Failed to create order", http.StatusInternalServerError)
		return
	}
	
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(order)
	
	log.Printf("Order created: %s", order.ID)
}

// GetOrder handles GET /orders/{id}
func (h *OrderHandler) GetOrder(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	
	order, err := h.service.GetOrder(id)
	if err != nil {
		http.Error(w, "Order not found", http.StatusNotFound)
		return
	}
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(order)
}

// GetAllOrders handles GET /orders
func (h *OrderHandler) GetAllOrders(w http.ResponseWriter, r *http.Request) {
	orders := h.service.GetAllOrders()
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(orders)
}
