package main

import (
	"embed"
	"fmt"
	"log"
	"os"
	"os/signal"

	"github.com/joho/godotenv"

	server "dependents.info/cmd"
	"dependents.info/internal/config"
	"dependents.info/internal/service"
	"dependents.info/pkg/utils"
)

//go:embed static/*
var static embed.FS

func main() {
	log.SetFlags(0)
	godotenv.Load()
	utils.LoadStylesFile(&static)

	cfg := config.New()
	services := service.BuildAll(cfg)
	app := server.Build(cfg, &static, services)

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		<-c
		_ = app.Shutdown()
	}()

	if err := app.Listen(fmt.Sprint(":", cfg.Port)); err != nil {
		log.Panic(err)
	}

	services.DatabaseService.Sync()
	services.DatabaseService.Close()
}
