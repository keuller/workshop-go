package main

import (
	"log"

	"github.com/keuller/exchange/internal/api"
)

func main() {
	log.Println("[INFO] Exchange Service")

	app := api.New()
	go app.Start()
	go app.StartRpc()
	app.Stop()
}
