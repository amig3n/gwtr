package worktree

import (
	"testing"
)

func TestWorktreeState_Add(t *testing.T) {
	state := NewWorktreeState([]Worktree{})
	
	worktree1 := Worktree{
		Path: "test/path",
		Branch: "test-branch",
		Deleted: false,
	}

	state.Add(worktree1)

	if len(state.Items()) != 1 {
		t.Errorf("Expected 1 worktree in state, got %d", len(state.Items()))
	}
}


//TODO implement test for deleted indexes detection
