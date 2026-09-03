package git

import (
	"testing"
)

func TestPorcelainParsing(t *testing.T) {
	output := []byte(`worktree /path/to/worktree1
branch refs/heads/branch1
HEAD 1234567890abcdef

worktree /path/to/worktree2
branch refs/heads/branch2
HEAD abcdef1234567890

`)

	worktrees, err := parsePorcelainOutput(output)
	if err != nil {
		t.Fatalf("Error parsing porcelain output: %v", err)
	}

	t.Logf("Parsed worktrees: %+v", worktrees)

	if len(worktrees) != 2 {
		t.Fatalf("Expected 2 worktrees, got %d", len(worktrees))
	}

	if worktrees[0].Path != "/path/to/worktree1" || worktrees[0].Branch != "branch1" || worktrees[0].Commit != "1234567890abcdef" {
		t.Errorf("Unexpected first worktree: %+v", worktrees[0])
	}

	if worktrees[1].Path != "/path/to/worktree2" || worktrees[1].Branch != "branch2" || worktrees[1].Commit != "abcdef1234567890" {
		t.Errorf("Unexpected second worktree: %+v", worktrees[1])
	}

}
