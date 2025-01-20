FROM node:20.18.0-slim AS builder

WORKDIR /app

COPY package*.json ./
RUN npm install --production

FROM alpine:3.21.2

RUN apk add --no-cache nodejs

WORKDIR /app

COPY --from=builder /app/node_modules ./node_modules
COPY . .

ENV PORT=7000
ENV LISTEN_ADDR=0.0.0.0

EXPOSE 7000
ENTRYPOINT ["node", "server.js"]
