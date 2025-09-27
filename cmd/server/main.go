package main

import (
	"gochat/internal/server"
	"log"
)

func main() {
	srv := server.NewServer(":8080")
	log.Fatal(srv.Start())
}
