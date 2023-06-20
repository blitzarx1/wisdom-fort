package hash

import (
	"crypto/sha256"
	"encoding/hex"
	"strconv"
)

const tokenZero = '0'

// GenerateHash returns the SHA-256 hash of the combination of token and solution as a hexadecimal string.
func GenerateHash(token string, sol uint64) string {
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

// CheckHash returns whether the given hash meets the difficulty requirement
func CheckHash(hash string, diff uint8) bool {
	// count the number of leading zeroes
	var leadingZeroes int
	for _, rune := range hash {
		if rune != tokenZero {
			break
		}

		leadingZeroes++
	}

	// return whether the number of leading zeroes meets the difficulty requirement
	return leadingZeroes >= int(diff)
}
