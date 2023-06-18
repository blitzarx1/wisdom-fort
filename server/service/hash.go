package service

import (
	"crypto/sha256"
	"encoding/hex"
	"strconv"
)

// generateHash returns the SHA-256 hash of the combination of token and solution as a hexadecimal string.
func generateHash(token Token, sol solution) string {
	// convert solution to string
	solStr := strconv.FormatUint(uint64(sol), 10)

	// combine token and solution into a single string
	data := string(token) + solStr

	// generate SHA-256 hash
	hasher := sha256.New()
	hasher.Write([]byte(data))
	hash := hasher.Sum(nil)

	// return hash as a hexadecimal string
	return hex.EncodeToString(hash)
}

// checkHash returns whether the given hash meets the difficulty requirement
func checkHash(hash string, diff difficulty) bool {
	// count the number of leading zeroes
	leadingZeroes := 0
	for _, rune := range hash {
		if rune != '0' {
			break
		}

		leadingZeroes++
	}

	// return whether the number of leading zeroes meets the difficulty requirement
	return leadingZeroes >= int(diff)
}
