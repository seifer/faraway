package shared

import (
	"fmt"
	"strings"
)

// Solution представляет собой решение Proof of Work
type Solution struct {
	Challenge Challenge // Исходная задача
	Nonce     uint64    // Найденное число-нонс
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
