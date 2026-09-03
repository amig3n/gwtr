package worktree

import (
	"encoding/json"
	"fmt"
	"os/exec"
	"os"
	"strings"
)

type StateStore struct {
	path string
}

func NewStateStore() (*StateStore, error) {
	// obtain repoistory root directory
	out, err := exec.Command("git", "rev-parse", "--show-toplevel").Output()
	if err != nil {
		return nil, err
	}
	// parse the bytes stream to proper path
	parsedPath := string(out)
	trimmedPath := strings.TrimSuffix(parsedPath, "\n")
	path := fmt.Sprintf("%s/.git/gwtr.json", trimmedPath)

	return &StateStore{
		path: path,
	}, nil
}

// NOTE create blank file if not exists, otherwise do nothing
func (s *StateStore) Init() error {
	_, err := os.Stat(s.path)	
	if os.IsNotExist(err) {
		// create the file
		file, err := os.Create(s.path)
		if err != nil {
			return fmt.Errorf("failed to initialize state file: %w", err)
		}
		defer file.Close()


	}

	return nil

}

func (s *StateStore) Load() ([]WorktreeState, error) {
	// read the file content
	stateFileContent, err := os.ReadFile(s.path)
	if os.IsNotExist(err) {
		// if file not found, init file and return empty slice
		err = s.Init()
		if err != nil {
			return nil, fmt.Errorf("failed to initialize state file: %w", err)
		}

		return []WorktreeState{}, nil
	}

	if err != nil {
		return nil, err
	}
	

	// if file not found, 
	
	// parse the file content to a slice of WorktreeState
	var state []WorktreeState
	err = json.Unmarshal(stateFileContent, &state)
	if err != nil {
		return nil, err
	}

	return state, nil
}

func (s *StateStore) Save(state []WorktreeState) error {
	// prepare current state as JSON
	data, err := json.Marshal(state)
	if err != nil {
		return err
	}

	// write the prepared data to parallel tmp file
	err = os.WriteFile(s.path+".tmp", data, 0644)
	if err != nil {
		return err
	}

	// NOTE: consider file syncing here

	// rename original file to backup
	err = os.Rename(s.path, s.path+".backup")
	if err != nil {
		return err
	}

	// rename tmp file to original file
	err = os.Rename(s.path+".tmp", s.path)
	if err != nil {
		return err
	}


	return nil
}
