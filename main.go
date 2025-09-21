package main

import (
	"log"
	"temotes/temotes"
	"temotes/temotes/api"
)

func main() {
	app := api.SetupServer()
	log.Fatal(app.Listen(temotes.GetConfig().ServerAddr))
}
