# order-service/Dockerfile
FROM golang:1.23-alpine AS builder
WORKDIR /app

# Установи git — нужен для go get из GitHub
RUN apk --no-cache add git

# Копируем go.mod и go.sum для кэширования зависимостей
COPY go.mod go.sum ./
RUN go mod download

# Копируем остальные файлы
COPY . .

# Собираем
RUN go build -o order-srv .

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app/order-srv .
EXPOSE 8081
CMD ["./order-srv"]
