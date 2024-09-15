# base go image
FROM golang:1.22.5-alpine3.20 AS builder

WORKDIR /app

COPY ./logger-service /app

RUN go build -o loggerApp ./src/api

RUN chmod +x /app/loggerApp

# build a tiny docker image
FROM alpine:3.20

WORKDIR /app

COPY --from=builder /app/loggerApp /app

EXPOSE 7070

CMD ["/app/loggerApp"]