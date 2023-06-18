package service

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"math/rand"
	"strings"
	"time"
)

const (
	separatorToken = "-"
)

type Token string

// newToken returns a new Token generated from the given IP address
// and the current time plus random part
func newToken(ip string) Token {
	token := fmt.Sprintf("%s"+separatorToken+"%d"+separatorToken+"%s", ip, time.Now().UnixNano(), generateRandomPart())
	return Token(token)
}

func (t Token) ip() string {
	return strings.Split(string(t), separatorToken)[0]
}

func generateRandomPart() string {
	randomNumber := rand.Int63()

	hasher := md5.New()
	hasher.Write([]byte(fmt.Sprintf("%d", randomNumber)))

	return hex.EncodeToString(hasher.Sum(nil))
}
