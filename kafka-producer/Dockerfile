# kafka-producer/Dokerfile
FROM golang:1.23-alpine AS builder
WORKDIR /app
COPY . .
RUN mod download
RUN build -o producer .

FROM alpine:latest
RUN apk --no-cache add ca-certificate
WORKDIR /root/
COPY --from=builder /app/producer .
EXPOSE 8085
CMD ["./producer"]
