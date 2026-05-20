# product-service/Dockerfile
FROM golang:1.21-alpine AS builder
WORKDIR /app

# Сначала копируем go.mod и go.sum — для кэширования зависимостей
COPY go.mod go.sum ./
RUN go mod download

# Копируем остальные файлы
COPY . .

# Отладка
RUN ls -la
RUN cat go.mod

# Собираем
RUN go build -o product-srv .

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app/product-srv .
EXPOSE 8080
CMD ["./product-srv"]
