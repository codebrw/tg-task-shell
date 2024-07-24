package main

import (
	"log"
	"log/slog"
	"tg-task-shell/config"
	"tg-task-shell/server"
)

func main() {
	slog.Info("init")

	app := server.New(config.NewConfig("TOKEN"))
	log.Panic( app.Start())
}