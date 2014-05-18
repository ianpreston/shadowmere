package main

import (
	"./shadowmere"
    "./kenny"
)

func main() {
    conf := shadowmere.LoadConfig("shadowmere.conf")
    err := shadowmere.ValidateConfig(conf)
    if err != nil {
        kenny.Fatal(err)
    }

    kenny.Info("Hello, world!")

	mere, err := shadowmere.NewServices(
        conf.Postgres.URI,
        conf.Link.Name,
        conf.Link.RemoteAddr,
        conf.Link.Password,
	)
	if err != nil {
		kenny.Fatal(err)
	}

	mere.Start()
}
