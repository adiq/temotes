package main

import (
	"log"
	"temotes/temotes"
	"temotes/temotes/api"
)

func main() {
	cfg := temotes.Load()
	temotes.SetConfig(cfg)

	app := api.SetupServer()

	log.Fatal(app.Listen(cfg.ServerAddr))
}
