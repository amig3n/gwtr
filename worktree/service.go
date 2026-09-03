package worktree

import (
	"fmt"
)

// NOTE contract for git provider
type GitProvider interface {
	ListWorktrees() ([]GitWorktree, error)
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
	git GitProvider
	state StateStore
}

func NewService(git GitProvider, state StateStore) *Service {
	return &Service{
		git: git,
		state: state,
	}
}

func (s *Service) Init() error {
	// init state file with current workspaces
	err := s.state.Init()
	if err != nil {
		return fmt.Errorf("init error: failed to create state file: %w", err)
	}

	// get current state file
	worktrees, err := s.git.ListWorktrees()
	if err != nil {	
		return fmt.Errorf("init error: failed to list worktrees: %w", err)
	}

	// prepare the state file content
	var state []WorktreeState
	for _, wt := range worktrees {
		state = append(state, WorktreeState{
			Path: wt.Path,
			Deleted: false,
		})
	}

	// save the state file
	err = s.state.Save(state)
	if err != nil {
		return fmt.Errorf("init error: failed to save state file: %w", err)
	}

	return nil
}

func (s *Service) List() ([]AppWorktree, error) {
	// load the state file 
	worktrees, err := s.git.ListWorktrees()
	if err != nil {
		return nil, fmt.Errorf("list error: failed to list worktrees: %w", err)
	}

	//load the state file
	stateFile, err := s.state.Load()
	if err != nil {
		return nil, fmt.Errorf("list error: failed to load state file: %w", err)
	}

	// if state file is empty, it must be initialized
	if len(stateFile) == 0 {
		return nil, fmt.Errorf("list error: please run 'gwtr init' to initialize the state file")
	}
	
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

	return appWorktrees, nil
}

