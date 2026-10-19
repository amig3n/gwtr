package worktree

import (
	"path/filepath"
	"fmt"
	"os"
)

// NOTE this struct represent a single worktree with all metadata saved in statefile
type Worktree struct {
	Path string
	Branch string
	Deleted bool
}

func NewWorktree(path string) (*Worktree, error) {
	// check if given path is empty, if so return error
	if path == "" {
		return nil, fmt.Errorf("path cannot be empty")
	}

	// if given path is not absolute combine it with CWD
	if !filepath.IsAbs(path) {
		cwd, err := os.Getwd()
		if err != nil {
			return nil, fmt.Errorf("failed to get current working directory: %v", err)
		}
		path = filepath.Join(cwd, path)
	}

	// return object
	return &Worktree{
		Path: path,
		Branch: "",
		Deleted: false,
	}, nil
}

func (wt *Worktree) SetBranch(branch string) {
	wt.Branch = branch
}


func (wt *Worktree) Delete() {
	wt.Deleted = true
}
