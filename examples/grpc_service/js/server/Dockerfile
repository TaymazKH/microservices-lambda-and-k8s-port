FROM node:20.18.0

WORKDIR /app

COPY . .

RUN npm install

EXPOSE 8080

ENV PORT=8080
ENV LISTEN_ADDR=0.0.0.0

CMD ["node", "server.js"]
