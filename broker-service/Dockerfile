# base go image
FROM golang:1.22.5-alpine3.20 AS builder

WORKDIR /app

COPY ./broker-service /app

RUN go build -o brokerApp ./src/api

RUN chmod +x /app/brokerApp

# build a tiny docker image
FROM alpine:3.20

WORKDIR /app

COPY --from=builder /app/brokerApp /app

EXPOSE 8080

CMD ["/app/brokerApp"]