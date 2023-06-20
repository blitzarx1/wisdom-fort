package client

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"strings"

	"blitzarx1/wisdom-fort/pkg/logger"
	"blitzarx1/wisdom-fort/pkg/scheme"
)

const protocol = "tcp"

type Client struct {
	cfg *Config
}

func New(ctx context.Context, cfg *Config) *Client {
	l, err := logger.FromCtx(ctx)
	if err != nil {
		l = logger.New(nil, "new")
	}

	l.Println("initializing client")

	return &Client{cfg: cfg}
}

// GetQuote gets a quote from the server. Underhood it will call GetChallenge, SolveChallenge and SubmitSolution.
func (c *Client) GetQuote(ctx context.Context) (*scheme.Quote, error) {
	l := logger.MustFromCtx(ctx)
	l.Println("getting quote from the server")

	challenge, err := c.GetChallenge(logger.WithCtx(ctx, l, "getChallenge"))
	if err != nil {
		return nil, err
	}

	l.Println("got challenge: ", challenge)

	solution := c.SolveChallenge(logger.WithCtx(ctx, l, "solveChallenge"), challenge)

	l.Println("submitting solution: ", solution)
	return c.SubmitSolution(logger.WithCtx(ctx, l, "submitSolution"), solution, challenge.Token)
}

// GetChallenge gets a challenge from the server.
func (c *Client) GetChallenge(ctx context.Context) (*Challenge, error) {
	l := logger.MustFromCtx(ctx)
	l.Println("getting challenge from the server")

	conn, err := c.connect(logger.WithCtx(ctx, l, "connect"), c.cfg.ServerHost, c.cfg.ServerPort)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	challengeRequest := scheme.Request{Action: scheme.ActionChallenge}
	challengeRequestBytes, err := json.Marshal(challengeRequest)
	if err != nil {
		return nil, err
	}
	conn.Write(challengeRequestBytes)

	challengeData, err := io.ReadAll(conn)
	if err != nil {
		return nil, err
	}

	var challengeResponse scheme.Response
	if err := json.Unmarshal(challengeData, &challengeResponse); err != nil {
		return nil, err
	}

	var challengePayload scheme.PayloadResponseChallenge
	if err := json.Unmarshal(challengeResponse.Payload, &challengePayload); err != nil {
		return nil, err
	}

	return &Challenge{
		Token:      challengeResponse.Token,
		Difficulty: challengePayload.Target,
	}, nil
}

// SubmitSolution submits a solution to the server and gets a quote.
func (c *Client) SubmitSolution(ctx context.Context, solution uint, token string) (*scheme.Quote, error) {
	l := logger.MustFromCtx(ctx)
	l.Println("sbmitting solution to the server")

	conn, err := c.connect(logger.WithCtx(ctx, l, "connect"), c.cfg.ServerHost, c.cfg.ServerPort)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	solutionPayload := scheme.PayloadRequestSolution{Solution: solution}
	solutionPayloadBytes, err := json.Marshal(solutionPayload)
	if err != nil {
		return nil, err
	}

	solutionRequest := scheme.Request{
		Token:   &token,
		Action:  scheme.ActionSolution,
		Payload: solutionPayloadBytes,
	}
	solutionRequestBytes, err := json.Marshal(solutionRequest)
	if err != nil {
		return nil, err
	}
	conn.Write(solutionRequestBytes)

	solutionData, err := io.ReadAll(conn)
	if err != nil {
		return nil, err
	}

	var solutionResponse scheme.Response
	if err := json.Unmarshal(solutionData, &solutionResponse); err != nil {
		return nil, err
	}

	var quotePayload scheme.PayloadResponseSolution
	if err := json.Unmarshal(solutionResponse.Payload, &quotePayload); err != nil {
		return nil, err
	}

	return &quotePayload.Quote, nil
}

// SolveChallenge solves a challenge.
func (c *Client) SolveChallenge(ctx context.Context, challenge *Challenge) uint {
	l := logger.MustFromCtx(ctx)
	l.Println("solving challenge: ", challenge)
	return solveChallenge(challenge.Token, challenge.Difficulty)
}

func (c *Client) connect(ctx context.Context, host string, port uint) (net.Conn, error) {
	l := logger.MustFromCtx(ctx)

	serverAddr := fmt.Sprintf("%s:%d", host, port)
	l.Println("establishing connection with the server: ", serverAddr)

	return net.Dial(protocol, serverAddr)
}

func solveChallenge(token string, difficulty uint8) uint {
	var nonce uint

	for {
		// generate hash of the solution combined with the token
		solution := fmt.Sprintf("%s%d", token, nonce)
		hash := sha256.Sum256([]byte(solution))
		hashHex := hex.EncodeToString(hash[:])

		// count leading zeroes
		zeroes := strings.IndexFunc(hashHex, func(r rune) bool { return r != '0' })

		if zeroes >= int(difficulty) {
			break
		}

		nonce++
	}

	return nonce
}
