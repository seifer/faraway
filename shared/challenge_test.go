package shared

import (
	"math"
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

func TestEstimateTime(t *testing.T) {
	// Проверяем, что функция EstimateTime возвращает корректные результаты
	testCases := []struct {
		difficulty      int
		hashesPerSecond float64
	}{
		{8, 1000000},  // ~256 попыток при 1M хешей/сек
		{16, 1000000}, // ~65536 попыток при 1M хешей/сек
	}

	for i, tc := range testCases {
		result := EstimateTime(tc.difficulty, tc.hashesPerSecond)
		// Проверяем формулу: 2^difficulty / hashesPerSecond
		expectedDuration := time.Duration(math.Pow(2, float64(tc.difficulty)) / tc.hashesPerSecond * float64(time.Second))

		if result != expectedDuration {
			t.Errorf("Тест %d: ожидалось %v, получено %v для сложности %d и %f хешей/сек",
				i, expectedDuration, result, tc.difficulty, tc.hashesPerSecond)
		}
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
