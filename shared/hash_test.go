package shared

import (
	"crypto/sha256"
	"testing"
)

func TestComputeHash(t *testing.T) {
	// Проверяем, что функция computeHash генерирует корректные хеши
	testCases := []struct {
		prefix string
		nonce  uint64
	}{
		{"test", 0},
		{"", 42},
		{"prefix123", 9999},
	}

	for i, tc := range testCases {
		hash := computeHash(tc.prefix, tc.nonce)

		// Создаем хеш вручную для проверки
		nonceBytes := make([]byte, 8)
		for j := 0; j < 8; j++ {
			nonceBytes[j] = byte(tc.nonce >> (j * 8))
		}
		data := append([]byte(tc.prefix), nonceBytes...)
		expectedHash := sha256.Sum256(data)

		// Проверяем совпадение хешей
		for j := range hash {
			if hash[j] != expectedHash[j] {
				t.Errorf("Тест %d: хеш не совпадает с ожидаемым для префикса '%s' и nonce %d",
					i, tc.prefix, tc.nonce)
				break
			}
		}
	}
}

func TestCountLeadingZeros(t *testing.T) {
	testCases := []struct {
		input    []byte
		expected int
	}{
		{[]byte{0, 0, 128, 0}, 16}, // 2 байта нулевых + 1 бит
		{[]byte{0, 0, 0, 0}, 32},   // 4 байта нулевых
		{[]byte{128, 0, 0, 0}, 0},  // 0 бит (первый бит равен 1)
		{[]byte{1, 0, 0, 0}, 7},    // 7 бит
		{[]byte{0, 1, 0, 0}, 15},   // 1 байт + 7 бит
	}

	for i, tc := range testCases {
		result := countLeadingZeros(tc.input)
		if result != tc.expected {
			t.Errorf("Тест %d: ожидалось %d нулевых бит, получено %d для %v",
				i, tc.expected, result, tc.input)
		}
	}
}
