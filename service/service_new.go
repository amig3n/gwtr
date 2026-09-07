package service_new

import (
	"log/slog"
	"github.com/amig3n/worktree"
	"github/com/amig3n/state"
)


// NOTE contract for git provider
type GitProvider interface {
	ListWorktrees() ([]GitWorktree, error)
	AddWorktree(branch string, path string) error
	DeleteWorktree(string) error
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

// ANCHOR loading state from git and statefile
func (s *service) LoadState() (WorktreeState, error) {
	// TODO potential goroutines here

	// NOTE load state from statefile
	s.logger.Debug("Loading state file")
	stateFile, err := s.state.Load()
	if err != nil {
		s.logger.Error("Failed to load state file", "error", err)
		return nil, err
	}
	s.logger.Debug("State file loaded successfully", "state", state)

	// NOTE load worktrees from git provider
	gitWorktrees, err := s.git.ListWorktrees()
	if err != nil {
		s.logger.Error("Failed to list worktrees from git provider", "error", err)
		return nil, err
	}

	// NOTE combine state from statefile and git provider

	var combinedState WorktreeState
	for _, wt := range stateFile {
		// form up Worktree struct
		worktree := Worktree{
			path: wt.Path,
			deleted: wt.Deleted,
			branch: "",
			path: "",
		}

		// do not add any infos if wt is marked as deleted
		if !wt.Deleted { 
			// add proper infos from git
			for _, gitWt := range gitWorktrees {
				if gitWt.Path == wt.Path {
					worktree.branch = gitWt.Branch
					worktree.commit = gitWt.Commit
					break
				}
			}
		}

		// add to combined state
		combinedState.items = append(combinedState.items, worktree)
	}

	return state, nil
}

// ANCHOR Saving state to to proper domains
func (s *Service) SaveState(state WorktreeState) error {
	// iterate over state to prepare data for saving - convert to Statefile
	var newStateFile []RawState
	for _, wt := range state.Items() {
		// current branch info is discarded, as it's determined during loading
		currentWorktree := RawState{	
			path: wt.Path,
			deleted: wt.Deleted,
		}

		newStateFile = append(newStateFile, currentWorktree)
	}

	s.logger.Debug("Saving state file", "state", state)
	err := s.state.Save(state)
	if err != nil {
		return fmt.Errorf("save state error: failed to save state file: %w", err)
	}
	s.logger.Debug("State file saved successfully")

	// TODO parse state to save data in proper places
	return nil
}
