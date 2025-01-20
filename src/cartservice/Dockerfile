FROM golang:1.22.3 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o cart_service .

FROM alpine:3.21.2

WORKDIR /app

RUN apk add --no-cache ca-certificates

COPY --from=builder /app/cart_service .
RUN chmod +x cart_service

ENV PORT=7070
ENV LISTEN_ADDR=0.0.0.0

EXPOSE 7070
ENTRYPOINT ["/app/cart_service"]
