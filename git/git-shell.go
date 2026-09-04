package git

import (
	"os/exec"
	"strings"
	"github.com/amig3n/gwtr/worktree"
	"log/slog"
	"path"
)


// Implements GitProvider interface using git shell commands
type GitShellWrapper struct {
	logger *slog.Logger
}

func NewGitShellWrapper(logger *slog.Logger) *GitShellWrapper {
	return &GitShellWrapper{
		logger: logger,
	}
}
// ANCHOR helper function: parsing porcelain output from git
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

// ANCHOR helper function: getting repository root path
func (wrp *GitShellWrapper) getRepoRootPath() (string, error) {
	wrp.logger.Debug("Getting repository root path")
	var rootPath string
	rootPathBytes, err := exec.Command("git", "rev-parse", "--show-toplevel").Output()
	if err != nil {
		wrp.logger.Error(
			"Failed to get repository root path",
			"error", err,
		)
		return "", err
	}

	rootPath = strings.TrimSpace(string(rootPathBytes))

	wrp.logger.Debug("Repository root path retrieved successfully", "rootPath", rootPath)
	return rootPath, nil
}

// ANCHOR GitProvider interface implementation
func (wrp *GitShellWrapper) ListWorktrees() ([]worktree.GitWorktree, error) {
	// execute git worktree list command and return the output
	// output is []Byte
	output, err := exec.Command("git", "worktree", "list", "--porcelain").Output()
	if err != nil {
		return nil, err
	}

	wrp.logger.Debug("Git worktree list output", "output", string(output))

	// parse the porcelain output to slice of GitWorktree
	wrp.logger.Debug("Parsing porcelain output")
	worktrees, err := parsePorcelainOutput(output)
	if err != nil {
		wrp.logger.Error("Failed to parse porcelain output", "error", err)
		return nil, err
	}

	wrp.logger.Debug("Parsed worktrees", "worktrees", worktrees)

	return worktrees, nil
}

func (wrp *GitShellWrapper) AddWorktree(branch string, wtPath string) error {
	// execute git worktree add command by adding new branch
	wrp.logger.Debug("Adding worktree", "branch", branch, "newWtPath", wtPath)

	// determine if path is absolute or relative, if relative, convert to absolute
	if !path.IsAbs(wtPath) {
		repoRootPath, err := wrp.getRepoRootPath()
		if err != nil {
			wrp.logger.Error(
				"Failed to get repository root path",
				"error", err,
			)
			return err
		}

		wtPath = path.Join(repoRootPath, wtPath)
	}

	// TODO check if given branch already exists, if yes, return error

	cmd := exec.Command("git", "worktree", "add", "-b", branch, wtPath)
	wrp.logger.Debug("Executing command", "command", cmd.String())

	output, err := cmd.CombinedOutput()
	if err != nil {
		wrp.logger.Error(
			"Git worktree add command failed", 
			"error", err, 
			"output", string(output),
		)
		return err
	}

	return nil
}


