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
	// simple appending of item to state
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



