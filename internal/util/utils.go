package util

import (
	"fmt"
	"math/rand"
	"os"
	"time"
)

const requestIDLength = 4

// GenerateRequestID generates a request ID with the format "RQ" + YYYYMMDDHHMMSS + random 4 digits.
func GenerateRequestID() string {
	now := time.Now()
	timestamp := now.Format("20060102150405") // Format: YYYYMMDDHHMMSS

	b := make([]rune, requestIDLength)
	var digitRunes = []rune("0123456789")

	for i := range b {
		b[i] = digitRunes[rand.Intn(len(digitRunes))]
	}
	randomDigits := string(b)

	return fmt.Sprintf("RQ%s%s", timestamp, randomDigits)
}

// WriteResultToFile writes a slice of strings to a file.
func WriteResultToFile(filePath string, results []string) error {
	file, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("failed to create result file '%s': %v", filePath, err)
	}
	defer file.Close()

	for _, result := range results {
		_, err = file.WriteString(result + "\n")
		if err != nil {
			return fmt.Errorf("failed to write to result file '%s': %v", filePath, err)
		}
	}

	return nil
}

// safeSubstring is a helper function to extract a substring safely
func SafeSubstring(s string, start, length int) string {
	if start < 0 || start >= len(s) {
		return ""
	}
	end := start + length
	if end > len(s) {
		end = len(s)
	}
	return s[start:end]
}
