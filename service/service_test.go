package service

import (
	"testing"
	"log/slog"
)

// FIXME not working, requires mocking of git due to git repo path relation
func TestPathSanitization(t *testing.T) {
	t.Skip("Skipping test due to dependency on actual git repository paths. TBF.")
	logger := slog.New(slog.NewTextHandler(nil, nil))
	service := NewService(logger, nil, nil)

	tests := []struct {
		input    string
		expected string
	}{
		{"./worktree1", "worktree1"},
		{"../worktree2", "worktree2"},
		{"/absolute/path/worktree3", "absolute/path/worktree3"},
		{"worktree4", "worktree4"},
	}

	for _, test := range tests {
		result, err := service.SanitizePath(test.input)
		if err != nil {
			t.Errorf("Unexpected error for input '%s': %v", test.input, err)
			continue
		}
		if result != test.expected {
			t.Errorf("For input '%s', expected '%s', but got '%s'", test.input, test.expected, result)
		}
	}
}
