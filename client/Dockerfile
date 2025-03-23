FROM golang:1.20-alpine AS builder

WORKDIR /app

# Копируем только необходимые файлы для сборки зависимостей
COPY go.mod ./
RUN go mod download

# Копируем исходный код
COPY . .

# Собираем приложение
RUN CGO_ENABLED=0 GOOS=linux go build -a -o client-app ./client/cmd

# Конечный образ
FROM alpine:latest

WORKDIR /app

# Копируем исполняемый файл из образа сборки
COPY --from=builder /app/client-app .

# Настройки среды выполнения по умолчанию
ENV SERVER_HOST=server
ENV SERVER_PORT=8080
ENV CONNECT_TIMEOUT=5
ENV RESPONSE_TIMEOUT=30

# Запускаем клиент
CMD ["./client-app"] 