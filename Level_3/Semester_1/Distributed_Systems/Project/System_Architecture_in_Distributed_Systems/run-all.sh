#!/bin/bash

# Script to run all services locally
# Usage: ./run-all.sh

echo "Starting all services..."

# Kill any existing processes on these ports
lsof -ti:8080 | xargs kill -9 2>/dev/null
lsof -ti:8081 | xargs kill -9 2>/dev/null
lsof -ti:8082 | xargs kill -9 2>/dev/null

# Start services in background
echo "Starting Order Service on port 8080..."
cd order-service && go run main.go > ../logs/order-service.log 2>&1 &
ORDER_PID=$!

echo "Starting Payment Service on port 8081..."
cd ../payment-service && go run main.go > ../logs/payment-service.log 2>&1 &
PAYMENT_PID=$!

echo "Starting Notification Service on port 8082..."
cd ../notification-service && go run main.go > ../logs/notification-service.log 2>&1 &
NOTIFICATION_PID=$!

cd ..

echo ""
echo "All services started!"
echo "Order Service PID: $ORDER_PID"
echo "Payment Service PID: $PAYMENT_PID"
echo "Notification Service PID: $NOTIFICATION_PID"
echo ""
echo "Logs are in ./logs/ directory"
echo ""
echo "To test the system, run:"
echo "  curl -X POST http://localhost:8080/orders -H 'Content-Type: application/json' -d '{\"customer_id\":\"cust-123\",\"items\":[\"item1\",\"item2\"],\"total\":99.99}'"
echo ""
echo "To stop all services, run: ./stop-all.sh"
