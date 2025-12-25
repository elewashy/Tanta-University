# Architecture Documentation

## System Overview

This distributed system demonstrates four key architectural patterns working together:

1. **Layered Architecture** (Order Service)
2. **Client-Server** (HTTP REST APIs)
3. **Event-Driven Architecture** (Async messaging)
4. **Microservices** (Independent services)

## High-Level Architecture Diagram

```mermaid
graph TB
    Client[Client Application]
    
    subgraph "Order Service - Port 8080"
        API[API Layer<br/>HTTP Handlers]
        Service[Service Layer<br/>Business Logic]
        Repo[Repository Layer<br/>Data Access]
        
        API --> Service
        Service --> Repo
    end
    
    subgraph "Event Bus"
        EventBus[In-Memory Event Bus<br/>Pub/Sub Pattern]
    end
    
    subgraph "Payment Service - Port 8081"
        PaymentHandler[Event Handler<br/>Payment Processing]
    end
    
    subgraph "Notification Service - Port 8082"
        NotifHandler[Event Handler<br/>Notification Sending]
    end
    
    Client -->|HTTP POST/GET| API
    Service -->|Publish Events| EventBus
    EventBus -->|OrderCreated| PaymentHandler
    EventBus -->|OrderCreated| NotifHandler
    PaymentHandler -->|PaymentProcessed| EventBus
    EventBus -->|PaymentProcessed| NotifHandler
    
    style API fill:#e1f5ff
    style Service fill:#e1f5ff
    style Repo fill:#e1f5ff
    style EventBus fill:#fff4e1
    style PaymentHandler fill:#f0e1ff
    style NotifHandler fill:#e1ffe1
```

## Detailed Event Flow

```mermaid
sequenceDiagram
    participant Client
    participant OrderAPI as Order Service<br/>(API Layer)
    participant OrderSvc as Order Service<br/>(Service Layer)
    participant OrderRepo as Order Service<br/>(Repository)
    participant EventBus as Event Bus
    participant Payment as Payment Service
    participant Notification as Notification Service
    
    Client->>OrderAPI: POST /orders
    activate OrderAPI
    OrderAPI->>OrderSvc: CreateOrder(request)
    activate OrderSvc
    OrderSvc->>OrderRepo: Save(order)
    activate OrderRepo
    OrderRepo-->>OrderSvc: order saved
    deactivate OrderRepo
    
    OrderSvc->>EventBus: Publish(OrderCreated)
    activate EventBus
    OrderSvc-->>OrderAPI: order created
    deactivate OrderSvc
    OrderAPI-->>Client: 201 Created
    deactivate OrderAPI
    
    EventBus-->>Payment: OrderCreated event
    EventBus-->>Notification: OrderCreated event
    deactivate EventBus
    
    activate Payment
    Payment->>Payment: Process payment
    Payment->>EventBus: Publish(PaymentProcessed)
    deactivate Payment
    
    activate Notification
    Notification->>Notification: Send order confirmation
    deactivate Notification
    
    activate EventBus
    EventBus-->>Notification: PaymentProcessed event
    deactivate EventBus
    
    activate Notification
    Notification->>Notification: Send payment notification
    deactivate Notification
```

## Architecture Patterns Explained

### 1. Layered Architecture (Order Service)

The Order Service implements a classic 3-tier architecture with clear separation of concerns:

```mermaid
graph LR
    subgraph "Layered Architecture"
        direction TB
        A[API Layer<br/>Presentation] --> B[Service Layer<br/>Business Logic]
        B --> C[Repository Layer<br/>Data Access]
    end
    
    style A fill:#ff9999
    style B fill:#99ccff
    style C fill:#99ff99
```

**API Layer (Presentation)**
- Handles HTTP requests/responses
- Input validation
- JSON serialization
- HTTP status codes

**Service Layer (Business Logic)**
- Order creation logic
- Event publishing
- Business rules enforcement
- Orchestration

**Repository Layer (Data Access)**
- Data persistence
- CRUD operations
- Thread-safe access
- Data retrieval

**Benefits:**
- ✅ Separation of concerns
- ✅ Easy to test each layer independently
- ✅ Clear responsibility boundaries
- ✅ Maintainable and scalable
- ✅ Can swap implementations (e.g., in-memory → PostgreSQL)

**Implementation:**
```
order-service/
├── api/handler.go          # Presentation layer
├── service/order_service.go # Business logic layer
└── repository/order_repository.go # Data access layer
```

### 2. Client-Server Pattern

Traditional request-response communication over HTTP:

```mermaid
graph LR
    Client[Client] -->|HTTP Request| Server[Order Service]
    Server -->|HTTP Response| Client
    
    style Client fill:#e1f5ff
    style Server fill:#ffe1e1
```

**Characteristics:**
- **Synchronous** - Client waits for response
- **Stateless** - Each request is independent
- **RESTful** - Standard HTTP methods (GET, POST, PATCH, DELETE)
- **Clear Contract** - Well-defined API endpoints

**Endpoints:**
- `POST /orders` - Create new order
- `GET /orders/{id}` - Get order by ID
- `GET /orders` - Get all orders
- `GET /health` - Health check

**Example Request/Response:**
```bash
# Request
POST /orders
Content-Type: application/json
{
  "customer_id": "cust-123",
  "items": ["laptop", "mouse"],
  "total": 1299.99
}

# Response
201 Created
{
  "id": "order-abc-123",
  "customer_id": "cust-123",
  "items": ["laptop", "mouse"],
  "total": 1299.99,
  "status": "pending"
}
```

### 3. Event-Driven Architecture

Services communicate asynchronously through events, enabling loose coupling:

```mermaid
graph TB
    subgraph "Event Publishers"
        OrderSvc[Order Service]
        PaymentSvc[Payment Service]
    end
    
    subgraph "Event Bus"
        EventBus[Event Bus<br/>Pub/Sub]
    end
    
    subgraph "Event Consumers"
        PaymentConsumer[Payment Service]
        NotifConsumer[Notification Service]
    end
    
    OrderSvc -->|Publish: OrderCreated| EventBus
    PaymentSvc -->|Publish: PaymentProcessed| EventBus
    
    EventBus -->|Subscribe: OrderCreated| PaymentConsumer
    EventBus -->|Subscribe: OrderCreated| NotifConsumer
    EventBus -->|Subscribe: PaymentProcessed| NotifConsumer
    
    style EventBus fill:#fff4e1
    style OrderSvc fill:#e1f5ff
    style PaymentSvc fill:#f0e1ff
    style PaymentConsumer fill:#f0e1ff
    style NotifConsumer fill:#e1ffe1
```

**Event Types:**

1. **OrderCreated Event**
   - Published by: Order Service
   - Consumed by: Payment Service, Notification Service
   - Payload: order_id, customer_id, items, total

2. **PaymentProcessed Event**
   - Published by: Payment Service
   - Consumed by: Notification Service
   - Payload: order_id, payment_id, amount, status

**Benefits:**
- ✅ **Loose Coupling** - Services don't know about each other
- ✅ **Asynchronous** - Non-blocking operations
- ✅ **Scalability** - Easy to add new consumers
- ✅ **Fault Tolerance** - One service failure doesn't block others
- ✅ **Flexibility** - Easy to add new event types

**Event Flow Example:**
```mermaid
stateDiagram-v2
    [*] --> OrderCreated: Client creates order
    OrderCreated --> PaymentProcessing: Payment Service receives event
    OrderCreated --> NotificationSent: Notification Service receives event
    PaymentProcessing --> PaymentCompleted: Payment processed
    PaymentCompleted --> NotificationSent: Notification Service receives event
    NotificationSent --> [*]
```

### 4. Microservices Architecture

Three independent services with clear boundaries and responsibilities:

```mermaid
graph TB
    subgraph "Microservices Architecture"
        subgraph "Order Service :8080"
            O1[Order Management]
            O2[REST API]
            O3[Event Publishing]
        end
        
        subgraph "Payment Service :8081"
            P1[Payment Processing]
            P2[Event Consumption]
            P3[Event Publishing]
        end
        
        subgraph "Notification Service :8082"
            N1[Email/SMS Sending]
            N2[Event Consumption]
            N3[Multi-channel Notifications]
        end
    end
    
    O3 -.->|Events| P2
    O3 -.->|Events| N2
    P3 -.->|Events| N2
    
    style O1 fill:#e1f5ff
    style P1 fill:#f0e1ff
    style N1 fill:#e1ffe1
```

**Service Characteristics:**

| Service | Port | Responsibility | Communication | Database |
|---------|------|----------------|---------------|----------|
| **Order Service** | 8080 | Order management, CRUD operations | HTTP + Events | In-memory (orders) |
| **Payment Service** | 8081 | Payment processing, validation | Events only | None |
| **Notification Service** | 8082 | Customer notifications | Events only | None |

**Microservice Principles Demonstrated:**

1. **Single Responsibility**
   - Each service has one clear purpose
   - Order Service: Manages orders
   - Payment Service: Processes payments
   - Notification Service: Sends notifications

2. **Independent Deployment**
   - Each service can be deployed separately
   - No shared codebase (except shared events)
   - Own Docker container

3. **Technology Independence**
   - Could use different languages/frameworks
   - Could use different databases
   - Could use different deployment strategies

4. **Fault Isolation**
   - If Payment Service crashes, Order Service continues
   - If Notification Service is down, payments still process
   - Each service has its own lifecycle

5. **Independent Scaling**
   - Scale Order Service for high order volume
   - Scale Payment Service for payment processing
   - Scale Notification Service for notification load

## Distributed System Concepts

### 1. Loose Coupling

Services are independent and communicate through well-defined interfaces:

```mermaid
graph LR
    A[Order Service] -.->|Events| B[Event Bus]
    B -.->|Events| C[Payment Service]
    B -.->|Events| D[Notification Service]
    
    style B fill:#fff4e1
```

**How it's achieved:**
- Services don't call each other directly
- Communication through events
- No shared database
- Services can be added/removed without affecting others

**Benefits:**
- Easy to modify one service without affecting others
- Easy to add new services
- Easy to replace implementations

### 2. Asynchronous Communication

Events are processed asynchronously without blocking:

```mermaid
graph TB
    A[Order Service] -->|1. Create Order| B[Save to DB]
    B -->|2. Publish Event| C[Event Bus]
    A -->|3. Return Response| D[Client]
    C -->|4. Async Processing| E[Payment Service]
    C -->|4. Async Processing| F[Notification Service]
    
    style C fill:#fff4e1
    style D fill:#e1f5ff
```

**Key Points:**
- Order creation returns immediately
- Payment processing happens in background
- Client doesn't wait for payment
- Non-blocking operations

### 3. Fault Isolation

Services fail independently without cascading failures:

```mermaid
graph TB
    subgraph "Scenario: Payment Service Down"
        A[Order Service] -->|✅ Still Working| B[Create Orders]
        C[Payment Service] -->|❌ Down| D[No Payment Processing]
        E[Notification Service] -->|✅ Still Working| F[Send Order Confirmations]
    end
    
    style A fill:#e1ffe1
    style C fill:#ffe1e1
    style E fill:#e1ffe1
```

**Isolation Mechanisms:**
- Separate processes
- Independent lifecycles
- No shared state
- Event-driven communication

### 4. Scalability

Services can scale independently based on load:

```mermaid
graph TB
    subgraph "Horizontal Scaling"
        LB[Load Balancer]
        O1[Order Service 1]
        O2[Order Service 2]
        O3[Order Service 3]
        
        LB --> O1
        LB --> O2
        LB --> O3
    end
    
    subgraph "Event Bus"
        EB[Event Bus]
    end
    
    O1 --> EB
    O2 --> EB
    O3 --> EB
    
    style LB fill:#fff4e1
    style EB fill:#fff4e1
```

**Scaling Strategies:**
- **Horizontal**: Add more instances
- **Vertical**: Increase resources per instance
- **Independent**: Scale each service based on its needs

## Go-Specific Features

### Goroutines for Concurrency

```mermaid
graph LR
    A[Main Thread] -->|spawn| B[Goroutine 1]
    A -->|spawn| C[Goroutine 2]
    A -->|spawn| D[Goroutine 3]
    A -->|spawn| E[Goroutine N]
    
    style A fill:#e1f5ff
    style B fill:#e1ffe1
    style C fill:#e1ffe1
    style D fill:#e1ffe1
    style E fill:#e1ffe1
```

**Implementation:**
```go
// Event handlers run in separate goroutines
for _, handler := range handlers {
    go func(h Handler) {
        h(event)  // Non-blocking!
    }(handler)
}
```

**Benefits:**
- Lightweight (2KB stack vs 1MB thread)
- Can handle thousands of concurrent events
- Built-in scheduler
- No manual thread management

### Channels for Communication

```go
// Safe communication between goroutines
events := make(chan Event, 100)

// Producer
go func() {
    events <- newEvent
}()

// Consumer
go func() {
    event := <-events
    process(event)
}()
```

**Benefits:**
- No race conditions
- No shared memory issues
- Clear data flow
- Type-safe

### sync.RWMutex for Thread Safety

```go
type OrderRepository struct {
    orders map[string]*Order
    mu     sync.RWMutex  // Concurrent safety
}

func (r *OrderRepository) Save(order *Order) error {
    r.mu.Lock()         // Exclusive write lock
    defer r.mu.Unlock()
    r.orders[order.ID] = order
    return nil
}

func (r *OrderRepository) FindByID(id string) (*Order, error) {
    r.mu.RLock()        // Shared read lock
    defer r.mu.RUnlock()
    return r.orders[id], nil
}
```

**Benefits:**
- Multiple readers, single writer
- Prevents race conditions
- Efficient concurrent access

## Production Considerations

### What's Simplified (Demo)

```mermaid
graph LR
    subgraph "Demo Implementation"
        A[In-Memory Event Bus]
        B[In-Memory Storage]
        C[Hardcoded Ports]
    end
    
    subgraph "Production Implementation"
        D[Kafka/RabbitMQ]
        E[PostgreSQL/MongoDB]
        F[Service Discovery]
    end
    
    A -.->|Replace with| D
    B -.->|Replace with| E
    C -.->|Replace with| F
    
    style A fill:#ffe1e1
    style B fill:#ffe1e1
    style C fill:#ffe1e1
    style D fill:#e1ffe1
    style E fill:#e1ffe1
    style F fill:#e1ffe1
```

### What's Production-Ready

✅ **Graceful Shutdown** - Proper signal handling
✅ **Structured Logging** - Clear log messages
✅ **Error Handling** - Proper error propagation
✅ **HTTP Timeouts** - Read/Write/Idle timeouts
✅ **Concurrent Safety** - Thread-safe data structures
✅ **Health Checks** - `/health` endpoints
✅ **Dockerization** - Multi-stage builds

## System Behavior Examples

### Success Flow

```mermaid
stateDiagram-v2
    [*] --> OrderReceived: Client POST /orders
    OrderReceived --> OrderSaved: Save to repository
    OrderSaved --> EventPublished: Publish OrderCreated
    EventPublished --> PaymentStarted: Payment Service receives event
    EventPublished --> NotificationSent: Notification Service receives event
    PaymentStarted --> PaymentCompleted: Process payment
    PaymentCompleted --> PaymentEventPublished: Publish PaymentProcessed
    PaymentEventPublished --> PaymentNotificationSent: Notification Service receives event
    PaymentNotificationSent --> [*]: Complete
```

### Failure Handling

```mermaid
stateDiagram-v2
    [*] --> OrderReceived: Client POST /orders
    OrderReceived --> OrderSaved: Save to repository
    OrderSaved --> EventPublished: Publish OrderCreated
    EventPublished --> PaymentStarted: Payment Service receives event
    PaymentStarted --> PaymentFailed: Payment declined
    PaymentFailed --> PaymentEventPublished: Publish PaymentProcessed (failed)
    PaymentEventPublished --> FailureNotificationSent: Notification Service receives event
    FailureNotificationSent --> [*]: Customer notified of failure
```

## Extending the System

### Adding a New Service

```mermaid
graph TB
    A[Define Service Responsibility] --> B[Create Service Directory]
    B --> C[Subscribe to Relevant Events]
    C --> D[Implement Business Logic]
    D --> E[Publish New Events if needed]
    E --> F[Add to docker-compose.yml]
    F --> G[Test Integration]
    
    style A fill:#e1f5ff
    style G fill:#e1ffe1
```

### Adding a New Event Type

```mermaid
graph LR
    A[Define Event in types.go] --> B[Publish from Producer]
    B --> C[Subscribe in Consumers]
    C --> D[Implement Handlers]
    
    style A fill:#e1f5ff
    style D fill:#e1ffe1
```

## Conclusion

This system demonstrates how distributed system patterns work together:

- **Layered Architecture** provides structure within services
- **Client-Server** enables synchronous communication
- **Event-Driven** enables asynchronous communication
- **Microservices** enables independent scaling and deployment

The implementation is educational but follows production best practices, making it a solid foundation for understanding distributed systems.
