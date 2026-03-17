FROM golang:1.26.1-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o shipment-service cmd/main.go

FROM alpine:latest

WORKDIR /app
COPY --from=builder /app/shipment-service .
COPY --from=builder /app/migrations ./migrations

CMD ["./shipment-service"]
