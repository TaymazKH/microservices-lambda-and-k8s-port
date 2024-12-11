FROM golang:1.22.3@sha256:f43c6f049f04cbbaeb28f0aad3eea15274a7d0a7899a617d0037aec48d7ab010

WORKDIR /app

COPY . .

RUN go mod download
RUN GOOS=linux go build -o product_catalog_service .

EXPOSE 3550

ENV PORT=3550
ENV LISTEN_ADDR=0.0.0.0

ENTRYPOINT ["/app/product_catalog_service"]
