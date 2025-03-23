package client

import (
	"bufio"
	"context"
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
func (c *Client) GetQuote(ctx context.Context) (string, error) {
	// Создаем контекст с таймаутом для подключения
	dialCtx, cancel := context.WithTimeout(ctx, c.config.ConnectTimeout)
	defer cancel()

	// Подключаемся к серверу
	addr := fmt.Sprintf("%s:%d", c.config.ServerHost, c.config.ServerPort)
	var dialer net.Dialer
	conn, err := dialer.DialContext(dialCtx, "tcp", addr)
	if err != nil {
		return "", fmt.Errorf("ошибка подключения к серверу: %v", err)
	}
	defer conn.Close()

	// Шаг 1: Получаем задачу от сервера
	conn.SetDeadline(time.Now().Add(c.config.ResponseTimeout))
	reader := bufio.NewReader(conn)

	challengeLine, err := reader.ReadString('\n')
	if err != nil {
		return "", fmt.Errorf("ошибка чтения задачи от сервера: %v", err)
	}

	// Проверка контекста
	if ctx.Err() != nil {
		return "", fmt.Errorf("операция отменена: %v", ctx.Err())
	}

	challengeLine = strings.TrimSpace(challengeLine)
	if !strings.HasPrefix(challengeLine, "CHALLENGE ") {
		return "", fmt.Errorf("неверный формат задачи: %s", challengeLine)
	}

	challengeStr := strings.TrimPrefix(challengeLine, "CHALLENGE ")
	challenge, err := shared.DecodeChallenge(challengeStr)
	if err != nil {
		return "", fmt.Errorf("ошибка декодирования задачи: %v", err)
	}

	fmt.Printf("Получена задача от сервера: %s (сложность %d)\n", challenge.Prefix, challenge.Difficulty)
	fmt.Printf("Решаем задачу Proof of Work (это может занять некоторое время)...\n")

	// Шаг 2: Решаем задачу, периодически проверяя контекст
	startTime := time.Now()

	// Проверяем, не отменён ли контекст перед вычислениями
	if ctx.Err() != nil {
		return "", fmt.Errorf("операция отменена: %v", ctx.Err())
	}

	// Вычисляем решение напрямую - по-хорошему здесь нужна реализация с поддержкой отмены
	// Но пока используем стандартный метод Solve
	solution := challenge.Solve()

	// Проверяем контекст после решения
	if ctx.Err() != nil {
		return "", fmt.Errorf("операция отменена: %v", ctx.Err())
	}

	elapsedTime := time.Since(startTime)
	fmt.Printf("Задача решена за %v, найден nonce: %d\n", elapsedTime, solution.Nonce)

	// Шаг 3: Отправляем решение серверу
	solutionStr := solution.Encode()
	conn.SetDeadline(time.Now().Add(c.config.ResponseTimeout))
	_, err = fmt.Fprintf(conn, "SOLUTION %s\n", solutionStr)
	if err != nil {
		return "", fmt.Errorf("ошибка отправки решения: %v", err)
	}

	// Шаг 4: Получаем ответ сервера
	conn.SetDeadline(time.Now().Add(c.config.ResponseTimeout))
	resultLine, err := reader.ReadString('\n')
	if err != nil {
		return "", fmt.Errorf("ошибка чтения ответа: %v", err)
	}

	// Проверка контекста
	if ctx.Err() != nil {
		return "", fmt.Errorf("операция отменена: %v", ctx.Err())
	}

	resultLine = strings.TrimSpace(resultLine)
	if strings.HasPrefix(resultLine, "ERROR: ") {
		return "", fmt.Errorf("ошибка сервера: %s", strings.TrimPrefix(resultLine, "ERROR: "))
	}

	if !strings.HasPrefix(resultLine, "QUOTE ") {
		return "", fmt.Errorf("неверный формат ответа: %s", resultLine)
	}

	quote := strings.TrimPrefix(resultLine, "QUOTE ")
	return quote, nil
}
