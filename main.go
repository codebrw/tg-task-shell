package main

import (
	"log"
	"log/slog"

	"tg-task-shell/config"
	"tg-task-shell/server"
	"tg-task-shell/shell"
)

func main() {
	config := config.Get()
	if config == nil {
		log.Fatal("failed to parse config")
	}

	slog.Info("starting server")
	app := server.New(config)
	
	slog.Info("starting shell executor")
	cmdChan := shell.Start()

	log.Panic(app.Start(cmdChan))
}
