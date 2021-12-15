package main

import (
	"log"
	"temotes/temotes/api"
)

func main() {
	app := api.SetupServer()
	log.Fatal(app.Listen(":3000"))
}
