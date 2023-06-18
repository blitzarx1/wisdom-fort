package main

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"strings"
)

type PayloadChallenge struct {
	Target uint8 `json:"target"`
}

type PayloadSolution struct {
	Solution uint64 `json:"solution"`
}

type PayloadQuote struct {
	Quote Quote `json:"quote"`
}

type Quote struct {
	Author string `json:"author"`
	Quote  string `json:"quote"`
}

type Request struct {
	Token   *string          `json:"token,omitempty"`
	Action  string           `json:"action"`
	Payload *json.RawMessage `json:"payload,omitempty"`
}

type Response struct {
	Token     string           `json:"token"`
	Payload   *json.RawMessage `json:"payload,omitempty"`
	ErrorCode *string          `json:"error_code,omitempty"`
	Error     *string          `json:"error,omitempty"`
}

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
	challengeRequest := Request{Action: "challenge"}
	challengeRequestBytes, _ := json.Marshal(challengeRequest)
	conn.Write(challengeRequestBytes)

	// Parse challenge response
	challengeData, _ := io.ReadAll(conn)
	conn.Close()

	var challengeResponse Response
	json.Unmarshal(challengeData, &challengeResponse)

	var challengePayload PayloadChallenge
	json.Unmarshal(*challengeResponse.Payload, &challengePayload)

	// Solve challenge
	solution := solveChallenge(challengeResponse.Token, challengePayload.Target)

	// Submit solution and get quote
	solutionPayload := PayloadSolution{Solution: solution}
	solutionPayloadBytes, _ := json.Marshal(solutionPayload)
	solutionRequest := Request{Token: &challengeResponse.Token, Action: "solution", Payload: (*json.RawMessage)(&solutionPayloadBytes)}
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
	var solutionResponse Response
	json.Unmarshal(solutionData, &solutionResponse)

	var quotePayload PayloadQuote
	json.Unmarshal(*solutionResponse.Payload, &quotePayload)

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
