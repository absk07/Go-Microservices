# base go image
FROM golang:1.22.5-alpine3.20 AS builder

WORKDIR /app

COPY ./listener-service /app

RUN go build -o listenerApp

RUN chmod +x /app/listenerApp

# build a tiny docker image
FROM alpine:3.20

WORKDIR /app

COPY --from=builder /app/listenerApp /app

EXPOSE 5050

CMD ["/app/listenerApp"]