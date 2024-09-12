# base go image
FROM golang:1.22.5-alpine3.20 AS builder

WORKDIR /app

COPY ./auth-service /app

RUN go build -o authApp ./src/api

RUN chmod +x /app/authApp

# build a tiny docker image
FROM alpine:3.20

WORKDIR /app

COPY --from=builder /app/authApp /app

EXPOSE 9090

CMD ["/app/authApp"]