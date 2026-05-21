# kafka-producer/Dockerfile
FROM golang:1.26-alpine AS builder
WORKDIR /app


COPY . .
RUN go mod download

# Генерируем документацию перед сборкой приложения
# RUN swag init

RUN go build -o producer .

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app/producer .
EXPOSE 8085
CMD ["./producer"]
