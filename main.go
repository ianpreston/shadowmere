package main

import (
	"./shadowmere"
	"log"
)

func main() {
    conf := shadowmere.LoadConfig("shadowmere.conf")
    err := shadowmere.ValidateConfig(conf)
    if err != nil {
        log.Fatalf(err.Error())
    }

	mere, err := shadowmere.NewServices(
        conf.Postgres.URI,
        conf.Link.Name,
        conf.Link.RemoteAddr,
        conf.Link.Password,
	)
	if err != nil {
		log.Fatal(err.Error())
	}

	mere.Start()
}
