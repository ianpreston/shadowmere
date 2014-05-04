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

	mere, err := shadowmere.NewServices(
		pgUrl,
		name,
		addr,
		pass,
	)
	if err != nil {
		log.Fatal(err.Error())
	}

	mere.Start()
}