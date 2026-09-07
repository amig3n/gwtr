package service

import (
	"fmt"
	"log/slog"
)

// NOTE contract for git provider
type GitProvider interface {
	ListWorktrees() ([]GitWorktree, error)
	AddWorktree(branch string, path string) error
	DeleteWorktree(string) error
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

// NOTE helper function to prepare combined state -> List of worktrees combined with state
func (s *Service) prepareState() ([]WorktreeState, error) {
	s.logger.Debug("Loading state file")
	state, err := s.state.Load()
	if err != nil {
		return nil, fmt.Errorf("load state error: failed to load state file: %w", err)
	}
	s.logger.Debug("State file loaded successfully", "state", state)
	return state, nil
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
// TODO Rewrite this method to use slog instead of fmt.error
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
// TODO Rewrite this method to use slog instead of fmt.error
func (s *Service) Add(branch string, path string) error {
	// load up current state
	s.logger.Debug("Adding new worktree", "branch", branch, "path", path)
	currentState , err := s.state.Load()
	if err != nil {
		return fmt.Errorf("worktree add error: state file error: %w", err)
	}
	s.logger.Debug("Current state loaded successfully", "currentState", currentState)

	// add new worktree to path
	//TODO add using first deleted slot instead of always appending
	currentState = append(currentState, WorktreeState{
		Path: path,
		Deleted: false,
	})

	err = s.state.Save(currentState)
	if err != nil {
		return fmt.Errorf("worktree add error: state file error: %w", err)	
	}

	// add git worktree
	err = s.git.AddWorktree(branch, path)	
	if err != nil {
		return fmt.Errorf("worktree add error: git error: %w", err)
	}
		
	return nil
}

// helper function to return the index of worktree based on identification value
func determineWorktreeIndex(wtId string, state []AppWorktree) (int, error) {
	switch id := wtId.(type) {
		case int:
		// int assumes to be index
		// check if given id is in range of statefile table
			if id < 0 || id >= len(state) {			
				return -1, fmt.Errorf("worktree determine error: index out of range")
			} else {
				// id should be integer
				return id, nil
			}

		case string:
		// string assumes path or branch name
			for i, wt := range state {
				// iterate over loaded state till find a match
				if wt.Path == wtId || wt.Branch == wtId {
					return i, nil
				}
			}

			return -1, fmt.Errorf("worktree determine error: no worktree found for identifier: %s", wtId)

		default:
			//TODO check if this should panic
			// invalid type passed - should panic?
			return -1, fmt.Errorf("worktree determine error: invalid identifier type: %T", wtId)

	}
	return index, nil
}

// ANCHOR delete existing worktree
func (s *Service) Delete(wtId string) error {
	// load up current state
	s.logger.Debug("Deleting worktree", "wtId", wtId)
	currentState , err := s.state.Load()
	if err != nil {
		return fmt.Errorf("worktree add error: state file error: %w", err)
	}
	s.logger.Debug("Current state loaded successfully", "currentState", currentState)

	// remove worktree from statefile or mark it as deleted
	s.logger.Debug("Determining worktree based on identifier", "wtId", wtId)
	// TODO move this whole typeswitch to a separate helper function - will be used in switching as well
	// TODO ^ should always return index, and always the valid value
	// typeswitch
	switch id := wtId.(type) {
		// string assumes path or branch name
		case string:
			for i, wt := range currentState {
				if wt.Path == wtId {
					// path matched
					s.logger.Debug("WT detected based on path", "wtId", wtID, "path", wt.Path)
					currentState[i].Deleted = true
					break
				} else if wt.Branch == wtId {
					// branch matched
					s.logger.Debug("WT detected based on branch", "wtId", wtID, "branch", wt.Branch)
					currentState[i].Deleted = true
					break
				}
			}

		// int assumes to be index
		case int:
			s.logger.Debug("WT detected based on index", "wtId", wtId, "index", id)
			// check if index is in range of statefile table
			if id < 0 || id >= len(currentState) {
				return fmt.Errorf("worktree delete error: index out of range")
			}

			// if index is in range, mark the worktree as deleted
			currentState[id].Deleted = true	
			s.logger.Debug("WT marked as deleted based on index", "wtId", wtId, "index", id, "path", currentState[id].Path)

		// other types are not valid
		default:
			s.logger.Error("WT determine error: invalid identifier type", "wtId", wtId, "type", fmt.Sprintf("%T", wtId))
			return fmt.Errorf("worktree delete error: invalid identifier type: %T", wtId)
	}

	// TODO clean the statefile

	// save statefile
	s.logger.Debug("Saving updated state file after marking worktree as deleted")
	err = s.state.Save(currentState)
	if err != nil {
		return fmt.Errorf("worktree add error: state file error: %w", err)	
	}
	s.logger.Debug("State file saved successfully")
	
	// remove worktree from git
	s.logger.Debug("Removing worktree from git", "wtId", wtId)
	err = s.git.RemoveWorktree(wtId)
	if err != nil {
		return fmt.Errorf("worktree delete error: git error: %w", err)
	}
	s.logger.Debug("Worktree removed from git successfully", "wtId", wtId)

	return nil
}
