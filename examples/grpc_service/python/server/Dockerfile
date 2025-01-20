FROM python:3.10-slim AS base

FROM base AS builder

RUN apt-get -qq update \
    && apt-get install -y --no-install-recommends gcc \
    && rm -rf /var/lib/apt/lists/*

WORKDIR /app

COPY requirements.txt .
RUN pip install --no-cache-dir --prefix=/install_dep -r requirements.txt

FROM base

WORKDIR /app

COPY --from=builder /install_dep /usr/local
COPY . .

ENV PORT=8080
ENV LISTEN_ADDR=0.0.0.0

EXPOSE 8080
ENTRYPOINT ["python", "server.py"]
