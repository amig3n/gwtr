package main
//
//import (
//	"log"
//	"log/slog"
//	"path/filepath"
//
//	"github.com/amig3n/gwtr/cli"
//	"github.com/amig3n/gwtr/git"
//	"github.com/amig3n/gwtr/service"
//	"github.com/amig3n/gwtr/state"
//)
//
//type App struct {
//	Logger *slog.Logger
//	cli	*cli.CLI
//}
//
//// ANCHOR App constructor
//func NewApp() *App {
//	// TODO allow to use flag for debug level
//	// init whole logger with debug level
//	logger := slog.New(
//		slog.NewTextHandler(
//			log.Writer(), 
//			&slog.HandlerOptions{
//				Level: slog.LevelDebug,
//			},
//		),
//	)
//
//	// init Service
//	logger.Debug("Initializing Service")
//	gitProvider := git.NewGitShellWrapper(logger)
//	statePath, err := gitProvider.GetRepoRootPath() 
//	if err != nil {
//		logger.Error("Error getting repository root path", "error", err)
//		panic(err)
//	}
//
//	// NOTE state file location, based on git repository root path
//	statePath = statePath + "/gwtr.json"
//
//	stateStore, err := state.NewJsonStateStore(statePath, logger)
//	if err != nil {
//		logger.Error("Error initializing state store", "error", err)
//		return nil
//	}
//	logger.Debug("State store initialized: ", "store", stateStore)
//
//	appService := service.NewService(logger, gitProvider, stateStore)
//	logger.Debug("Service initialized: ", "service", appService)
//
//	// pass initiated service to CLI
//	return &App{
//		Logger: logger,
//		cli: cli.NewCLI(appService),
//	}
//}
//
//// ANCHOR app runner
//func (app *App) Run() {
//	err := app.cli.RootCmd.Execute()
//	if err != nil {
//		app.Logger.Error("Error executing command", "error", err)
//	}
//}
//
//
//func InitizalizeService(logger *slog.Logger) (*service.Service, error) {
//
//	// init git provider and obtain repo path
//	logger.Debug("Initializing Service")
//	gitProvider := git.NewGitShellWrapper(logger)
//	statePath, err := gitProvider.GetRepoRootPath() 
//	if err != nil {
//		logger.Error("Error getting repository root path", "error", err)
//		panic(err)
//	}
//
//	stateFilePath := filepath.Join(statePath, "gwtr.json")
//
//	// init state store
//	stateStore, err := state.NewJsonStateStore(stateFilePath, logger)
//	if err != nil {
//		logger.Error("Error initializing state store", "error", err)
//		return nil, err
//	}
//	logger.Debug("State store initialized: ", "store", stateStore)
//
//	// construct service with logger, git provider and state store
//	appService := service.NewService(logger, gitProvider, stateStore)
//	logger.Debug("Service initialized: ", "service", appService)
//
//	return appService, nil
//}
