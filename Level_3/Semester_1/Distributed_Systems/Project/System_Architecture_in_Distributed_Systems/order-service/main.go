package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"order-service/api"
	"order-service/repository"
	"order-service/service"
	"shared/events"

	"github.com/gorilla/mux"
)

func main() {
	log.Println("Starting Order Service on port 8080...")
	
	// Initialize layers (Layered Architecture)
	eventBus := events.GetEventBus()
	repo := repository.NewOrderRepository()
	svc := service.NewOrderService(repo, eventBus)
	handler := api.NewOrderHandler(svc)
	
	// Setup HTTP router
	router := mux.NewRouter()
	router.HandleFunc("/orders", handler.CreateOrder).Methods("POST")
	router.HandleFunc("/orders/{id}", handler.GetOrder).Methods("GET")
	router.HandleFunc("/orders", handler.GetAllOrders).Methods("GET")
	
	// Health check endpoint
	router.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	}).Methods("GET")
	
	// Create HTTP server
	srv := &http.Server{
		Addr:         ":8080",
		Handler:      router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}
	
	// Start server in goroutine
	go func() {
		log.Println("Order Service listening on :8080")
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server error: %v", err)
		}
	}()
	
	// Graceful shutdown
	// Demonstrates: Production-ready service lifecycle management
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	
	log.Println("Shutting down Order Service...")
	
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}
	
	log.Println("Order Service stopped")
}
