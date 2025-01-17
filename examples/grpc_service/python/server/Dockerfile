FROM python:3.10@sha256:3ba2e48b887586835af6a0c35fc6fc6086fb4881e963082330ab0a35f3f42c16

WORKDIR /app

COPY . .

RUN pip install -r requirements.txt

EXPOSE 8080

ENV PORT=8080
ENV LISTEN_ADDR=0.0.0.0

ENTRYPOINT ["python", "server.py"]
