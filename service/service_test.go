package service

import (
	"testing"
)

func TestPathNormalization(t *testing.T) {
	service := NewService(nil, nil, nil)

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
		result, err := service.NormalizeAndValidatePath(test.input)
		if err != nil {
			t.Errorf("Unexpected error for input '%s': %v", test.input, err)
			continue
		}
		if result != test.expected {
			t.Errorf("For input '%s', expected '%s', but got '%s'", test.input, test.expected, result)
		}
	}
}
