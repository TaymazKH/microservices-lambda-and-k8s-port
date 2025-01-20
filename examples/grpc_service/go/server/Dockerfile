FROM golang:1.22.3 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o greeter .

FROM alpine:3.21.2

WORKDIR /app

RUN apk add --no-cache ca-certificates

COPY --from=builder /app/greeter .
RUN chmod +x greeter

ENV PORT=8080
ENV LISTEN_ADDR=0.0.0.0

EXPOSE 8080
ENTRYPOINT ["/app/greeter"]
