package cli

import (
	"bytes"
	"strings"
	"testing"
)

func testApp(ran *string) App {
	return App{
		Name: "app",
		Commands: []Command{
			{Name: "sort-files", Short: "Sort files", Long: "Copies files.", Run: func() error {
				*ran = "sort-files"
				return nil
			}},
			{Name: "aspect-ratio", Short: "Aspect ratio", Run: func() error {
				*ran = "aspect-ratio"
				return nil
			}},
		},
	}
}

func run(t *testing.T, args ...string) (string, string, error) {
	t.Helper()
	var buf bytes.Buffer
	out = &buf
	var ran string
	err := testApp(&ran).Run(args)
	return buf.String(), ran, err
}

func TestHelp(t *testing.T) {
	for _, arg := range []string{"--help", "-h", "help"} {
		text, ran, err := run(t, arg)
		if err != nil || ran != "" {
			t.Fatalf("%s: ran=%q err=%v", arg, ran, err)
		}
		if !strings.Contains(text, "sort-files    Sort files") || !strings.Contains(text, "aspect-ratio  Aspect ratio") {
			t.Errorf("%s: help without aligned commands:\n%s", arg, text)
		}
	}
}

func TestCommandHelp(t *testing.T) {
	text, ran, err := run(t, "sort-files", "--help")
	if err != nil || ran != "" {
		t.Fatalf("ran=%q err=%v", ran, err)
	}
	for _, want := range []string{"app sort-files — Sort files", "Copies files.", "Usage:\n  app sort-files"} {
		if !strings.Contains(text, want) {
			t.Errorf("command help has no %q:\n%s", want, text)
		}
	}
}

func TestRunCommand(t *testing.T) {
	if _, ran, err := run(t, "aspect-ratio"); err != nil || ran != "aspect-ratio" {
		t.Errorf("ran=%q err=%v; want aspect-ratio, nil", ran, err)
	}
}

func TestErrors(t *testing.T) {
	if _, ran, err := run(t, "nope"); err == nil || ran != "" {
		t.Errorf("unknown command: ran=%q err=%v; want error", ran, err)
	}
	if _, ran, err := run(t, "sort-files", "extra"); err == nil || ran != "" {
		t.Errorf("extra args: ran=%q err=%v; want error", ran, err)
	}
}
