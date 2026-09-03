package git

import (
	"os/exec"
	"strings"
	"github.com/amig3n/gwtr/worktree"
)


// Implements GitProvider interface using git shell commands
type GitShellWrapper struct {
}

func parsePorcelainOutput(output []byte) ([]worktree.GitWorktree, error) {
	// parse the output to slice of strings
	lines := strings.Split(string(output), "\n")

	// prepare required resources
	var worktrees []worktree.GitWorktree
	var currentWorktree worktree.GitWorktree
	var inBlock bool = false

	// parse the lines to slice of GitWorktree
	for _, line := range lines {
		// if current line is blank, it means that block has been finished, so append the current worktree, reset the object and skip iteration
		if line == "" {
			if inBlock {
				worktrees = append(worktrees, currentWorktree)
				currentWorktree = worktree.GitWorktree{}
				inBlock = false
			}
			continue
		}

		if strings.HasPrefix(line, "worktree ") {
			currentWorktree.Path = strings.TrimPrefix(line, "worktree ")
			inBlock = true
		}

		if strings.HasPrefix(line, "branch ") {
			currentWorktree.Branch = strings.TrimPrefix(line, "branch refs/heads/")
			inBlock = true
		}

		if strings.HasPrefix(line, "HEAD ") {
			currentWorktree.Commit = strings.TrimPrefix(line, "HEAD ")
			inBlock = true
		}
	}

	// close active block if present
	if inBlock {
		worktrees = append(worktrees, currentWorktree)
	}

	return worktrees, nil

}

func (repo *GitShellWrapper) ListWorktrees() ([]worktree.GitWorktree, error) {
	// execute git worktree list command and return the output
	// output is []Byte
	output, err := exec.Command("git", "worktree", "list", "--porcelain").Output()
	if err != nil {
		return nil, err
	}

	// parse the porcelain output to slice of GitWorktree
	worktrees, err := parsePorcelainOutput(output)
	if err != nil {
		return nil, err
	}

	return worktrees, nil
}


