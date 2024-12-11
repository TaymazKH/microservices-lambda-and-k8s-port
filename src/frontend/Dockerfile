FROM golang:1.22.3@sha256:f43c6f049f04cbbaeb28f0aad3eea15274a7d0a7899a617d0037aec48d7ab010

WORKDIR /app

COPY . .

RUN go mod download
RUN GOOS=linux go build -o frontend_service .

EXPOSE 8080

ENV PORT=8080
ENV LISTEN_ADDR=0.0.0.0
ENV AD_SERVICE_ADDR="http://my-ad-service:9555"
ENV CART_SERVICE_ADDR="http://my-cart-service:7070"
ENV CHECKOUT_SERVICE_ADDR="http://my-checkout-service:5050"
ENV CURRENCY_SERVICE_ADDR="http://my-currency-service:7000"
ENV PRODUCT_CATALOG_SERVICE_ADDR="http://my-product-catalog-service:3550"
ENV RECOMMENDATION_SERVICE_ADDR="http://my-recommendation-service:8082"
ENV SHIPPING_SERVICE_ADDR="http://my-shipping-service:50053"

ENTRYPOINT ["/app/frontend_service"]
