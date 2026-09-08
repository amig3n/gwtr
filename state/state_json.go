package state

import (
	"encoding/json"
	"fmt"
	"os"
	"log/slog"
	// imports StateStore and RawState from service package to avoid circular dependency
	"github.com/amig3n/gwtr/service"
)

type JsonStateStore struct {
	path string
	logger *slog.Logger
}

func NewJsonStateStore(path string, logger *slog.Logger) (*JsonStateStore, error) {
	
	// TODO check if given path is valid

	return &JsonStateStore{
		path: path,
		logger: logger,
	}, nil
}


// NOTE create blank file if not exists, otherwise do nothing
func (s *JsonStateStore) Init() error {
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

func (s *JsonStateStore) Load() ([]service.RawState, error) {
	// read the file content
	stateFileContent, err := os.ReadFile(s.path)
	if os.IsNotExist(err) {
		// if file not found, init file and return empty slice
		err = s.Init()
		if err != nil {
			return nil, fmt.Errorf("failed to initialize state file: %w", err)
		}

		return []service.RawState{}, nil
	}

	if err != nil {
		return nil, err
	}
	
	// if file not found, 
	// parse the file content to a slice of RawState
	var state []service.RawState
	err = json.Unmarshal(stateFileContent, &state)
	if err != nil {
		return nil, err
	}

	return state, nil
}

func (s *JsonStateStore) Save(state []service.RawState) error {
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
