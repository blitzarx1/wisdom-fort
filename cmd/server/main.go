package main

import (
	"context"
	"log"

	"blitzarx1/wisdom-fort/server"
)

func main() {
	cfg, err := server.NewConfigFromEnv()
	if err != nil {
		log.Fatal(err)
	}

	srv, err := server.New(context.Background(), cfg)
	if err != nil {
		log.Fatal(err)
	}

	srv.Run(context.Background())
}
