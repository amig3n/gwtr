package worktree

// NOTE this struct represent a single worktree with all metadata saved in statefile
// NOTE constructur is senseless here; can be formed in-flight
type Worktree struct {
	Path string
	Branch string
	Deleted bool
}

func (wt *worktree) Delete() {
	wt.Deleted = true
}
