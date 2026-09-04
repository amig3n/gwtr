package main

import (
	"log/slog"
	"log"
	"github.com/amig3n/gwtr/cli"
	"github.com/amig3n/gwtr/worktree"
	"github.com/amig3n/gwtr/git"
)

type App struct {
	Logger *slog.Logger
	cli	*cli.CLI
}

// ANCHOR App constructor
func NewApp() *App {
	// TODO allow to use flag for debug level
	// init whole logger with debug level
	logger := slog.New(
		slog.NewTextHandler(
			log.Writer(), 
			&slog.HandlerOptions{
				Level: slog.LevelInfo,
			},
		),
	)

	// init Service
	logger.Debug("Initializing Service")
	gitProvider := git.NewGitShellWrapper(logger)

	stateStore, err := worktree.NewStateStore(logger)
	if err != nil {
		logger.Error("Error initializing state store", "error", err)
		return nil
	}
	logger.Debug("State store initialized: ", "store", stateStore)

	appService := worktree.NewService(logger, gitProvider, *stateStore)
	logger.Debug("Service initialized: ", "service", appService)

	// pass initiated service to CLI
	return &App{
		Logger: logger,
		cli: cli.NewCLI(appService),
	}
}

// ANCHOR app runner
func (app *App) Run() {
	err := app.cli.RootCmd.Execute()
	if err != nil {
		app.Logger.Error("Error executing command", "error", err)
	}
}

