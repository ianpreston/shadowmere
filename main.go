package main

import (
	"./puzzle"
	"log"
)

func main() {
	// TODO - Load configuration values from a file
	name := "noveria.0x-1.com"
	addr := "localhost:6667"
	pass := "foo"

	server, err := puzzle.NewServer(name, addr, pass)
	if err != nil {
		log.Fatal(err.Error())
	}

	server.Start()
}