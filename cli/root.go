package cli

import (
	"fmt"
	"log"
	"log/slog"
	"path/filepath"

	"github.com/spf13/cobra"

	"github.com/amig3n/gwtr/git"
	"github.com/amig3n/gwtr/service"
	"github.com/amig3n/gwtr/state"
)

type CLI struct {
	RootCmd *cobra.Command
	logger *slog.Logger
	level *slog.LevelVar
}

func (c *CLI) PrepareService(logger *slog.Logger) (*service.Service, error) {
	// init git provider and obtain repo path
	logger.Debug("Initializing Service")
	gitProvider := git.NewGitShellWrapper(logger)
	statePath, err := gitProvider.GetRepoRootPath() 
	if err != nil {
		return nil, err
	}

	stateFilePath := filepath.Join(statePath, "gwtr.json")

	// init state store
	stateStore, err := state.NewJsonStateStore(stateFilePath)
	if err != nil {
		logger.Error("Error initializing state store", "error", err)
		return nil, err
	}
	logger.Debug("State store initialized: ", "store", stateStore)

	// construct service with logger, git provider and state store
	appService := service.NewService(logger, gitProvider, stateStore)
	logger.Debug("Service initialized: ", "service", appService)

	return appService, nil
}

func NewCLI() *CLI {
	// define loglevel holder
	level := new(slog.LevelVar)
	level.Set(slog.LevelWarn) // default log level

	cli := CLI{
		level: level,
	}

	// initialize logger with default settings
	cli.logger = slog.New(
		slog.NewTextHandler(
			log.Writer(), 
			&slog.HandlerOptions{
				Level: level,
			},
		),
	)

	var loglevel string

	// init CLI
	rootCmd := &cobra.Command{
		Use:   "gwtr",
		Short: "GWTR",
		Long:  "Worktres manager for git. For anybody who works with multiple branches at one time.",

		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			switch loglevel {
				case "debug":
					cli.level.Set(slog.LevelDebug)
				case "info":
					cli.level.Set(slog.LevelInfo)
				case "warn":
					cli.level.Set(slog.LevelWarn)
				case "error":
					cli.level.Set(slog.LevelError)
				default:
					return fmt.Errorf("invalid log level: %s", loglevel)
			}
			return nil
		},
	}

	// add global flags
	rootCmd.PersistentFlags().StringVar(
		&loglevel,
		"log-level",
		"warn",
		"Set the logging level (debug, info, warn, error)",
	)

	// add subcommands
	rootCmd.AddCommand(cli.addListCmd())
	rootCmd.AddCommand(cli.addInitCmd())
	rootCmd.AddCommand(cli.addAddCmd())
	rootCmd.AddCommand(cli.addRemoveCmd())

	// assign prepared command tree to CLI object
	cli.RootCmd = rootCmd
	
	return &cli
}

