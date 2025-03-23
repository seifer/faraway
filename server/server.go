package server

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"strings"
	"time"

	"faraway/quotes"
	"faraway/shared"
)

// Config содержит настройки сервера
type Config struct {
	Host         string
	Port         int
	Difficulty   int
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
}

// DefaultConfig возвращает конфигурацию сервера по умолчанию
func DefaultConfig() Config {
	return Config{
		Host:         "0.0.0.0",
		Port:         8080,
		Difficulty:   20, // Это значение можно настраивать в зависимости от требуемого уровня защиты
		ReadTimeout:  time.Second * 10,
		WriteTimeout: time.Second * 10,
	}
}

// Server представляет TCP сервер Word of Wisdom
type Server struct {
	config   Config
	listener net.Listener
}

// NewServer создает новый экземпляр сервера с указанной конфигурацией
func NewServer(config Config) *Server {
	return &Server{
		config: config,
	}
}

// Start запускает сервер и начинает прослушивание соединений
func (s *Server) Start() error {
	addr := fmt.Sprintf("%s:%d", s.config.Host, s.config.Port)
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		return fmt.Errorf("ошибка запуска сервера: %v", err)
	}
	s.listener = listener

	log.Printf("Сервер 'Word of Wisdom' запущен на %s", addr)
	log.Printf("Сложность Proof of Work: %d", s.config.Difficulty)

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("Ошибка принятия соединения: %v", err)
			continue
		}

		go s.handleConnection(conn)
	}
}

// Stop останавливает сервер
func (s *Server) Stop() error {
	if s.listener != nil {
		return s.listener.Close()
	}
	return nil
}

// handleConnection обрабатывает входящее TCP соединение
func (s *Server) handleConnection(conn net.Conn) {
	defer conn.Close()
	addr := conn.RemoteAddr().String()
	log.Printf("Новое соединение от %s", addr)

	// Устанавливаем таймауты
	conn.SetReadDeadline(time.Now().Add(s.config.ReadTimeout))
	conn.SetWriteDeadline(time.Now().Add(s.config.WriteTimeout))

	// Шаг 1: Отправляем клиенту задачу Proof of Work
	challenge := shared.GenerateChallenge(s.config.Difficulty)
	challengeStr := challenge.Encode()
	_, err := fmt.Fprintf(conn, "CHALLENGE %s\n", challengeStr)
	if err != nil {
		log.Printf("Ошибка отправки задачи клиенту %s: %v", addr, err)
		return
	}
	log.Printf("Отправлена задача клиенту %s: %s", addr, challengeStr)

	// Шаг 2: Ожидаем решение от клиента
	reader := bufio.NewReader(conn)
	line, err := reader.ReadString('\n')
	if err != nil {
		log.Printf("Ошибка чтения ответа от клиента %s: %v", addr, err)
		return
	}

	line = strings.TrimSpace(line)
	if !strings.HasPrefix(line, "SOLUTION ") {
		log.Printf("Получен неверный формат решения от клиента %s: %s", addr, line)
		conn.Write([]byte("ERROR: неверный формат решения\n"))
		return
	}

	solutionStr := strings.TrimPrefix(line, "SOLUTION ")
	solution, err := shared.DecodeSolution(solutionStr)
	if err != nil {
		log.Printf("Ошибка декодирования решения от клиента %s: %v", addr, err)
		conn.Write([]byte(fmt.Sprintf("ERROR: %v\n", err)))
		return
	}

	log.Printf("Получено решение от клиента %s: nonce=%d", addr, solution.Nonce)

	// Шаг 3: Проверяем решение
	if !solution.Verify() || solution.Challenge.Prefix != challenge.Prefix || solution.Challenge.Difficulty != challenge.Difficulty {
		log.Printf("Клиент %s предоставил неверное решение", addr)
		conn.Write([]byte("ERROR: неверное решение\n"))
		return
	}

	log.Printf("Клиент %s успешно решил задачу Proof of Work", addr)

	// Шаг 4: Отправляем цитату
	quote := quotes.GetRandomQuote()
	_, err = fmt.Fprintf(conn, "QUOTE %s\n", quote)
	if err != nil {
		log.Printf("Ошибка отправки цитаты клиенту %s: %v", addr, err)
		return
	}

	log.Printf("Отправлена цитата клиенту %s", addr)
}
