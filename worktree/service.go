package worktree

import (
	"fmt"
	"log/slog"
)

// NOTE contract for git provider
type GitProvider interface {
	ListWorktrees() ([]GitWorktree, error)
	AddWorktree(branch string, path string) error
}

// NOTE data model that comes straight from git provider
type GitWorktree struct {
	Path string
	Branch string
	Commit string // FIXME to be deleted?
}


// NOTE data model used by the state file
type WorktreeState struct {
	Path string `json:path`
	Deleted bool `json:deleted`
}

// NOTE data model that comes from combination between git and state
type AppWorktree struct {
	Path string
	Branch string
	Commit string // FIXME to be deleted?
	Deleted bool 
}

type Service struct {
	logger *slog.Logger
	git GitProvider
	state StateStore
}

func NewService(logger *slog.Logger, git GitProvider, state StateStore) *Service {
	return &Service{
		logger: logger,
		git: git,
		state: state,
	}
}

// ANCHOR listing all worktrees
func (s *Service) Init() error {
	// init state file with current workspaces
	s.logger.Debug("Initializing state file")
	err := s.state.Init()
	if err != nil {
		return fmt.Errorf("init error: failed to create state file: %w", err)
	}
	s.logger.Debug("State file initialized successfully")

	// get current state file
	s.logger.Debug("Getting current worktrees")
worktrees, err := s.git.ListWorktrees()
	if err != nil {	
		return fmt.Errorf("init error: failed to list worktrees: %w", err)
	}
	s.logger.Debug("Current worktrees retrieved successfully", "worktrees", worktrees)

	// prepare the state file content
	s.logger.Debug("Preparing state file content")
	var state []WorktreeState
	for _, wt := range worktrees {
		s.logger.Debug("Adding worktree to state", "worktree", wt)
		state = append(state, WorktreeState{
			Path: wt.Path,
			Deleted: false,
		})
	}
	s.logger.Debug("State file content prepared successfully", "state", state)

	// save the state file
	s.logger.Debug("Saving state file")
	err = s.state.Save(state)
	if err != nil {
		return fmt.Errorf("init error: failed to save state file: %w", err)
	}
	s.logger.Debug("State file saved successfully")

	return nil
}

// ANCHOR listing all worktrees
func (s *Service) List() ([]AppWorktree, error) {
	// load the state file 
	s.logger.Debug("Listing worktrees")
	worktrees, err := s.git.ListWorktrees()
	if err != nil {
		return nil, fmt.Errorf("list error: failed to list worktrees: %w", err)
	}
	s.logger.Debug("Git worktrees retrieved successfully", "worktrees", worktrees)

	//load the state file
	s.logger.Debug("Loading state file")
	stateFile, err := s.state.Load()
	if err != nil {
		return nil, fmt.Errorf("list error: failed to load state file: %w", err)
	}
	s.logger.Debug("State file loaded successfully", "stateFile", stateFile)

	// if state file is empty, it must be initialized
	s.logger.Debug("Checking if state file is empty")
	if len(stateFile) == 0 {
		return nil, fmt.Errorf("list error: please run 'gwtr init' to initialize the state file")
	}
	s.logger.Debug("State file is not empty, matching worktrees with state file")
	
	// match the gitWorktree with State using path as key
	var appWorktrees []AppWorktree
	// NOTE iterate over statefile index to easily preserve the order
	for _, wt := range stateFile {
		var currentWorktree AppWorktree	= AppWorktree{
			Path: wt.Path,
			Deleted: wt.Deleted,
		}

		// iterate over git worktrees to find the matching path
		for _, gitWorktree := range worktrees {
			// if the path matches, get branch and commit
			if wt.Path == gitWorktree.Path {
				currentWorktree.Branch = gitWorktree.Branch
				currentWorktree.Commit = gitWorktree.Commit
				break
			}
		}

		// TODO: remove the found worktree from worktrees to speed up matching

		// TODO: handle remaining worktrees that are not matched to state file

		// append prepared object to result list
		appWorktrees = append(appWorktrees, currentWorktree)

	}

	s.logger.Debug("Worktrees matched with state file successfully", "appWorktrees", appWorktrees)

	return appWorktrees, nil
}

// ANCHOR adding new worktree
func (s *Service) Add(branch string, path string) error {
	// load up current state
	s.logger.Debug("Adding new worktree", "branch", branch, "path", path)
	currentState , err := s.state.Load()
	if err != nil {
		return fmt.Errorf("worktree add error: state file error: %w", err)
	}
	s.logger.Debug("Current state loaded successfully", "currentState", currentState)


	// add git worktree
	err = s.git.AddWorktree(branch, path)	
	if err != nil {
		return fmt.Errorf("worktree add error: git error: %w", err)
	}

	currentState = append(currentState, WorktreeState{
		Path: path,
		Deleted: false,
	})

	err = s.state.Save(currentState)
	if err != nil {
		return fmt.Errorf("worktree add error: state file error: %w", err)	
	}
		

	return nil
}
