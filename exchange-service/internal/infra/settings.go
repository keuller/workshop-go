package infra

import (
	"log"

	"github.com/golobby/config/v3"
	"github.com/golobby/config/v3/pkg/feeder"
)

type Settings struct {
	Host string
	Port string
}

func NewSettings() Settings {
	settings := Settings{}
	jsonFeed := feeder.Json{Path: "config/server.json"}
	c := config.New()
	c.AddFeeder(jsonFeed).AddStruct(&settings)
	if err := c.Feed(); err != nil {
		log.Println("[ERROR] fail to load settings file -", err.Error())
	}
	return settings
}
