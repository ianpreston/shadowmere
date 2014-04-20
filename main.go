package main

import (
	"./puzzle"
	"log"
)

func main() {
	server, err := puzzle.NewServer()
	if err != nil {
		log.Fatal(err.Error())
	}

	server.Start()
}