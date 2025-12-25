@echo off
REM Script to test the distributed system on Windows
REM Usage: test-system.bat

echo Testing Distributed System...
echo.

echo 1. Checking service health...
curl -s http://localhost:8080/health
echo  - Order Service
curl -s http://localhost:8081/health
echo  - Payment Service
curl -s http://localhost:8082/health
echo  - Notification Service
echo.

echo 2. Creating a new order...
curl -X POST http://localhost:8080/orders -H "Content-Type: application/json" -d "{\"customer_id\":\"cust-123\",\"items\":[\"laptop\",\"mouse\",\"keyboard\"],\"total\":1299.99}"
echo.
echo.

echo 3. Waiting for event processing (5 seconds)...
timeout /t 5 /nobreak >nul
echo.

echo 4. Creating another order...
curl -X POST http://localhost:8080/orders -H "Content-Type: application/json" -d "{\"customer_id\":\"cust-456\",\"items\":[\"phone\"],\"total\":899.99}"
echo.
echo.

echo 5. Retrieving all orders...
curl http://localhost:8080/orders
echo.
echo.

echo Test completed! Check service windows to see the event flow.
pause
