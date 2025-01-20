FROM golang:1.22.3 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o frontend_service .

FROM alpine:3.21.2

WORKDIR /app

RUN apk add --no-cache ca-certificates

COPY --from=builder /app/frontend_service ./
RUN chmod +x frontend_service

COPY ./templates ./templates
COPY ./static ./static

ENV PORT=8080
ENV LISTEN_ADDR=0.0.0.0
ENV AD_SERVICE_ADDR="http://my-ad-service:9555"
ENV CART_SERVICE_ADDR="http://my-cart-service:7070"
ENV CHECKOUT_SERVICE_ADDR="http://my-checkout-service:5050"
ENV CURRENCY_SERVICE_ADDR="http://my-currency-service:7000"
ENV PRODUCT_CATALOG_SERVICE_ADDR="http://my-product-catalog-service:3550"
ENV RECOMMENDATION_SERVICE_ADDR="http://my-recommendation-service:8082"
ENV SHIPPING_SERVICE_ADDR="http://my-shipping-service:50053"

EXPOSE 8080
ENTRYPOINT ["/app/frontend_service"]
