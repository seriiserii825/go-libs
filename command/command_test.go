package command

import (
	"strings"
	"testing"
)

func TestOutput(t *testing.T) {
	if out, err := Output("echo", "  hi  "); out != "hi" || err != nil {
		t.Errorf("Output(echo) = %q, %v; want hi, nil", out, err)
	}
	_, err := Output("sh", "-c", "echo oops >&2; exit 3")
	if err == nil || !strings.Contains(err.Error(), "oops") {
		t.Errorf("Output error = %v; want stderr text in it", err)
	}
	if code := ExitCode(err); code != 3 {
		t.Errorf("ExitCode = %d; want 3", code)
	}
}

func TestExitCode(t *testing.T) {
	if code := ExitCode(nil); code != 0 {
		t.Errorf("ExitCode(nil) = %d; want 0", code)
	}
	_, err := Output("no-such-program-xyz")
	if code := ExitCode(err); code != -1 {
		t.Errorf("ExitCode(not found) = %d; want -1", code)
	}
}

func TestExists(t *testing.T) {
	if !Exists("sh") {
		t.Error("Exists(sh) = false; want true")
	}
	if Exists("no-such-program-xyz") {
		t.Error("Exists(no-such-program-xyz) = true; want false")
	}
}

func TestQuote(t *testing.T) {
	for _, s := range []string{"plain", "with space", "it's", `a"b$c`, ""} {
		out, err := Output("sh", "-c", "printf %s "+Quote(s))
		if out != strings.TrimSpace(s) || err != nil {
			t.Errorf("Quote(%q) through sh = %q, %v", s, out, err)
		}
	}
}
