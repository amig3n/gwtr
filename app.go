package main

import (
	"log/slog"
	"log"
)

type App struct {
	Logger *slog.Logger
}

func NewApp() *App {
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
	}
}
