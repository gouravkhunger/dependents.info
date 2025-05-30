package main

import (
	app "dependents-img/cmd"
	"github.com/joho/godotenv"
	"log"
)

func main() {
	log.SetFlags(0)
	godotenv.Load()
	log.Fatal(app.Start())
}
