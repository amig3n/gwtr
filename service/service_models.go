package service

// NOTE data model to aggregate from git provider
type GitWorktree struct {
	Path string
	Branch string
}

// NOTE data model used by the state file
type RawState struct {
	Path    string `json:"path"`
	Deleted bool   `json:"deleted"`
}

// NOTE contract for git provider
type GitProvider interface {
	GetRepoRootPath() (string, error)
	ListWorktrees() ([]GitWorktree, error)
	AddWorktree(branch string, path string) error
	DeleteWorktree(string) error
}

// NOTE contract for state Store
type StateStore interface {
	Load() ([]RawState, error)
	Save([]RawState) error
	Init() error
}
