package main

import (
	"context"

	"blitzarx1/wisdom-fort/pkg/logger"
	"blitzarx1/wisdom-fort/server"
)

func main() {
	l := logger.New(nil, "main")
	l.Println("initializing server")

	cfg, err := server.NewConfigFromEnv()
	if err != nil {
		l.Fatal(err)
	}

	ctx := context.Background()

	srv, err := server.New(logger.WithCtx(ctx, l, "serverNew"), cfg)
	if err != nil {
		l.Fatal(err)
	}

	if err := srv.Run(logger.WithCtx(ctx, l, "serverRun")); err != nil {
		l.Fatal(err)
	}
}
