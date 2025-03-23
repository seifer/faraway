package client

import (
	"bufio"
	"fmt"
	"net"
	"strings"
	"time"

	"faraway/shared"
)

// Config содержит настройки клиента
type Config struct {
	ServerHost      string
	ServerPort      int
	ConnectTimeout  time.Duration
	ResponseTimeout time.Duration
}

// DefaultConfig возвращает конфигурацию клиента по умолчанию
func DefaultConfig() Config {
	return Config{
		ServerHost:      "localhost",
		ServerPort:      8080,
		ConnectTimeout:  time.Second * 5,
		ResponseTimeout: time.Second * 30,
	}
}

// Client представляет клиент для "Word of Wisdom" сервера
type Client struct {
	config Config
}

// NewClient создает новый экземпляр клиента с указанной конфигурацией
func NewClient(config Config) *Client {
	return &Client{
		config: config,
	}
}

// GetQuote подключается к серверу, решает задачу PoW и получает цитату
func (c *Client) GetQuote() (string, error) {
	// Подключение к серверу
	addr := fmt.Sprintf("%s:%d", c.config.ServerHost, c.config.ServerPort)
	conn, err := net.DialTimeout("tcp", addr, c.config.ConnectTimeout)
	if err != nil {
		return "", fmt.Errorf("ошибка подключения к серверу: %v", err)
	}
	defer conn.Close()

	// Устанавливаем таймаут на чтение ответа
	conn.SetReadDeadline(time.Now().Add(c.config.ResponseTimeout))

	// Читаем задачу от сервера
	reader := bufio.NewReader(conn)
	challengeLine, err := reader.ReadString('\n')
	if err != nil {
		return "", fmt.Errorf("ошибка чтения задачи от сервера: %v", err)
	}

	challengeLine = strings.TrimSpace(challengeLine)
	if !strings.HasPrefix(challengeLine, "CHALLENGE ") {
		return "", fmt.Errorf("получен неверный формат задачи: %s", challengeLine)
	}

	challengeStr := strings.TrimPrefix(challengeLine, "CHALLENGE ")
	challenge, err := shared.DecodeChallenge(challengeStr)
	if err != nil {
		return "", fmt.Errorf("ошибка декодирования задачи: %v", err)
	}

	fmt.Printf("Получена задача от сервера: %s (сложность %d)\n", challenge.Prefix, challenge.Difficulty)
	fmt.Printf("Решаем задачу Proof of Work (это может занять некоторое время)...\n")

	// Засекаем время начала решения
	startTime := time.Now()

	// Решаем задачу
	solution := challenge.Solve()

	// Вычисляем затраченное время
	elapsedTime := time.Since(startTime)
	fmt.Printf("Задача решена за %v, найден nonce: %d\n", elapsedTime, solution.Nonce)

	// Отправляем решение серверу
	solutionStr := solution.Encode()
	_, err = fmt.Fprintf(conn, "SOLUTION %s\n", solutionStr)
	if err != nil {
		return "", fmt.Errorf("ошибка отправки решения серверу: %v", err)
	}

	// Читаем ответ сервера (цитату или ошибку)
	resultLine, err := reader.ReadString('\n')
	if err != nil {
		return "", fmt.Errorf("ошибка чтения ответа от сервера: %v", err)
	}

	resultLine = strings.TrimSpace(resultLine)
	if strings.HasPrefix(resultLine, "ERROR: ") {
		return "", fmt.Errorf("сервер вернул ошибку: %s", strings.TrimPrefix(resultLine, "ERROR: "))
	}

	if !strings.HasPrefix(resultLine, "QUOTE ") {
		return "", fmt.Errorf("получен неверный формат ответа: %s", resultLine)
	}

	quote := strings.TrimPrefix(resultLine, "QUOTE ")
	return quote, nil
}
