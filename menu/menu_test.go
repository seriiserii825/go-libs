package menu

import "testing"

func TestParseMenuOutput(t *testing.T) {
	tests := []struct {
		name    string
		out     string
		hasExit bool
		want    int
		wantOk  bool
	}{
		{"pick by enter", "\n02) Second\n", true, 2, true},
		{"pick by filter", "sec\n02) Second\n", true, 2, true},
		{"number in query", "3\n03) Third\n", true, 3, true},
		{"number wins over selection", "1\n01) First\n", true, 1, true},
		{"exit by number", "0\n00) Exit\n", true, 0, true},
		{"zero without exit", "0\n", false, 0, false},
		{"number too big", "9\n", true, 0, false},
		{"exit by enter", "\n00) Exit\n", true, 0, true},
		{"colored selection", "\n\x1b[35m02)\x1b[0m Second\n", true, 2, true},
		{"no match", "zzz\n", true, 0, false},
		{"empty", "", true, 0, false},
	}
	for _, tt := range tests {
		got, ok := parseMenuOutput(tt.out, 3, tt.hasExit)
		if got != tt.want || ok != tt.wantOk {
			t.Errorf("%s: parseMenuOutput(%q) = %d, %v; want %d, %v", tt.name, tt.out, got, ok, tt.want, tt.wantOk)
		}
	}
}
