# Distributed System Demo in Go

## 🎯 Overview

This project demonstrates core distributed system architectural patterns:
- **Layered Architecture** - Clean separation of concerns (API → Service → Repository)
- **Client-Server** - HTTP-based request/response communication
- **Event-Driven Architecture** - Asynchronous event publishing and consumption
- **Microservices** - Independent, loosely-coupled services

## 📋 Table of Contents

- [Why Go for Distributed Systems?](#-why-go-for-distributed-systems)
- [Architecture](#-architecture)
- [Quick Start](#-quick-start)
- [Testing the System](#-testing-the-system)
- [Architecture Mapping](#-architecture-mapping)
- [Docker Support](#-docker-support)
- [Production Features](#-production-features)
- [Documentation](#-documentation)
- [Next Steps](#-next-steps)

## 🔧 Why Go for Distributed Systems?

Go is ideal for distributed systems because:

1. **Goroutines & Channels** - Lightweight concurrency primitives for handling thousands of concurrent connections
2. **Static Binaries** - Single executable with no runtime dependencies, perfect for containerization
3. **Fast Compilation** - Quick iteration during development
4. **Built-in HTTP Server** - Production-ready `net/http` package
5. **Strong Standard Library** - JSON, networking, and concurrency support out of the box
6. **Performance** - Near C-level performance with memory safety
7. **Simple Deployment** - No VM or interpreter needed

## 🏗️ Architecture

### System Overview

```mermaid
graph TB
    Client[Client Application]
    
    subgraph "Order Service :8080"
        API[API Layer]
        Service[Service Layer]
        Repo[Repository Layer]
        API --> Service --> Repo
    end
    
    subgraph "Event Bus"
        EventBus[Pub/Sub Event Bus]
    end
    
    subgraph "Payment Service :8081"
        Payment[Payment Processor]
    end
    
    subgraph "Notification Service :8082"
        Notification[Notification Sender]
    end
    
    Client -->|HTTP POST/GET| API
    Service -->|Publish| EventBus
    EventBus -->|OrderCreated| Payment
    EventBus -->|OrderCreated| Notification
    Payment -->|PaymentProcessed| EventBus
    EventBus -->|PaymentProcessed| Notification
    
    style EventBus fill:#ffd700,stroke:#333,stroke-width:2px,color:#000
    style Client fill:#4a9eff,stroke:#333,stroke-width:2px,color:#fff
```

### Services

1. **Order Service** (Port 8080)
   - Layered architecture: API → Service → Repository
   - Manages order creation and retrieval
   - Publishes `OrderCreated` events

2. **Payment Service** (Port 8081)
   - Subscribes to `OrderCreated` events
   - Processes payments asynchronously
   - Publishes `PaymentProcessed` events

3. **Notification Service** (Port 8082)
   - Subscribes to `OrderCreated` and `PaymentProcessed` events
   - Sends notifications (simulated)

### Event Flow

```mermaid
sequenceDiagram
    participant C as Client
    participant O as Order Service
    participant E as Event Bus
    participant P as Payment Service
    participant N as Notification Service
    
    C->>O: POST /orders
    O->>O: Create & Save Order
    O->>E: Publish OrderCreated
    O-->>C: 201 Created
    
    E-->>P: OrderCreated Event
    E-->>N: OrderCreated Event
    
    P->>P: Process Payment
    P->>E: Publish PaymentProcessed
    N->>N: Send Order Confirmation
    
    E-->>N: PaymentProcessed Event
    N->>N: Send Payment Notification
```

### Distributed Concepts Demonstrated

```mermaid
graph LR
    A[Loose Coupling] --> B[Services communicate via events]
    C[Async Communication] --> D[Non-blocking event processing]
    E[Fault Isolation] --> F[Independent service failures]
    G[Scalability] --> H[Horizontal scaling per service]
    I[Service Boundaries] --> J[Clear responsibilities]
    
    style A fill:#4a9eff,stroke:#333,stroke-width:2px,color:#fff
    style C fill:#ffd700,stroke:#333,stroke-width:2px,color:#000
    style E fill:#ff6b6b,stroke:#333,stroke-width:2px,color:#fff
    style G fill:#51cf66,stroke:#333,stroke-width:2px,color:#000
    style I fill:#a78bfa,stroke:#333,stroke-width:2px,color:#fff
```

- **Loose Coupling** - Services communicate via events, not direct calls
- **Asynchronous Communication** - Event-driven processing
- **Fault Isolation** - Each service can fail independently
- **Scalability** - Services can be scaled horizontally
- **Service Boundaries** - Clear separation of responsibilities

## 🚀 Running the System

### Prerequisites
- Go 1.21+ installed

### Start All Services

```bash
# Terminal 1 - Order Service
cd order-service
go run main.go

# Terminal 2 - Payment Service
cd payment-service
go run main.go

# Terminal 3 - Notification Service
cd notification-service
go run main.go
```

### Test the System

```bash
# Create an order
curl -X POST http://localhost:8080/orders \
  -H "Content-Type: application/json" \
  -d '{"customer_id":"cust-123","items":["item1","item2"],"total":99.99}'

# Get order by ID
curl http://localhost:8080/orders/{order-id}
```

Watch the logs in each terminal to see the event flow!

## 📊 Architecture Mapping

```mermaid
graph TB
    subgraph "Architectural Patterns"
        L[Layered Architecture]
        CS[Client-Server]
        ED[Event-Driven]
        MS[Microservices]
    end
    
    subgraph "Implementation"
        L1[Order Service: API → Service → Repository]
        CS1[HTTP REST Endpoints]
        ED1[Event Bus + Pub/Sub]
        MS1[3 Independent Services]
    end
    
    L -.-> L1
    CS -.-> CS1
    ED -.-> ED1
    MS -.-> MS1
    
    style L fill:#4a9eff,stroke:#333,stroke-width:2px,color:#fff
    style CS fill:#ffd700,stroke:#333,stroke-width:2px,color:#000
    style ED fill:#ff6b6b,stroke:#333,stroke-width:2px,color:#fff
    style MS fill:#51cf66,stroke:#333,stroke-width:2px,color:#000
```

| Pattern | Implementation |
|---------|---------------|
| **Layered Architecture** | Order Service: API → Service → Repository layers |
| **Client-Server** | HTTP REST endpoints for synchronous communication |
| **Event-Driven** | In-memory event bus for async communication |
| **Microservices** | Three independent services with clear boundaries |

## 🐳 Docker Support (Bonus)

Each service includes a Dockerfile for containerization:

```bash
# Build and run with Docker
docker build -t order-service ./order-service
docker run -p 8080:8080 order-service
```

## 🔄 Production Features

- **Graceful Shutdown** - Services handle SIGTERM/SIGINT properly
- **Structured Logging** - Clear log output for debugging
- **Error Handling** - Proper error propagation and HTTP status codes
- **Concurrent Safety** - Thread-safe data structures (sync.RWMutex)

## 📝 Notes

This demo uses an **in-memory event bus** for simplicity. In production, you would use:

```mermaid
graph LR
    Demo[In-Memory Event Bus] -.->|Replace with| Prod1[Kafka]
    Demo -.->|Replace with| Prod2[RabbitMQ]
    Demo -.->|Replace with| Prod3[NATS]
    Demo -.->|Replace with| Prod4[Redis Streams]
    
    style Demo fill:#ff6b6b,stroke:#333,stroke-width:2px,color:#fff
    style Prod1 fill:#51cf66,stroke:#333,stroke-width:2px,color:#000
    style Prod2 fill:#51cf66,stroke:#333,stroke-width:2px,color:#000
    style Prod3 fill:#51cf66,stroke:#333,stroke-width:2px,color:#000
    style Prod4 fill:#51cf66,stroke:#333,stroke-width:2px,color:#000
```

- **Kafka** - High-throughput, distributed event streaming
- **RabbitMQ** - Message queuing with routing
- **NATS** - Lightweight, cloud-native messaging
- **Redis Streams** - Simple pub/sub with persistence

The architecture principles remain the same regardless of the messaging technology.

## 📚 Documentation

```mermaid
graph TB
    Start([Start Here]) --> README[📖 README.md<br/>Quick Overview]
    
    README --> Choice{What Next?}
    
    Choice -->|Learn Architecture| ARCH[🏗️ ARCHITECTURE.md<br/>Detailed Patterns & Diagrams]
    Choice -->|Get Started| QUICK[🚀 QUICKSTART.md<br/>Step-by-Step Guide]
    Choice -->|Understand Go| WHY[💡 WHY_GO.md<br/>Go for Distributed Systems]
    Choice -->|Go to Production| PROD[🏭 PRODUCTION_GUIDE.md<br/>Production Evolution]
    
    style Start fill:#4a9eff,stroke:#333,stroke-width:2px,color:#fff
    style README fill:#ffd700,stroke:#333,stroke-width:2px,color:#000
    style ARCH fill:#ff6b6b,stroke:#333,stroke-width:2px,color:#fff
    style QUICK fill:#51cf66,stroke:#333,stroke-width:2px,color:#000
    style WHY fill:#a78bfa,stroke:#333,stroke-width:2px,color:#fff
    style PROD fill:#ff6bcf,stroke:#333,stroke-width:2px,color:#fff
    style EX fill:#20c997,stroke:#333,stroke-width:2px,color:#fff
```

| Document | Purpose |
|----------|---------|
| **[README.md](README.md)** | This file - overview and quick start |
| **[docs/ARCHITECTURE.md](docs/ARCHITECTURE.md)** | Detailed architecture with Mermaid diagrams |
| **[docs/QUICKSTART.md](docs/QUICKSTART.md)** | Step-by-step getting started guide |
| **[docs/WHY_GO.md](docs/WHY_GO.md)** | Why Go is ideal for distributed systems |
| **[docs/PRODUCTION_GUIDE.md](docs/PRODUCTION_GUIDE.md)** | Evolving to production |

## 🎓 Next Steps

1. **Run the System** - Follow the [Quick Start](#-running-the-system) above
2. **Understand Architecture** - Read [docs/ARCHITECTURE.md](docs/ARCHITECTURE.md)
4. **Explore Production** - Learn production patterns in [docs/PRODUCTION_GUIDE.md](docs/PRODUCTION_GUIDE.md)

## 🌟 Key Features

✅ **Educational but Realistic** - Production-style code, not toy examples
✅ **Clear Architecture** - Each pattern clearly demonstrated
✅ **Well Documented** - Comprehensive guides with visual diagrams
✅ **Extensible** - Easy to add new services and features
✅ **Production Ready** - Graceful shutdown, health checks, error handling

---

**Happy Learning! 🚀**

For detailed information, explore the [documentation](#-documentation) above.
