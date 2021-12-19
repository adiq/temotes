package main

import (
	"log"
	"os"
	"temotes/temotes/api"
)

func main() {
	app := api.SetupServer()
	log.Fatal(app.Listen(os.Getenv("SERVER_ADDR")))
}
