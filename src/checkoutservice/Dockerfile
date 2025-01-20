FROM golang:1.22.3 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o checkout_service .

FROM alpine:3.21.2

WORKDIR /app

RUN apk add --no-cache ca-certificates

COPY --from=builder /app/checkout_service .
RUN chmod +x checkout_service

ENV PORT=5050
ENV LISTEN_ADDR=0.0.0.0
ENV CART_SERVICE_ADDR="http://my-cart-service:7070"
ENV CURRENCY_SERVICE_ADDR="http://my-currency-service:7000"
ENV EMAIL_SERVICE_ADDR="http://my-email-service:8081"
ENV PAYMENT_SERVICE_ADDR="http://my-payment-service:50052"
ENV PRODUCT_CATALOG_SERVICE_ADDR="http://my-product-catalog-service:3550"
ENV SHIPPING_SERVICE_ADDR="http://my-shipping-service:50053"

EXPOSE 5050
ENTRYPOINT ["/app/checkout_service"]
