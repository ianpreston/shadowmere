package shadowmere

import (
	"code.google.com/p/gcfg"
	"errors"
)

var InvalidLinkName = errors.New("Invalid Link Name")
var InvalidLinkAddr = errors.New("Invalid Link RemoteAddr")
var InvalidLinkPassword = errors.New("Invalid Link Password")
var InvalidPostgresURI = errors.New("Invalid Postgres URI")

type Config struct {
	Link struct {
		Name       string `gcfg:"name"`
		RemoteAddr string `gcfg:"remote-addr"`
		Password   string `gcfg:"password"`
	}

	Postgres struct {
		URI string `gcfg:"uri"`
	}
}

func LoadConfig(path string) Config {
	var config Config
	gcfg.ReadFileInto(&config, path)
	return config
}

func ValidateConfig(config Config) error {
	if len(config.Link.Name) == 0 {
		return InvalidLinkName
	}
	if len(config.Link.RemoteAddr) == 0 {
		return InvalidLinkAddr
	}
	if len(config.Link.Password) == 0 {
		return InvalidLinkPassword
	}
	if len(config.Postgres.URI) == 0 {
		return InvalidPostgresURI
	}
	return nil
}
