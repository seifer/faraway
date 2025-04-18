FROM golang:1.20-alpine AS builder

WORKDIR /app

# Копируем только необходимые файлы для сборки зависимостей
COPY go.mod ./
RUN go mod download

# Копируем исходный код
COPY . .

# Собираем приложение
RUN CGO_ENABLED=0 GOOS=linux go build -a -o server-app ./server/cmd

# Конечный образ
FROM alpine:latest

WORKDIR /app

# Добавляем tzdata для корректной работы с датами
RUN apk --no-cache add tzdata

# Копируем исполняемый файл из образа сборки
COPY --from=builder /app/server-app .

# Настройки среды выполнения
ENV SERVER_HOST=0.0.0.0
ENV SERVER_PORT=8080
ENV POW_DIFFICULTY=10
# READ_TIMEOUT используется как дополнительное время к расчетному таймауту PoW
ENV READ_TIMEOUT=10
ENV WRITE_TIMEOUT=10

# Открываем порт
EXPOSE 8080

# Запускаем сервер
CMD ["./server-app"] 