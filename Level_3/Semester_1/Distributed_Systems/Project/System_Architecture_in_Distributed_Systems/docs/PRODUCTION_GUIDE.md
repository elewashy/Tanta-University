# Production Deployment Guide

## From Demo to Production

This guide explains how to evolve this demo system into a production-ready distributed system.

## Current State vs Production

| Component | Demo Implementation | Production Implementation |
|-----------|-------------------|--------------------------|
| **Event Bus** | In-memory channels | Kafka, RabbitMQ, NATS, AWS SNS/SQS |
| **Database** | In-memory map | PostgreSQL, MongoDB, Cassandra |
| **Service Discovery** | Hardcoded ports | Consul, etcd, Kubernetes DNS |
| **API Gateway** | None | Kong, Traefik, AWS API Gateway |
| **Authentication** | None | JWT, OAuth2, mTLS |
| **Monitoring** | Basic logs | Prometheus, Grafana, ELK Stack |
| **Tracing** | None | Jaeger, Zipkin, AWS X-Ray |
| **Configuration** | Hardcoded | Consul, etcd, AWS Parameter Store |
| **Secrets** | None | Vault, AWS Secrets Manager |
| **Load Balancing** | None | Nginx, HAProxy, AWS ALB |

## Step-by-Step Production Evolution

### Phase 1: Add Real Message Broker

#### Replace In-Memory Event Bus with Kafka

**Install Kafka:**
```yaml
# docker-compose.yml
services:
  zookeeper:
    image: confluentinc/cp-zookeeper:latest
    environment:
      ZOOKEEPER_CLIENT_PORT: 2181

  kafka:
    image: confluentinc/cp-kafka:latest
    depends_on:
      - zookeeper
    environment:
      KAFKA_ZOOKEEPER_CONNECT: zookeeper:2181
      KAFKA_ADVERTISED_LISTENERS: PLAINTEXT://kafka:9092
```

**Update Event Bus:**
```go
// shared/events/kafka_bus.go
package events

import (
    "context"
    "encoding/json"
    "github.com/segmentio/kafka-go"
)

type KafkaEventBus struct {
    writer *kafka.Writer
    readers map[EventType]*kafka.Reader
}

func NewKafkaEventBus(brokers []string) *KafkaEventBus {
    return &KafkaEventBus{
        writer: &kafka.Writer{
            Addr:     kafka.TCP(brokers...),
            Balancer: &kafka.LeastBytes{},
        },
        readers: make(map[EventType]*kafka.Reader),
    }
}

func (kb *KafkaEventBus) Publish(event Event) error {
    data, _ := json.Marshal(event)
    return kb.writer.WriteMessages(context.Background(),
        kafka.Message{
            Topic: string(event.Type),
            Value: data,
        },
    )
}

func (kb *KafkaEventBus) Subscribe(eventType EventType, handler Handler) {
    reader := kafka.NewReader(kafka.ReaderConfig{
        Brokers: []string{"localhost:9092"},
        Topic:   string(eventType),
        GroupID: "consumer-group",
    })
    
    go func() {
        for {
            msg, _ := reader.ReadMessage(context.Background())
            var event Event
            json.Unmarshal(msg.Value, &event)
            handler(event)
        }
    }()
}
```

### Phase 2: Add Database Persistence

#### Replace In-Memory Repository with PostgreSQL

**Install PostgreSQL:**
```yaml
# docker-compose.yml
services:
  postgres:
    image: postgres:15
    environment:
      POSTGRES_DB: orders
      POSTGRES_USER: admin
      POSTGRES_PASSWORD: secret
    ports:
      - "5432:5432"
```

**Update Repository:**
```go
// order-service/repository/postgres_repository.go
package repository

import (
    "database/sql"
    _ "github.com/lib/pq"
)

type PostgresOrderRepository struct {
    db *sql.DB
}

func NewPostgresOrderRepository(connStr string) (*PostgresOrderRepository, error) {
    db, err := sql.Open("postgres", connStr)
    if err != nil {
        return nil, err
    }
    
    // Create table
    _, err = db.Exec(`
        CREATE TABLE IF NOT EXISTS orders (
            id VARCHAR(255) PRIMARY KEY,
            customer_id VARCHAR(255),
            items TEXT[],
            total DECIMAL(10,2),
            status VARCHAR(50)
        )
    `)
    
    return &PostgresOrderRepository{db: db}, err
}

func (r *PostgresOrderRepository) Save(order *Order) error {
    _, err := r.db.Exec(
        "INSERT INTO orders (id, customer_id, items, total, status) VALUES ($1, $2, $3, $4, $5)",
        order.ID, order.CustomerID, pq.Array(order.Items), order.Total, order.Status,
    )
    return err
}

func (r *PostgresOrderRepository) FindByID(id string) (*Order, error) {
    var order Order
    err := r.db.QueryRow(
        "SELECT id, customer_id, items, total, status FROM orders WHERE id = $1",
        id,
    ).Scan(&order.ID, &order.CustomerID, pq.Array(&order.Items), &order.Total, &order.Status)
    
    return &order, err
}
```

### Phase 3: Add Service Discovery

#### Use Consul for Service Discovery

**Install Consul:**
```yaml
# docker-compose.yml
services:
  consul:
    image: consul:latest
    ports:
      - "8500:8500"
```

**Register Service:**
```go
// shared/discovery/consul.go
package discovery

import (
    "github.com/hashicorp/consul/api"
)

func RegisterService(name string, port int) error {
    config := api.DefaultConfig()
    client, _ := api.NewClient(config)
    
    registration := &api.AgentServiceRegistration{
        ID:      name,
        Name:    name,
        Port:    port,
        Address: "localhost",
        Check: &api.AgentServiceCheck{
            HTTP:     fmt.Sprintf("http://localhost:%d/health", port),
            Interval: "10s",
        },
    }
    
    return client.Agent().ServiceRegister(registration)
}

func DiscoverService(name string) (string, error) {
    config := api.DefaultConfig()
    client, _ := api.NewClient(config)
    
    services, _, _ := client.Health().Service(name, "", true, nil)
    if len(services) > 0 {
        return fmt.Sprintf("%s:%d", 
            services[0].Service.Address, 
            services[0].Service.Port), nil
    }
    
    return "", errors.New("service not found")
}
```

### Phase 4: Add API Gateway

#### Use Kong or Custom Gateway

**Custom API Gateway:**
```go
// api-gateway/main.go
package main

import (
    "net/http"
    "net/http/httputil"
    "net/url"
)

func main() {
    // Order Service proxy
    orderURL, _ := url.Parse("http://localhost:8080")
    orderProxy := httputil.NewSingleHostReverseProxy(orderURL)
    
    http.HandleFunc("/api/orders", func(w http.ResponseWriter, r *http.Request) {
        // Add authentication
        if !authenticate(r) {
            http.Error(w, "Unauthorized", http.StatusUnauthorized)
            return
        }
        
        // Add rate limiting
        if !checkRateLimit(r) {
            http.Error(w, "Too Many Requests", http.StatusTooManyRequests)
            return
        }
        
        // Forward to service
        orderProxy.ServeHTTP(w, r)
    })
    
    http.ListenAndServe(":8000", nil)
}
```

### Phase 5: Add Monitoring & Observability

#### Prometheus Metrics

**Add Metrics:**
```go
// shared/metrics/prometheus.go
package metrics

import (
    "github.com/prometheus/client_golang/prometheus"
    "github.com/prometheus/client_golang/prometheus/promauto"
)

var (
    OrdersCreated = promauto.NewCounter(prometheus.CounterOpts{
        Name: "orders_created_total",
        Help: "Total number of orders created",
    })
    
    OrderDuration = promauto.NewHistogram(prometheus.HistogramOpts{
        Name: "order_duration_seconds",
        Help: "Order processing duration",
    })
)

// In handler
func (h *OrderHandler) CreateOrder(w http.ResponseWriter, r *http.Request) {
    timer := prometheus.NewTimer(metrics.OrderDuration)
    defer timer.ObserveDuration()
    
    // ... create order logic ...
    
    metrics.OrdersCreated.Inc()
}
```

**Expose Metrics:**
```go
// Add to main.go
import "github.com/prometheus/client_golang/prometheus/promhttp"

http.Handle("/metrics", promhttp.Handler())
```

#### Distributed Tracing with Jaeger

**Add Tracing:**
```go
// shared/tracing/jaeger.go
package tracing

import (
    "github.com/opentracing/opentracing-go"
    "github.com/uber/jaeger-client-go"
    "github.com/uber/jaeger-client-go/config"
)

func InitTracer(serviceName string) (opentracing.Tracer, io.Closer) {
    cfg := config.Configuration{
        ServiceName: serviceName,
        Sampler: &config.SamplerConfig{
            Type:  jaeger.SamplerTypeConst,
            Param: 1,
        },
        Reporter: &config.ReporterConfig{
            LogSpans: true,
            LocalAgentHostPort: "localhost:6831",
        },
    }
    
    tracer, closer, _ := cfg.NewTracer()
    return tracer, closer
}

// In handler
func (h *OrderHandler) CreateOrder(w http.ResponseWriter, r *http.Request) {
    span := opentracing.StartSpan("create_order")
    defer span.Finish()
    
    // ... create order logic ...
}
```

### Phase 6: Add Resilience Patterns

#### Circuit Breaker

```go
// shared/resilience/circuit_breaker.go
package resilience

import (
    "github.com/sony/gobreaker"
)

var cb *gobreaker.CircuitBreaker

func init() {
    cb = gobreaker.NewCircuitBreaker(gobreaker.Settings{
        Name:        "payment-service",
        MaxRequests: 3,
        Interval:    60,
        Timeout:     30,
    })
}

func CallWithCircuitBreaker(fn func() error) error {
    _, err := cb.Execute(func() (interface{}, error) {
        return nil, fn()
    })
    return err
}
```

#### Retry with Exponential Backoff

```go
// shared/resilience/retry.go
package resilience

import (
    "time"
    "github.com/cenkalti/backoff/v4"
)

func RetryWithBackoff(operation func() error) error {
    exponentialBackoff := backoff.NewExponentialBackOff()
    exponentialBackoff.MaxElapsedTime = 5 * time.Minute
    
    return backoff.Retry(operation, exponentialBackoff)
}
```

### Phase 7: Add Security

#### JWT Authentication

```go
// shared/auth/jwt.go
package auth

import (
    "github.com/golang-jwt/jwt/v5"
)

func ValidateToken(tokenString string) (*jwt.Token, error) {
    return jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
        return []byte("your-secret-key"), nil
    })
}

// Middleware
func AuthMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        token := r.Header.Get("Authorization")
        if token == "" {
            http.Error(w, "Unauthorized", http.StatusUnauthorized)
            return
        }
        
        _, err := ValidateToken(token)
        if err != nil {
            http.Error(w, "Invalid token", http.StatusUnauthorized)
            return
        }
        
        next.ServeHTTP(w, r)
    })
}
```

## Kubernetes Deployment

### Deployment Manifests

```yaml
# k8s/order-service-deployment.yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: order-service
spec:
  replicas: 3
  selector:
    matchLabels:
      app: order-service
  template:
    metadata:
      labels:
        app: order-service
    spec:
      containers:
      - name: order-service
        image: order-service:latest
        ports:
        - containerPort: 8080
        env:
        - name: DATABASE_URL
          valueFrom:
            secretKeyRef:
              name: db-secret
              key: url
        livenessProbe:
          httpGet:
            path: /health
            port: 8080
          initialDelaySeconds: 30
          periodSeconds: 10
        readinessProbe:
          httpGet:
            path: /health
            port: 8080
          initialDelaySeconds: 5
          periodSeconds: 5
---
apiVersion: v1
kind: Service
metadata:
  name: order-service
spec:
  selector:
    app: order-service
  ports:
  - port: 80
    targetPort: 8080
  type: LoadBalancer
```

## Production Checklist

- [ ] Replace in-memory event bus with Kafka/RabbitMQ
- [ ] Add database persistence (PostgreSQL/MongoDB)
- [ ] Implement service discovery (Consul/Kubernetes)
- [ ] Add API Gateway (Kong/Traefik)
- [ ] Implement authentication (JWT/OAuth2)
- [ ] Add monitoring (Prometheus/Grafana)
- [ ] Add distributed tracing (Jaeger/Zipkin)
- [ ] Implement circuit breaker pattern
- [ ] Add retry logic with exponential backoff
- [ ] Implement rate limiting
- [ ] Add request/response validation
- [ ] Implement proper error handling
- [ ] Add structured logging (Zap/Logrus)
- [ ] Implement health checks
- [ ] Add graceful shutdown (✓ Already implemented)
- [ ] Configure timeouts (✓ Already implemented)
- [ ] Add CORS handling
- [ ] Implement request ID tracking
- [ ] Add database migrations
- [ ] Implement backup strategy
- [ ] Add disaster recovery plan
- [ ] Configure auto-scaling
- [ ] Implement blue-green deployment
- [ ] Add canary deployment strategy
- [ ] Configure SSL/TLS
- [ ] Implement secrets management
- [ ] Add compliance logging (GDPR, etc.)
- [ ] Implement data encryption at rest
- [ ] Add DDoS protection
- [ ] Configure CDN for static assets
- [ ] Implement caching strategy (Redis)
- [ ] Add performance testing
- [ ] Implement chaos engineering tests

## Recommended Tools & Libraries

### Go Libraries
- **HTTP Router**: gorilla/mux, gin, echo
- **Database**: gorm, sqlx, pgx
- **Messaging**: kafka-go, amqp, nats.go
- **Monitoring**: prometheus/client_golang
- **Tracing**: opentracing-go, jaeger-client-go
- **Logging**: zap, logrus
- **Configuration**: viper, envconfig
- **Validation**: validator/v10
- **Testing**: testify, gomock

### Infrastructure
- **Container Orchestration**: Kubernetes, Docker Swarm
- **Service Mesh**: Istio, Linkerd
- **CI/CD**: GitHub Actions, GitLab CI, Jenkins
- **Monitoring**: Prometheus, Grafana, Datadog
- **Logging**: ELK Stack, Loki, Splunk
- **Tracing**: Jaeger, Zipkin, AWS X-Ray
- **Message Broker**: Kafka, RabbitMQ, NATS, AWS SQS
- **Database**: PostgreSQL, MongoDB, Cassandra
- **Cache**: Redis, Memcached
- **API Gateway**: Kong, Traefik, AWS API Gateway

## Cost Optimization

1. **Right-size resources** - Don't over-provision
2. **Use auto-scaling** - Scale based on demand
3. **Implement caching** - Reduce database load
4. **Use spot instances** - For non-critical workloads
5. **Optimize database queries** - Add indexes, use connection pooling
6. **Compress responses** - Reduce bandwidth costs
7. **Use CDN** - Cache static content closer to users
8. **Monitor costs** - Set up billing alerts

## Conclusion

This demo provides a solid foundation. The architecture patterns remain the same in production - you're just replacing the implementations with production-grade tools while maintaining the same distributed system principles.

Start small, iterate, and evolve based on your actual needs!
