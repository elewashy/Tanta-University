package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"shared/events"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

// PaymentService handles payment processing
// Demonstrates: Microservice with event-driven architecture
type PaymentService struct {
	eventBus *events.EventBus
}

func NewPaymentService(eventBus *events.EventBus) *PaymentService {
	return &PaymentService{
		eventBus: eventBus,
	}
}

// ProcessOrderCreated handles order created events
// Demonstrates: Event-driven processing, Asynchronous communication
func (ps *PaymentService) ProcessOrderCreated(event events.Event) {
	log.Printf("[Payment] Received OrderCreated event: %s", event.ID)
	
	orderID := event.Data["order_id"].(string)
	total := event.Data["total"].(float64)
	
	// Simulate payment processing
	log.Printf("[Payment] Processing payment for order %s, amount: $%.2f", orderID, total)
	time.Sleep(2 * time.Second) // Simulate processing time
	
	// Simulate payment success (90% success rate)
	status := "success"
	if time.Now().Unix()%10 == 0 {
		status = "failed"
	}
	
	log.Printf("[Payment] Payment %s for order %s", status, orderID)
	
	// Publish payment processed event
	// Demonstrates: Event chain - one event triggers another
	paymentEvent := events.Event{
		ID:        uuid.New().String(),
		Type:      events.PaymentProcessedEvent,
		Timestamp: time.Now(),
		Data: map[string]interface{}{
			"order_id":   orderID,
			"payment_id": uuid.New().String(),
			"amount":     total,
			"status":     status,
		},
	}
	
	ps.eventBus.Publish(paymentEvent)
}

func main() {
	log.Println("Starting Payment Service on port 8081...")
	
	// Initialize event bus
	eventBus := events.GetEventBus()
	paymentService := NewPaymentService(eventBus)
	
	// Subscribe to OrderCreated events
	// Demonstrates: Loose coupling - Payment service doesn't know about Order service
	eventBus.Subscribe(events.OrderCreatedEvent, paymentService.ProcessOrderCreated)
	
	// Setup HTTP server for health checks
	router := mux.NewRouter()
	router.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	}).Methods("GET")
	
	srv := &http.Server{
		Addr:         ":8081",
		Handler:      router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}
	
	// Start server
	go func() {
		log.Println("Payment Service listening on :8081")
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server error: %v", err)
		}
	}()
	
	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	
	log.Println("Shutting down Payment Service...")
	
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}
	
	log.Println("Payment Service stopped")
}
