FROM node:20.18.0@sha256:a7a3b7ec6de4b11bb2d673b31de9d28c6da09c557ee65453672c8e4f754c23fc

WORKDIR /app

COPY . .

RUN npm install

EXPOSE 7000

ENV PORT=7000
ENV LISTEN_ADDR=0.0.0.0

ENTRYPOINT ["node", "server.js"]
