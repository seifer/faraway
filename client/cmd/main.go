package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"faraway/client"
)

func main() {
	// Настройка конфигурации клиента
	config := client.DefaultConfig()

	// Параметры командной строки имеют приоритет над конфигурацией по умолчанию
	flag.StringVar(&config.ServerHost, "host", config.ServerHost, "хост сервера")
	flag.IntVar(&config.ServerPort, "port", config.ServerPort, "порт сервера")
	connectTimeout := flag.Int("connect-timeout", int(config.ConnectTimeout.Seconds()), "таймаут подключения в секундах")
	responseTimeout := flag.Int("response-timeout", int(config.ResponseTimeout.Seconds()), "таймаут ожидания ответа в секундах")
	flag.Parse()

	// Переопределяем конфигурацию из переменных окружения, если они установлены
	if host := os.Getenv("SERVER_HOST"); host != "" {
		config.ServerHost = host
	}
	if portStr := os.Getenv("SERVER_PORT"); portStr != "" {
		if port, err := strconv.Atoi(portStr); err == nil {
			config.ServerPort = port
		}
	}
	if connectTimeoutStr := os.Getenv("CONNECT_TIMEOUT"); connectTimeoutStr != "" {
		if timeout, err := strconv.Atoi(connectTimeoutStr); err == nil {
			config.ConnectTimeout = time.Duration(timeout) * time.Second
		}
	} else {
		config.ConnectTimeout = time.Duration(*connectTimeout) * time.Second
	}
	if responseTimeoutStr := os.Getenv("RESPONSE_TIMEOUT"); responseTimeoutStr != "" {
		if timeout, err := strconv.Atoi(responseTimeoutStr); err == nil {
			config.ResponseTimeout = time.Duration(timeout) * time.Second
		}
	} else {
		config.ResponseTimeout = time.Duration(*responseTimeout) * time.Second
	}

	// Создаем клиент
	c := client.NewClient(config)

	// Получаем цитату
	fmt.Printf("Подключение к серверу %s:%d...\n", config.ServerHost, config.ServerPort)

	quote, err := c.GetQuote()
	if err != nil {
		log.Fatalf("Ошибка: %v", err)
	}

	fmt.Printf("\nПолучена цитата мудрости:\n\n%s\n", quote)
}
