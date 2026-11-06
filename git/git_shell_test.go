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

worktree /path/to/worktree3
HEAD sa092n031ni14b3b5b0i11
detached

`)

	worktrees, err := parsePorcelainOutput(output)
	if err != nil {
		t.Fatalf("Error parsing porcelain output: %v", err)
	}

	t.Logf("Parsed worktrees: %+v", worktrees)

	if len(worktrees) != 3 {
		t.Fatalf("Expected 3 worktrees, got %d", len(worktrees))
	}

	if worktrees[0].Path != "/path/to/worktree1" {
		t.Errorf("First worktree parsed incorrectly: %+v", worktrees[0])
	}

	if worktrees[0].Branch != "branch1" {
		t.Errorf("First worktree branch parsed incorrectly: %+v", worktrees[0])
	}

	if worktrees[1].Path != "/path/to/worktree2" {
		t.Errorf("Second worktree parsed incorrectly: %+v", worktrees[1])
	}

	if worktrees[1].Branch != "branch2" {
		t.Errorf("Second worktree branch incorrectly: %+v", worktrees[1])
	}

	if worktrees[2].Path != "/path/to/worktree3" || worktrees[2].Branch != "detached" {
		t.Errorf("Third worktree branch should be detached: %+v", worktrees[2])
	}

}
