package client

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"strings"

	"blitzarx1/wisdom-fort/pkg/scheme"
)

const (
	ServerAddr = "localhost:8080"
)

func main() {
	// Connect to the server
	conn, err := net.Dial("tcp", ServerAddr)
	if err != nil {
		fmt.Println(err)
		return
	}

	// Start challenge
	challengeRequest := scheme.Request{Action: scheme.ActionChallenge}
	challengeRequestBytes, _ := json.Marshal(challengeRequest)
	conn.Write(challengeRequestBytes)

	// Parse challenge response
	challengeData, _ := io.ReadAll(conn)
	conn.Close()

	var challengeResponse scheme.Response
	json.Unmarshal(challengeData, &challengeResponse)

	var challengePayload scheme.PayloadResponseChallenge
	json.Unmarshal(challengeResponse.Payload, &challengePayload)

	// Solve challenge
	solution := solveChallenge(challengeResponse.Token, challengePayload.Target)

	// Submit solution and get quote
	solutionPayload := scheme.PayloadRequestSolution{Solution: solution}
	solutionPayloadBytes, _ := json.Marshal(solutionPayload)
	solutionRequest := scheme.Request{
		Token:   &challengeResponse.Token,
		Action:  scheme.ActionSolution,
		Payload: solutionPayloadBytes,
	}
	solutionRequestBytes, _ := json.Marshal(solutionRequest)

	// Connect to the server
	conn, err = net.Dial("tcp", ServerAddr)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer conn.Close()
	conn.Write(solutionRequestBytes)

	// Parse solution response
	solutionData, _ := io.ReadAll(conn)
	var solutionResponse scheme.Response
	json.Unmarshal(solutionData, &solutionResponse)

	var quotePayload scheme.PayloadResponseSolution
	json.Unmarshal(solutionResponse.Payload, &quotePayload)

	// Print quote
	fmt.Println(quotePayload.Quote)
	conn.Close()
}

func solveChallenge(token string, difficulty uint8) uint64 {
	nonce := uint64(0)

	for {
		// Generate hash of the solution combined with the token
		solution := fmt.Sprintf("%s%d", token, nonce)
		hash := sha256.Sum256([]byte(solution))
		hashHex := hex.EncodeToString(hash[:])

		// Count leading zeroes
		zeroes := strings.IndexFunc(hashHex, func(r rune) bool { return r != '0' })

		if zeroes >= int(difficulty) {
			break
		}

		nonce++
	}

	return nonce
}
