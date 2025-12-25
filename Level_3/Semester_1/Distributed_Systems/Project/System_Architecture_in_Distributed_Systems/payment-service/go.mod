module payment-service

go 1.21

replace shared => ../shared

require (
	github.com/google/uuid v1.6.0
	github.com/gorilla/mux v1.8.1
	shared v0.0.0
)
