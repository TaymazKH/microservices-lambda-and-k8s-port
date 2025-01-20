FROM golang:1.22.3 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o product_catalog_service .

FROM alpine:3.21.2

WORKDIR /app

RUN apk add --no-cache ca-certificates

COPY --from=builder /app/product_catalog_service .
RUN chmod +x product_catalog_service

COPY products.json .

ENV PORT=3550
ENV LISTEN_ADDR=0.0.0.0

EXPOSE 3550
ENTRYPOINT ["/app/product_catalog_service"]
