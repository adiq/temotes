package main

import (
	"log"
	"temotes/temotes/api"
)

func main() {
	app := api.SetupServer()
	log.Fatal(app.Listen("0.0.0.0:80"))
}
