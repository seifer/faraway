package shared

import (
	"crypto/sha256"
	"encoding/binary"
)

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
