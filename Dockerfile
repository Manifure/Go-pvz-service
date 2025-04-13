FROM golang:1.23 AS builder
LABEL authors="manifure"

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

ARG SERVICE_PATH
RUN go build -o main $SERVICE_PATH

FROM frolvlad/alpine-glibc

WORKDIR /root/

COPY --from=builder /app/main .

EXPOSE 8080:8080

# Запускаем приложение, sleep на 30 секунд, что бы необходимые контейнеры успели загрузиться
CMD ["sh", "-c", "sleep 30 && ./main"]