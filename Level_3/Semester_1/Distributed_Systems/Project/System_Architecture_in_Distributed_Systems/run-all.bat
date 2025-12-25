@echo off
REM Script to run all services on Windows
REM Usage: run-all.bat

echo Starting all services...

REM Create logs directory if it doesn't exist
if not exist logs mkdir logs

REM Start Order Service
echo Starting Order Service on port 8080...
start "Order Service" cmd /k "cd order-service && go run main.go"
timeout /t 2 /nobreak >nul

REM Start Payment Service
echo Starting Payment Service on port 8081...
start "Payment Service" cmd /k "cd payment-service && go run main.go"
timeout /t 2 /nobreak >nul

REM Start Notification Service
echo Starting Notification Service on port 8082...
start "Notification Service" cmd /k "cd notification-service && go run main.go"

echo.
echo All services started in separate windows!
echo.
echo To test the system, run:
echo   curl -X POST http://localhost:8080/orders -H "Content-Type: application/json" -d "{\"customer_id\":\"cust-123\",\"items\":[\"item1\",\"item2\"],\"total\":99.99}"
echo.
echo Close the service windows to stop them.
