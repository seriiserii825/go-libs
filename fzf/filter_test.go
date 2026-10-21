package fzf

import "testing"

func TestMatches(t *testing.T) {
	tests := []struct {
		query, line string
		want        bool
	}{
		{"resdow", "Restore from Downloads", true},
		{"rfd", "Restore from Downloads", false},
		{"", "anything", true},
		{"down rest", "Restore from Downloads", true},
		{"RES", "restore", true},
		{"res", "\x1b[32mRestore\x1b[0m", true},
		{"02", "\x1b[35m02)\x1b[0m Second", true},
		{"восзаг", "Восстановить из Загрузок", true},
		{"взз", "Восстановить из Загрузок", false},
	}
	for _, tt := range tests {
		if got := Matches(tt.query, tt.line); got != tt.want {
			t.Errorf("Matches(%q, %q) = %v, want %v", tt.query, tt.line, got, tt.want)
		}
	}
}
