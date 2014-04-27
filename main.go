package main

import (
	"./shadowmere"
	"log"
)

func main() {
	// TODO - Load configuration values from a file
	name := "noveria.0x-1.com"
	addr := "localhost:6668"
	pass := "foo"
	pgUrl := "postgres://localhost/shadowmere?sslmode=disable"

	datastore, err := shadowmere.NewDatastore(pgUrl)
	if err != nil {
		log.Fatal(err.Error())
	}

	server, err := shadowmere.NewServer(name, addr, pass, datastore)
	if err != nil {
		log.Fatal(err.Error())
	}

	server.Start()
}