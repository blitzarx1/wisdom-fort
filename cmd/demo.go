package main

import (
	"context"
	"sync"

	"blitzarx1/wisdom-fort/pkg/logger"
	"blitzarx1/wisdom-fort/server"
)

func startServer(ctx context.Context, wg *sync.WaitGroup) {
	defer wg.Done()
	l := logger.MustFromCtx(ctx)
	l.Println("initializing server")

	cfg, err := server.NewConfigFromEnv()
	if err != nil {
		l.Fatal(err)
	}

	srv, err := server.New(logger.WithCtx(ctx, l, "serverNew"), cfg)
	if err != nil {
		l.Fatal(err)
	}

	srv.Run(logger.WithCtx(ctx, l, "serverRun"))
}

func main() {
	wg := &sync.WaitGroup{}
	defer wg.Wait()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	l := logger.New(nil, "demo")

	wg.Add(1)
	go startServer(logger.WithCtx(ctx, l, "server"), wg)
}
