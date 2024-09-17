# base go image
FROM golang:1.22.5-alpine3.20 AS builder

WORKDIR /app

COPY ./mail-service /app

RUN go build -o mailApp ./src/api

RUN chmod +x /app/mailApp

# build a tiny docker image
FROM alpine:3.20

WORKDIR /app

COPY --from=builder /app/mailApp /app
COPY --from=builder /app/templates /templates

EXPOSE 6060

CMD ["/app/mailApp"]