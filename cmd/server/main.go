package main

import (
	"context"
	"log"

	"blitzarx1/wisdom-fort/server"
)

func main() {
	srv, err := server.New(context.Background())
	if err != nil {
		log.Fatal(err)
	}
	srv.Run(context.Background())
}
