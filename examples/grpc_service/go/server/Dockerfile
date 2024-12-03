FROM golang:1.22.3

WORKDIR /app

COPY . .

RUN go mod download

EXPOSE 8080

ENV PORT=8080
ENV LISTEN_ADDR=0.0.0.0

CMD ["go", "run", "."]
