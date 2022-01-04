package infra

import (
	"log"

	"github.com/golobby/config/v3"
	"github.com/golobby/config/v3/pkg/feeder"
)

var (
	settings Settings = Settings{}
	isLoaded bool     = false
)

type Settings struct {
	Host        string
	Port        string
	ExchangeUrl string
	DbFile      string
}

func _loadConfigFile() {
	log.Println("[INFO] Loading configuration file.")
	jsonFeed := feeder.Json{Path: "config/server.json"}
	c := config.New()
	c.AddFeeder(jsonFeed).AddStruct(&settings)
	if err := c.Feed(); err != nil {
		log.Println("[ERROR] fail to load settings file -", err.Error())
	}
}

// GetConfig read a key configuration in 'server.json' file
func GetConfig() Settings {
	if isLoaded {
		return settings
	}

	_loadConfigFile()

	isLoaded = true
	return settings
}
