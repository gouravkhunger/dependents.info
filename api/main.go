package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"

	"github.com/joho/godotenv"

	server "dependents-img/cmd"
	"dependents-img/internal/config"
	"dependents-img/internal/service"
)

func main() {
	log.SetFlags(0)
	godotenv.Load()

	cfg := config.New()
	services := service.BuildAll(cfg)
	app := server.Build(cfg, services)

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		<-c
		_ = app.Shutdown()
	}()

	if err := app.Listen(fmt.Sprintf(":%s", cfg.Port)); err != nil {
		log.Panic(err)
	}

	services.DatabaseService.Close()
}
