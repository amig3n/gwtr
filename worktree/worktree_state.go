package worktree

import (
	"fmt"
)

// NOTE prepared state created from combining statefile and Git info
type WorktreeState struct {
	items []Worktree
}

func NewWorktreeState(items []Worktree) *WorktreeState {
	return &WorktreeState{
		items: items,
	}
}

func (ws *WorktreeState) Items() []Worktree {
	return ws.items
}

func (ws *WorktreeState) Add(item Worktree) {
	for i, wt := range ws.items {
		if wt.Deleted {
			// reuse the free index with early return
			ws.items[i] = item
			return
		}
	}
	// simple appending of item to state if no free index is available
	ws.items = append(ws.items, item)
}

func (ws *WorktreeState) Delete(id int) error {
	if id < 0 || id >= len(ws.items) {
		return fmt.Errorf("worktree delete error: index out of range")
	}
	ws.items[id].Deleted = true
	return nil
}

func (ws *WorktreeState) GetByID(id int) (*Worktree, error) {
	if id < 0 || id >= len(ws.items) {
		return nil, fmt.Errorf("worktree get error: index out of range")
	}

	if ws.items[id].Deleted {
		return nil, fmt.Errorf("worktree get error: worktree at index %d is marked as deleted", id)
	}

	return &ws.items[id], nil
}

func (ws *WorktreeState) GetByString(id string) (*Worktree, error) {
	for _, item := range ws.items {
		if item.Path == id || item.Branch == id {
			return &item, nil
		}
	}
	return nil, fmt.Errorf("worktree get error: no worktree found with path or branch '%s'", id)
}

// clean all deleted worktrees from the end of the state file (returns the number of deleted worktrees)
func (ws *WorktreeState) CleanDeleted() int {
	stateLength := len(ws.items)

	deletedCount := 0
	// iterate from the end of slice till first non-deleted WT found
	for i := stateLength - 1; i >= 0; i-- {
		if ws.items[i].Deleted {
			// drop the last element from the slice
			ws.items = ws.items[:i]
			deletedCount++
		} else {
			break
		}
	}

	return deletedCount
}



