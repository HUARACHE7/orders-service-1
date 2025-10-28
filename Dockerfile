FROM golang:1.25-alpine AS builder

WORKDIR /app

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . .

RUN go build -o /orders-service main.go


FROM alpine:latest
WORKDIR /root/

COPY --from=builder /orders-service .
COPY static static
COPY model.json model.json

CMD ["./orders-service"]