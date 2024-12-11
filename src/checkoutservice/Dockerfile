FROM golang:1.22.3@sha256:f43c6f049f04cbbaeb28f0aad3eea15274a7d0a7899a617d0037aec48d7ab010

WORKDIR /app

COPY . .

RUN go mod download
RUN GOOS=linux go build -o checkout_service .

EXPOSE 5050

ENV PORT=5050
ENV LISTEN_ADDR=0.0.0.0
ENV CART_SERVICE_ADDR="http://my-cart-service:7070"
ENV CURRENCY_SERVICE_ADDR="http://my-currency-service:7000"
ENV EMAIL_SERVICE_ADDR="http://my-email-service:8081"
ENV PAYMENT_SERVICE_ADDR="http://my-payment-service:50052"
ENV PRODUCT_CATALOG_SERVICE_ADDR="http://my-product-catalog-service:3550"
ENV SHIPPING_SERVICE_ADDR="http://my-shipping-service:50053"

ENTRYPOINT ["/app/checkout_service"]
