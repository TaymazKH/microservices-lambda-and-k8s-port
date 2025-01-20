FROM node:20.18.0-slim AS builder

WORKDIR /app

COPY package*.json ./
RUN npm install --production

FROM alpine:latest

RUN apk add --no-cache nodejs

WORKDIR /app

COPY --from=builder /app/node_modules ./node_modules
COPY . .

ENV PORT=8080
ENV LISTEN_ADDR=0.0.0.0

EXPOSE 8080
ENTRYPOINT ["node", "server.js"]
