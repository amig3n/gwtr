package git

import (
	"os/exec"
	"strings"
	"log/slog"
	"path"
	"errors"
	"github.com/amig3n/gwtr/service"
	"fmt"
)


// Implements GitProvider interface using git shell commands
type GitShellWrapper struct {}

func NewGitShellWrapper(logger *slog.Logger) *GitShellWrapper {
	return &GitShellWrapper{}
}

// ANCHOR helper function: parsing porcelain output from git
func parsePorcelainOutput(output []byte) ([]service.GitWorktree, error) {
	// parse the output to slice of strings
	lines := strings.Split(string(output), "\n")

	// prepare required resources
	var worktrees []service.GitWorktree
	var currentWorktree service.GitWorktree
	var inBlock bool = false

	// parse the lines to slice of GitWorktree
	for _, line := range lines {
		// if current line is blank, it means that block has been finished, so append the current worktree, reset the object and skip iteration
		if line == "" {
			if inBlock {
				worktrees = append(worktrees, currentWorktree)
				currentWorktree = service.GitWorktree{}
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
	}

	// close active block if present
	if inBlock {
		worktrees = append(worktrees, currentWorktree)
	}

	return worktrees, nil
}

// ANCHOR helper function: getting repository root path
func (wrp *GitShellWrapper) GetRepoRootPath() (string, error) {
	var rootPath string

	rootPathBytes, err := exec.Command("git", "rev-parse", "--path-format=absolute", "--git-common-dir").Output()
	if err != nil {
		var exitError *exec.ExitError

		if errors.As(err, &exitError) {
			if exitError.ExitCode() == 128 {
				return "", errors.New("not a git repository")
			}
		}

		return "", err
	}

	// trim the output to get clear path
	rootPath = strings.TrimSpace(string(rootPathBytes))
	return rootPath, nil
}

// ANCHOR GitProvider interface implementation
func (wrp *GitShellWrapper) ListWorktrees() ([]service.GitWorktree, error) {
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

func (wrp *GitShellWrapper) AddWorktree(branch string, wtPath string) error {
	// determine if path is absolute or relative, if relative, convert to absolute
	if !path.IsAbs(wtPath) {
		repoRootPath, err := wrp.GetRepoRootPath()
		if err != nil {
			return err
		}

		// this should return absolute path for the worktree
		wtPath = path.Join(repoRootPath, wtPath)
	}

	var cmd *exec.Cmd

	// TODO check if given branch already exists, and throw error if so
	checkBranchCmd := exec.Command("git", "show-ref", "--heads", branch, "--quiet")
	
	branchNotExist := checkBranchCmd.Run()

	// check if branch already exists, if so - add worktree without creating new branch
	if branchNotExist == nil {
		
		cmd = exec.Command("git", "worktree", "add", wtPath, branch)
	} else {
		
		cmd = exec.Command("git", "worktree", "add", "-b", branch, wtPath)
	}

	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to add worktree: %v, output: %s", err, string(output))
	}

	return nil
}

func (wrp *GitShellWrapper) DeleteWorktree(wtPath string) error {
	
	// TODO check if given path exists, if not - early fail
	cmd := exec.Command("git","worktree", "remove", wtPath)
	
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to add worktree: %v, output: %s", err, string(output))
	}

	return nil
}

