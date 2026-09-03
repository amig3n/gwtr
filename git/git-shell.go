package git

import (
	"os/exec"
	"strings"
	"github.com/amig3n/gwtr/worktree"
)


// Implements GitProvider interface using git shell commands
type GitShellWrapper struct {
}

func (repo *GitShellWrapper) ListWorktrees() ([]worktree.GitWorktree, error) {
	// execute git worktree list command and return the output
	// output is []Byte
	output, err := exec.Command("git", "worktree", "list", "--porcelain").Output()
	if err != nil {
		return nil, err
	}

	// parse the output to slice of strings
	lines := strings.Split(string(output), "\n")

	// parse the lines to slice of GitWorktree
	var worktrees []worktree.GitWorktree
	var currentWorktree worktree.GitWorktree
	for _, line := range lines {
		// if current line is blank, it means that block has been finished, so append the current worktree, reset the object and skip iteration
		if line == "" {
			worktrees = append(worktrees, currentWorktree)
			currentWorktree = worktree.GitWorktree{}
			continue
		}

		if strings.HasPrefix(line, "worktree ") {
			currentWorktree.Path = strings.TrimPrefix(line, "worktree ")
		}

		if strings.HasPrefix(line, "branch ") {
			currentWorktree.Branch = strings.TrimPrefix(line, "branch /refs/heads/")
		}

		if strings.HasPrefix(line, "HEAD ") {
			currentWorktree.Commit = strings.TrimPrefix(line, "HEAD ")
		}
	}

	return worktrees, nil
}

