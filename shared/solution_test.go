package shared

import (
	"testing"
)

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

func TestSolutionVerify(t *testing.T) {
	// Создаем задачу с низкой сложностью
	challenge := GenerateChallenge(4)

	// Решаем задачу
	solution := challenge.Solve()

	// Проверяем, что решение проходит проверку
	if !solution.Verify() {
		t.Errorf("Правильное решение не прошло верификацию")
	}

	// Создаем неправильное решение
	invalidSolution := Solution{
		Challenge: challenge,
		Nonce:     solution.Nonce + 1,
	}

	// Проверяем, что неправильное решение не проходит проверку
	if invalidSolution.Verify() {
		t.Errorf("Неправильное решение прошло верификацию")
	}
}
