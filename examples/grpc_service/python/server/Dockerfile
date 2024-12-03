FROM python:3.10

WORKDIR /app

COPY . .

RUN pip install -r requirements.txt

EXPOSE 8080

ENV PORT=8080
ENV LISTEN_ADDR=0.0.0.0

CMD ["python", "server.py"]
