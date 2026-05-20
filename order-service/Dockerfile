# order-service/Dockerfile
FROM golang:1.21-alpine AS builder
WORKDIR /app

# Копируем go.mod и go.sum сначала — для кэширования зависимостей
COPY go.mod go.sum ./
RUN go mod download

# Копируем остальные файлы
COPY . .

# Добавим отладку: покажем, что внутри
RUN ls -la
RUN cat go.mod
RUN cat main.go

# Теперь пытаемся собрать — с выводом ошибки
RUN go build -o order-srv .

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app/order-srv .
EXPOSE 8081
CMD ["./order-srv"]
t