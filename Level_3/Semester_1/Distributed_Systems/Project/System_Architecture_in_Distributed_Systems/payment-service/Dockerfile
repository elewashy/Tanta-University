# Build stage
FROM golang:1.21-alpine AS builder

WORKDIR /app

# Copy shared module first
COPY shared/ /app/shared/

# Copy service files
COPY payment-service/go.mod payment-service/go.sum* ./
RUN go mod download

COPY payment-service/ ./

# Build static binary
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .

# Runtime stage
FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /root/

COPY --from=builder /app/main .

EXPOSE 8081

CMD ["./main"]
