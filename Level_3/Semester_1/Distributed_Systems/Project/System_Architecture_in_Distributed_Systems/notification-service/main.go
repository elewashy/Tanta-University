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

	"github.com/gorilla/mux"
)

// NotificationService handles sending notifications
// Demonstrates: Microservice that consumes multiple event types
type NotificationService struct {
	eventBus *events.EventBus
}

func NewNotificationService(eventBus *events.EventBus) *NotificationService {
	return &NotificationService{
		eventBus: eventBus,
	}
}

// NotifyOrderCreated sends notification when order is created
// Demonstrates: Multiple services can react to the same event
func (ns *NotificationService) NotifyOrderCreated(event events.Event) {
	log.Printf("[Notification] Received OrderCreated event: %s", event.ID)
	
	orderID := event.Data["order_id"].(string)
	customerID := event.Data["customer_id"].(string)
	
	// Simulate sending notification
	log.Printf("[Notification] Sending order confirmation to customer %s for order %s", customerID, orderID)
	time.Sleep(500 * time.Millisecond) // Simulate email/SMS sending
	
	log.Printf("[Notification] ✓ Order confirmation sent to customer %s", customerID)
}

// NotifyPaymentProcessed sends notification when payment is processed
// Demonstrates: Event-driven workflow across multiple services
func (ns *NotificationService) NotifyPaymentProcessed(event events.Event) {
	log.Printf("[Notification] Received PaymentProcessed event: %s", event.ID)
	
	orderID := event.Data["order_id"].(string)
	status := event.Data["status"].(string)
	amount := event.Data["amount"].(float64)
	
	// Simulate sending notification
	if status == "success" {
		log.Printf("[Notification] Sending payment success notification for order %s (amount: $%.2f)", orderID, amount)
	} else {
		log.Printf("[Notification] Sending payment failure notification for order %s", orderID)
	}
	
	time.Sleep(500 * time.Millisecond)
	log.Printf("[Notification] ✓ Payment notification sent for order %s", orderID)
}

func main() {
	log.Println("Starting Notification Service on port 8082...")
	
	// Initialize event bus
	eventBus := events.GetEventBus()
	notificationService := NewNotificationService(eventBus)
	
	// Subscribe to multiple event types
	// Demonstrates: Service can react to multiple events, Fault isolation
	eventBus.Subscribe(events.OrderCreatedEvent, notificationService.NotifyOrderCreated)
	eventBus.Subscribe(events.PaymentProcessedEvent, notificationService.NotifyPaymentProcessed)
	
	// Setup HTTP server for health checks
	router := mux.NewRouter()
	router.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	}).Methods("GET")
	
	srv := &http.Server{
		Addr:         ":8082",
		Handler:      router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}
	
	// Start server
	go func() {
		log.Println("Notification Service listening on :8082")
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server error: %v", err)
		}
	}()
	
	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	
	log.Println("Shutting down Notification Service...")
	
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}
	
	log.Println("Notification Service stopped")
}
