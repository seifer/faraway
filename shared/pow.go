package shared

import (
	"crypto/sha256"
	"encoding/base64"
	"encoding/binary"
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

// Solution представляет собой решение Proof of Work
type Solution struct {
	Challenge Challenge // Исходная задача
	Nonce     uint64    // Найденное число-нонс
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

// Encode кодирует решение в строку для передачи по сети
func (s Solution) Encode() string {
	return fmt.Sprintf("%s:%d", s.Challenge.Encode(), s.Nonce)
}

// DecodeSolution декодирует решение из строки
func DecodeSolution(encoded string) (Solution, error) {
	lastColonIndex := strings.LastIndex(encoded, ":")
	if lastColonIndex == -1 {
		return Solution{}, fmt.Errorf("неверный формат решения")
	}

	challengeStr := encoded[:lastColonIndex]
	nonceStr := encoded[lastColonIndex+1:]

	challenge, err := DecodeChallenge(challengeStr)
	if err != nil {
		return Solution{}, err
	}

	var nonce uint64
	_, err = fmt.Sscanf(nonceStr, "%d", &nonce)
	if err != nil {
		return Solution{}, fmt.Errorf("неверный формат nonce: %v", err)
	}

	return Solution{
		Challenge: challenge,
		Nonce:     nonce,
	}, nil
}

// Verify проверяет, является ли решение действительным
func (s Solution) Verify() bool {
	return s.Challenge.Verify(s.Nonce)
}

// computeHash вычисляет SHA-256 хеш от префикса и nonce
func computeHash(prefix string, nonce uint64) []byte {
	nonceBytes := make([]byte, 8)
	binary.LittleEndian.PutUint64(nonceBytes, nonce)

	data := append([]byte(prefix), nonceBytes...)
	hash := sha256.Sum256(data)
	return hash[:]
}

// countLeadingZeros подсчитывает количество ведущих нулевых бит в хеше
func countLeadingZeros(hash []byte) int {
	count := 0
	for _, b := range hash {
		if b == 0 {
			count += 8
			continue
		}

		// Подсчет ведущих нулевых бит в байте
		zeros := 0
		mask := byte(128) // 10000000
		for i := 0; i < 8; i++ {
			if b&mask == 0 {
				zeros++
				mask >>= 1
			} else {
				break
			}
		}

		count += zeros
		break
	}

	return count
}

// EstimateTime оценивает среднее время для решения задачи с заданной сложностью
func EstimateTime(difficulty int, hashesPerSecond float64) time.Duration {
	averageAttempts := math.Pow(2, float64(difficulty))
	seconds := averageAttempts / hashesPerSecond
	return time.Duration(seconds * float64(time.Second))
}
