#!/bin/bash

# Script to test the distributed system
# Usage: ./test-system.sh

echo "Testing Distributed System..."
echo ""

# Check if services are running
echo "1. Checking service health..."
curl -s http://localhost:8080/health && echo " ✓ Order Service is healthy"
curl -s http://localhost:8081/health && echo " ✓ Payment Service is healthy"
curl -s http://localhost:8082/health && echo " ✓ Notification Service is healthy"
echo ""

# Create an order
echo "2. Creating a new order..."
ORDER_RESPONSE=$(curl -s -X POST http://localhost:8080/orders \
  -H "Content-Type: application/json" \
  -d '{
    "customer_id": "cust-123",
    "items": ["laptop", "mouse", "keyboard"],
    "total": 1299.99
  }')

echo "Order Response:"
echo $ORDER_RESPONSE | jq .
ORDER_ID=$(echo $ORDER_RESPONSE | jq -r '.id')
echo ""

# Wait for async processing
echo "3. Waiting for event processing (5 seconds)..."
sleep 5
echo ""

# Retrieve the order
echo "4. Retrieving order $ORDER_ID..."
curl -s http://localhost:8080/orders/$ORDER_ID | jq .
echo ""

# Create another order
echo "5. Creating another order..."
curl -s -X POST http://localhost:8080/orders \
  -H "Content-Type: application/json" \
  -d '{
    "customer_id": "cust-456",
    "items": ["phone"],
    "total": 899.99
  }' | jq .
echo ""

# Wait for processing
sleep 3

# Get all orders
echo "6. Retrieving all orders..."
curl -s http://localhost:8080/orders | jq .
echo ""

echo "✓ Test completed! Check service logs to see the event flow."
