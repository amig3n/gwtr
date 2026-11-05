package worktree

import (
	"testing"
)

func TestWorktreeStateAdd(t *testing.T) {
	state := NewWorktreeState([]Worktree{})
	
	worktree1 := Worktree{
		Path: "test/path",
		Branch: "test-branch",
		Deleted: false,
	}

	state.Add(worktree1, false)

	if len(state.Items()) != 1 {
		t.Errorf("Expected 1 worktree in state, got %d", len(state.Items()))
	}
}

func TestWorktreeStateReuseIndex(t *testing.T) {
	// prepare state
	var state []Worktree = []Worktree{
		Worktree{Path: "path1", Branch: "branch1", Deleted: false},
		Worktree{Path: "path2", Branch: "branch2", Deleted: true}, // this one is marked as deleted
		Worktree{Path: "path3", Branch: "branch3", Deleted: false},
	}

	// create WorktreeState
	worktreeState := NewWorktreeState(state)

	// add a new worktree, which should reuse the index of the deleted worktree
	newWorktree := Worktree{Path: "path4", Branch: "branch4", Deleted: false}
	worktreeState.Add(newWorktree, true)

	// check the length of the state after addition
	if len(worktreeState.Items()) != 3 {
		t.Errorf("Expected 3 worktrees in state after addition, got %d", len(worktreeState.Items()))
	}

	// check that the new worktree has replaced the deleted one
	if worktreeState.Items()[1].Path != "path4" || worktreeState.Items()[1].Branch != "branch4" {
		t.Errorf("Expected worktree at index 1 to be replaced with new worktree, got path '%s' and branch '%s'", worktreeState.Items()[1].Path, worktreeState.Items()[1].Branch)
	}

}

func TestWorktreeStateClean(t *testing.T) {
	// prepare state
	state := WorktreeState{
		// prepare 5 dummy worktrees, with 2 last marked as deleted
		items: []Worktree{
			Worktree{Path: "path1", Branch: "branch1", Deleted: false},
			Worktree{Path: "path2", Branch: "branch2", Deleted: true}, // should not be deleted
			Worktree{Path: "path3", Branch: "branch3", Deleted: false},
			Worktree{Path: "path4", Branch: "branch4", Deleted: true}, // should be deleted
			Worktree{Path: "path5", Branch: "branch5", Deleted: true}, // should be deleted
		},
	}

	if len(state.Items()) != 5 {
		t.Errorf("Expected 5 worktrees before cleaning, got %d", len(state.Items()))
	}

	// perform cleaning of deleted worktrees
	state.CleanDeleted()

	// check the length of the state after cleaning
	if len(state.Items()) != 3 {
		t.Errorf("Expected 3 worktrees in state after cleaning, got %d", len(state.Items()))
	}

	// check that the remaining worktrees are the correct ones
	expectedPaths := []string{"path1", "path2", "path3"}

	for i, wt := range state.Items() {
		if wt.Path != expectedPaths[i] {
			t.Errorf("Expected worktree path '%s', got '%s'", expectedPaths[i], wt.Path)
		}
	}

	// add a new worktree after cleaning
	state.Add(Worktree{Path: "path6", Branch: "branch6", Deleted: false}, true)

	// shoould be 3 worktrees still
	if len(state.Items()) != 3 {
		t.Errorf("Expected 3 worktrees in state after adding new worktree, got %d", len(state.Items()))
	}

	// should be on index 1
	if state.Items()[1].Path != "path6" {
		t.Errorf("Expected worktree at index 1 to be 'path6', got '%s'", state.Items()[1].Path)
	}

	if state.Items()[1].Branch != "branch6" {
		t.Errorf("Expected worktree at index 1 to have branch 'branch6', got '%s'", state.Items()[1].Branch)
	}
}

