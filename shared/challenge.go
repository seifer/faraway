package shared

import (
	"encoding/base64"
	"fmt"
	"math"
	"math/rand"
	"strings"
	"time"
)

// Challenge представляет собой задачу Proof of Work
type Challenge struct {
	Prefix     string // Префикс для генерации
	Difficulty int    // Сложность задачи (количество нулевых бит в начале)
}

// GenerateChallenge создает новую задачу с указанной сложностью
func GenerateChallenge(difficulty int) Challenge {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	prefix := make([]byte, 16)
	r.Read(prefix)
	return Challenge{
		Prefix:     base64.StdEncoding.EncodeToString(prefix),
		Difficulty: difficulty,
	}
}

// Encode кодирует задачу в строку для передачи по сети
func (c Challenge) Encode() string {
	return fmt.Sprintf("%s:%d", c.Prefix, c.Difficulty)
}

// DecodeChallenge декодирует задачу из строки
func DecodeChallenge(encoded string) (Challenge, error) {
	parts := strings.Split(encoded, ":")
	if len(parts) != 2 {
		return Challenge{}, fmt.Errorf("неверный формат задачи")
	}

	difficulty := 0
	_, err := fmt.Sscanf(parts[1], "%d", &difficulty)
	if err != nil {
		return Challenge{}, fmt.Errorf("неверный формат сложности: %v", err)
	}

	return Challenge{
		Prefix:     parts[0],
		Difficulty: difficulty,
	}, nil
}

// Solve решает задачу Proof of Work
func (c Challenge) Solve() Solution {
	var nonce uint64 = 0
	for {
		if c.Verify(nonce) {
			break
		}
		nonce++
	}
	return Solution{
		Challenge: c,
		Nonce:     nonce,
	}
}

// Verify проверяет, является ли данный nonce решением задачи
func (c Challenge) Verify(nonce uint64) bool {
	hash := computeHash(c.Prefix, nonce)

	// Проверяем, что первые c.Difficulty бит равны нулю
	return countLeadingZeros(hash) >= c.Difficulty
}

// EstimateTime оценивает среднее время для решения задачи с заданной сложностью
func EstimateTime(difficulty int, hashesPerSecond float64) time.Duration {
	averageAttempts := math.Pow(2, float64(difficulty))
	seconds := averageAttempts / hashesPerSecond
	return time.Duration(seconds * float64(time.Second))
}
