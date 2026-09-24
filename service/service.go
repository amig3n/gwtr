package service

import (
	"log/slog"
	"github.com/amig3n/gwtr/worktree"
	"os"
	"path/filepath"
)



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

	worktrees, err := s.git.ListWorktrees()
	if err != nil {
		s.logger.Error("Failed to list worktrees from git provider", "error", err)
		return err
	}
	
	var initialState []RawState

	for _, wt := range worktrees {
		s.logger.Debug("Found worktree from git provider", "path", wt.Path, "branch", wt.Branch)
		
		initialState = append(initialState, RawState{
			Path: wt.Path,
			Deleted: false,
		})
	}

	s.logger.Debug("Initial state prepared", "wt_count", len(initialState))

	// transform git worktrees to raw state as it goes
	err = s.state.Save(initialState)
	if err != nil {
		s.logger.Error("Failed to save initial state file", "error", err)
		return err
	}

	s.logger.Debug("State file initialized successfully")

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

	// obtain worktree list from git provider
	wts, err := s.git.ListWorktrees()
	if err != nil {
		s.logger.Error("save state error: failed to list worktrees from git provider", "error", err)
		return err
	}

	// iteraate over state to determine required git actions
	for _, wt := range state.Items() {
		// check if given worktree exists in git repo	
		existsInGit := false

		for _, gitWt := range wts {
			if gitWt.Path == wt.Path {
				existsInGit = true
				break
			}
		}
		// if worktree exists in git and is not marked as deleted in state, do nothing
		if existsInGit && !wt.Deleted {
			s.logger.Debug("save state: git worktree: worktree exists and is up-to-date", "path", wt.Path, "branch", wt.Branch)
			continue
		}

		// if worktree exists in git but is marked as deleted in state, delete it from git (autofix after potential crash)
		if existsInGit && wt.Deleted {
			s.logger.Debug("saving state: git worktree: deleting worktree from git", "path", wt.Path, "branch", wt.Branch)
			err = s.git.DeleteWorktree(wt.Path)
			if err != nil {
				s.logger.Error("save state error: failed to delete worktree from git", "path", wt.Path, "branch", wt.Branch, "error", err)
				return err
			}
		}

		// if worktree does not exist in git but is not marked as deleted in state, add it to git
		if !existsInGit && !wt.Deleted {
			s.logger.Debug("saving state: git worktree: adding worktree to git", "path", wt.Path, "branch", wt.Branch)
			err = s.git.AddWorktree(wt.Branch, wt.Path)
			if err != nil {	
				s.logger.Error("save state error: failed to add worktree to git", "path", wt.Path, "branch", wt.Branch, "error", err)
				return err
			}
		}

		// if not exists in git and not marked as deleted - panic
		if !existsInGit && wt.Deleted {
			s.logger.Error("Save state error: worktree marked as deleted in state but does not exist in git", "path", wt.Path)
			// TODO this panic should be replaced with state fixing underneath
			panic("Save state error: state corrupted")
		}
	}

	// pop all worktrees from the end of the state file 
	// that are marked as deleted before saving
	cntDeleted := state.CleanDeleted()
	s.logger.Debug("Cleaned deleted worktrees from state", "deleted_count", cntDeleted)

	// transform worktreeState to rawState for saving
	// FIXME would be nice to have some method in worktreeState for this
	var newState []RawState
	for _, wt := range state.Items() {
		newState = append(newState, RawState{
			Path: wt.Path,
			Deleted: wt.Deleted,
		})
	}

	s.logger.Debug("Saving current state to state file")
	err = s.state.Save(newState)
	if err != nil {
		s.logger.Error("save state error: failed to save state file", "error", err)
		return err
	}

	return nil
}

// ANCHOR determine the target path form
func (s *Service) SanitizePath(inputPath string) (string, error) {
	// check if path is absolute or relative

	if !filepath.IsAbs(inputPath) {
		// if relative, convert to absolute with equal level to git repo
		s.logger.Debug("path is relative, converting to absolute", "inputPath", inputPath)

		// obtain repo root path from git provider (.git directory)
		repoPath, err := s.git.GetRepoRootPath()
		if err != nil {
			s.logger.Error("Failed to get repository root path", "error", err)
			return "", err
		}

		repoPath = filepath.Dir(repoPath) // get parent directory of .git

		// join repo path with input path
		// /home/user/repo + ../worktree1 => /home/user/worktree1
		inputPath = filepath.Join(repoPath, inputPath)
		s.logger.Debug("path converted to absolute", "absolutePath", inputPath)

	} else {
		// if absolute, don't mutate it
		s.logger.Debug("path is absolute", "inputPath", inputPath)
	}

	// check if path is valid
	s.logger.Debug("checking if path exists", "inputPath", inputPath)
	_, err := os.Stat(inputPath)
	if err != nil {
		if os.IsExist(err) {
			s.logger.Error("Path already exists", "inputPath", inputPath)
			return "", err
		}
	}
	s.logger.Debug("path is valid", "inputPath", inputPath)

	return inputPath, nil	
}

