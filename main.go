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
	pgUrl := "postgres://localhost/puzzle?sslmode=disable"

	datastore, err := puzzle.NewDatastore(pgUrl)
	if err != nil {
		log.Fatal(err.Error())
	}

	server, err := puzzle.NewServer(name, addr, pass, datastore)
	if err != nil {
		log.Fatal(err.Error())
	}

	server.Start()
}