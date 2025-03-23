package main

import (
	"flag"
	"log"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"faraway/server"
)

func main() {
	// Настраиваем логирование на stdout для лучшей видимости в Docker
	log.SetOutput(os.Stdout)
	log.SetFlags(log.Ldate | log.Ltime | log.Lmsgprefix)
	log.SetPrefix("[SERVER] ")

	// Настройка конфигурации сервера
	config := server.DefaultConfig()

	// Параметры командной строки имеют приоритет над конфигурацией по умолчанию
	flag.StringVar(&config.Host, "host", config.Host, "хост для привязки сервера")
	flag.IntVar(&config.Port, "port", config.Port, "порт для прослушивания")
	flag.IntVar(&config.Difficulty, "difficulty", config.Difficulty, "сложность Proof of Work (количество нулевых бит)")
	readTimeout := flag.Int("read-timeout", int(config.ReadTimeout.Seconds()), "таймаут чтения в секундах")
	writeTimeout := flag.Int("write-timeout", int(config.WriteTimeout.Seconds()), "таймаут записи в секундах")
	flag.Parse()

	// Переопределяем конфигурацию из переменных окружения, если они установлены
	if host := os.Getenv("SERVER_HOST"); host != "" {
		config.Host = host
	}
	if portStr := os.Getenv("SERVER_PORT"); portStr != "" {
		if port, err := strconv.Atoi(portStr); err == nil {
			config.Port = port
		}
	}
	if difficultyStr := os.Getenv("POW_DIFFICULTY"); difficultyStr != "" {
		if difficulty, err := strconv.Atoi(difficultyStr); err == nil {
			config.Difficulty = difficulty
		}
	}
	if readTimeoutStr := os.Getenv("READ_TIMEOUT"); readTimeoutStr != "" {
		if timeout, err := strconv.Atoi(readTimeoutStr); err == nil {
			config.ReadTimeout = time.Duration(timeout) * time.Second
		}
	} else {
		config.ReadTimeout = time.Duration(*readTimeout) * time.Second
	}
	if writeTimeoutStr := os.Getenv("WRITE_TIMEOUT"); writeTimeoutStr != "" {
		if timeout, err := strconv.Atoi(writeTimeoutStr); err == nil {
			config.WriteTimeout = time.Duration(timeout) * time.Second
		}
	} else {
		config.WriteTimeout = time.Duration(*writeTimeout) * time.Second
	}

	// Создаем и запускаем сервер
	srv := server.NewServer(config)
	if err := srv.Start(); err != nil {
		log.Printf("Ошибка запуска сервера: %v", err)
	}

	// Обработка сигналов для корректного завершения
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	sig := <-sigChan
	log.Printf("Получен сигнал завершения (%v), начинаем корректное завершение...", sig)

	srv.Stop()
}
