package service

import (
	"log/slog"
	"github.com/amig3n/gwtr/worktree"
)


// NOTE data model to aggregate from git provider
type GitWorktree struct {
	Path string
	Branch string
}

// NOTE data model used by the state file
type RawState struct {
	Path    string `json:"path"`
	Deleted bool   `json:"deleted"`
}

// NOTE contract for git provider
type GitProvider interface {
	ListWorktrees() ([]GitWorktree, error)
	AddWorktree(branch string, path string) error
	DeleteWorktree(string) error
}

// NOTE contract for state Store
type StateStore interface {
	Load() ([]RawState, error)
	Save([]RawState) error
	Init() error
}

// NOTE infrastructure-based composition root
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

func CombineState(worktrees []GitWorktree, rawState []RawState) (worktree.WorktreeState, error) {
	var combinedState worktree.WorktreeState
	for _, wt := range rawState {
		// form up Worktree struct
		worktree := worktree.Worktree{
			Path: wt.Path,
			Deleted: wt.Deleted,
			Branch: "",
		}

		// do not add any infos if wt is marked as deleted
		if !wt.Deleted { 
			// add proper infos from git
			for _, gitWt := range worktrees {
				if gitWt.Path == wt.Path {
					worktree.Branch = gitWt.Branch
					break
				}
			}
		}

		// add to combined state
		combinedState.Add(worktree)
	}

	//TODO check if any WT from git is not connected with statefile

	return combinedState, nil
}

func (s *Service) InitState() error {
	s.logger.Debug("Initializing state file")
	err := s.state.Init()
	if err != nil {
		s.logger.Error("Failed to initialize state file", "error", err)
		return err
	}

	return nil
}


// ANCHOR loading state from git and statefile
func (s *Service) LoadState() (worktree.WorktreeState, error) {
	// TODO potential goroutines here

	// NOTE load state from statefile
	s.logger.Debug("Loading state file")
	stateFile, err := s.state.Load()
	if err != nil {
		s.logger.Error("Failed to load state file", "error", err)
		return worktree.WorktreeState{}, err
	}
	s.logger.Debug("State file loaded successfully")

	// NOTE load worktrees from git provider
	gitWorktrees, err := s.git.ListWorktrees()
	if err != nil {
		s.logger.Error("Failed to list worktrees from git provider", "error", err)
		return worktree.WorktreeState{}, err
	}

	readyState, err := CombineState(gitWorktrees, stateFile)
	if err != nil {
		s.logger.Error("Failed to combine state from git and statefile", "error", err)
		return worktree.WorktreeState{}, err
	}

	return readyState, nil
}

// ANCHOR Saving state to to proper domains
func (s *Service) SaveState(state worktree.WorktreeState) error {
	// iterate over state to prepare data for saving - convert to Statefile
	var newStateFile []RawState
	for _, wt := range state.Items() {
		// current branch info is discarded, as it's determined during loading
		currentWorktree := RawState{	
			Path: wt.Path,
			Deleted: wt.Deleted,
		}

		newStateFile = append(newStateFile, currentWorktree)
	}

	s.logger.Debug("Saving state file", "state", state)

	err := s.state.Save(newStateFile)
	if err != nil {
		s.logger.Error("save state error: failed to save state file", "error", err)
		return err
	}

	s.logger.Debug("State file saved successfully")

	// TODO parse state to save data in proper places
	return nil
}
