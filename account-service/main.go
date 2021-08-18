package main

import (
	"log"

	"github.com/keuller/account/internal/api"
	"github.com/keuller/account/internal/infra"
)

func main() {
	log.Println("[INFO] Account Service")
	if err := infra.InitDB(); err != nil {
		panic(err)
	}

	app := api.New()
	go app.Start()
	app.Stop()
}
