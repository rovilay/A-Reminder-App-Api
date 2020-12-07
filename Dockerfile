FROM golang:latest

WORKDIR /app

ENV PG_HOST=localhost \
    PG_PORT=5432 \
    PG_USER=postgres \
    PG_PASSWORD=postgres \
    DB_NAME=reminder_app_2

COPY . /app

RUN go build -o bin/tmp .

CMD [ "./bin/tmp" ]
