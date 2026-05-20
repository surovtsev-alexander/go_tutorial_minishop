# go_tutorial_minishop/order-service/Dockerfile
FROM golang:1.21-alpine AS builder
WORKDIR /app
COPY . .
RUN go mod download
RUN go build -o order-srv .

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app/order-srv .
EXPOSE 8081
CMD ["./order-srv"]
