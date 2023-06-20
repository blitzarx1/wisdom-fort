package main

import (
	"context"

	"blitzarx1/wisdom-fort/client"
	"blitzarx1/wisdom-fort/pkg/logger"
)

const (
	port = 8080
	host = "localhost"
)

// demoChallenge is a demonstration of the basics of getting solution and solving it.
func demoChallenge(ctx context.Context, c *client.Client) (uint, string, error) {
	l := logger.MustFromCtx(ctx)
	l.Println("demo challenge")

	challenge, err := c.GetChallenge(logger.WithCtx(ctx, l, "getChallenge"))
	if err != nil {
		return 0, "", err
	}

	l.Println("got challenge: ", challenge)

	solution := c.SolveChallenge(logger.WithCtx(ctx, l, "solveChallenge"), challenge)

	l.Println("found solution: ", solution)

	return solution, challenge.Token, nil
}

// demoSubmitSolution is a demonstration of the basics of submitting solution and getting a quote.
func demoSubmitSolution(ctx context.Context, c *client.Client, solution uint, token string) error {
	l := logger.MustFromCtx(ctx)
	l.Println("demo submit solution")

	quote, err := c.SubmitSolution(logger.WithCtx(ctx, l, "submitSolution"), solution, token)
	if err != nil {
		return err
	}

	l.Println("got quote: ", quote)

	return nil
}

// demoLoad is a demonstration of the behavior of the server under load from the same token - it
// increases the difficulty of the challenge.
//
// The client will make 3 consecutive requests to the server with the same token.
func demoLoad(ctx context.Context, c *client.Client) {
	l := logger.MustFromCtx(ctx)
	l.Println("demo submit solution")

	for i := 0; i < 3; i++ {
		quote, err := c.GetQuote(logger.WithCtx(ctx, l, "getQuote"))
		if err != nil {
			l.Fatal(err)
		}

		l.Println("got quote: ", quote)
	}
}

func main() {
	ctx := context.Background()
	l := logger.New(nil, "demo")

	l.Println("initializing client")

	c := client.New(ctx, &client.Config{
		ServerHost: host,
		ServerPort: port,
	})

	solution, token, err := demoChallenge(logger.WithCtx(ctx, l, "demoChallenge"), c)
	if err != nil {
		l.Fatal(err)
	}

	if err := demoSubmitSolution(logger.WithCtx(ctx, l, "demoSubmitSolution"), c, solution, token); err != nil {
		l.Fatal(err)
	}

	demoLoad(logger.WithCtx(ctx, l, "demoLoad"), c)
}
