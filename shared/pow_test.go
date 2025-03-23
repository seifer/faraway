package shared

import (
	"testing"
	"time"
)

func TestChallengeVerify(t *testing.T) {
	challenge := GenerateChallenge(4) // Используем низкую сложность для тестов

	// Проверяем, что сгенерированная задача имеет корректные параметры
	if challenge.Prefix == "" {
		t.Error("Префикс не должен быть пустым")
	}
	if challenge.Difficulty != 4 {
		t.Errorf("Сложность должна быть 4, получено %d", challenge.Difficulty)
	}

	// Находим решение задачи
	solution := challenge.Solve()

	// Проверяем, что решение действительно решает задачу
	if !challenge.Verify(solution.Nonce) {
		t.Errorf("Решение %d не прошло верификацию для префикса %s и сложности %d",
			solution.Nonce, challenge.Prefix, challenge.Difficulty)
	}

	// Проверяем, что неправильное решение не проходит верификацию
	if challenge.Verify(solution.Nonce + 1) {
		t.Errorf("Неправильное решение %d прошло верификацию для префикса %s и сложности %d",
			solution.Nonce+1, challenge.Prefix, challenge.Difficulty)
	}
}

func TestChallengeEncodeDecode(t *testing.T) {
	original := Challenge{
		Prefix:     "test-prefix",
		Difficulty: 10,
	}

	// Кодируем задачу в строку
	encoded := original.Encode()

	// Декодируем задачу из строки
	decoded, err := DecodeChallenge(encoded)
	if err != nil {
		t.Errorf("Ошибка декодирования задачи: %v", err)
	}

	// Проверяем, что декодированная задача совпадает с оригинальной
	if decoded.Prefix != original.Prefix {
		t.Errorf("Префикс декодированной задачи (%s) не совпадает с оригинальным (%s)",
			decoded.Prefix, original.Prefix)
	}
	if decoded.Difficulty != original.Difficulty {
		t.Errorf("Сложность декодированной задачи (%d) не совпадает с оригинальной (%d)",
			decoded.Difficulty, original.Difficulty)
	}

	// Проверяем обработку некорректных строк
	_, err = DecodeChallenge("invalid-format")
	if err == nil {
		t.Error("Декодирование некорректной строки должно возвращать ошибку")
	}
}

func TestSolutionEncodeDecode(t *testing.T) {
	originalChallenge := Challenge{
		Prefix:     "test-prefix",
		Difficulty: 10,
	}
	original := Solution{
		Challenge: originalChallenge,
		Nonce:     12345,
	}

	// Кодируем решение в строку
	encoded := original.Encode()

	// Декодируем решение из строки
	decoded, err := DecodeSolution(encoded)
	if err != nil {
		t.Errorf("Ошибка декодирования решения: %v", err)
	}

	// Проверяем, что декодированное решение совпадает с оригинальным
	if decoded.Challenge.Prefix != original.Challenge.Prefix {
		t.Errorf("Префикс декодированного решения (%s) не совпадает с оригинальным (%s)",
			decoded.Challenge.Prefix, original.Challenge.Prefix)
	}
	if decoded.Challenge.Difficulty != original.Challenge.Difficulty {
		t.Errorf("Сложность декодированного решения (%d) не совпадает с оригинальной (%d)",
			decoded.Challenge.Difficulty, original.Challenge.Difficulty)
	}
	if decoded.Nonce != original.Nonce {
		t.Errorf("Nonce декодированного решения (%d) не совпадает с оригинальным (%d)",
			decoded.Nonce, original.Nonce)
	}

	// Проверяем обработку некорректных строк
	_, err = DecodeSolution("invalid-format")
	if err == nil {
		t.Error("Декодирование некорректной строки должно возвращать ошибку")
	}
}

func TestProofOfWorkPerformance(t *testing.T) {
	// Тестируем производительность PoW с разной сложностью
	difficulties := []int{8, 12, 16}

	for _, difficulty := range difficulties {
		challenge := GenerateChallenge(difficulty)

		start := time.Now()
		solution := challenge.Solve()
		elapsed := time.Since(start)

		t.Logf("Сложность %d: решение найдено за %v, nonce = %d",
			difficulty, elapsed, solution.Nonce)

		// Проверяем, что решение корректно
		if !solution.Verify() {
			t.Errorf("Решение не прошло верификацию для сложности %d", difficulty)
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
