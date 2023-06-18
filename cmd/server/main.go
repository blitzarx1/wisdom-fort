package main

import (
	"log"

	"blitzarx1/wisdom-fort/server"
)

func main() {
	srv, err := server.New()
	if err != nil {
		log.Fatal(err)
	}
	srv.Run()
}
