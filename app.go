package main

import (
	"log/slog"
	"log"
	"github.com/amig3n/gwtr/cli"
)

type App struct {
	Logger *slog.Logger
	cli	*cli.CLI
}

func NewApp() *App {
	// TODO allow to use flag for debug level
	// init whole logger with debug level
	logger := slog.New(
		slog.NewTextHandler(
			log.Writer(), 
			&slog.HandlerOptions{
				Level: slog.LevelDebug,
			},
		),
	)

	return &App{
		Logger: logger,
		cli: cli.NewCLI(),
	}
}

func (app *App) Run() {
	app.Logger.Info("GWTR started")	

	err := app.cli.RootCmd.Execute()
	if err != nil {
		app.Logger.Error("Error executing command", "error", err)
	}
}
