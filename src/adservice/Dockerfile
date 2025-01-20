FROM golang:1.22.3 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o ad_service .

FROM alpine:3.21.2

WORKDIR /app

RUN apk add --no-cache ca-certificates

COPY --from=builder /app/ad_service .
RUN chmod +x ad_service

ENV PORT=9555
ENV LISTEN_ADDR=0.0.0.0

EXPOSE 9555
ENTRYPOINT ["/app/ad_service"]
