# go_tutorial_minishop/product-service/Dockerfile
FROM golang:1.21-alpine AS builder
WORKDIR /app
COPY . .
RUN go mod download
RUN go build -o product-srv .

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app/product-srv .
EXPOSE 8080
CMD ["./product-srv"]
