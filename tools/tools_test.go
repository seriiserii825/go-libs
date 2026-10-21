package tools

import (
	"strings"
	"testing"
)

func TestRequire(t *testing.T) {
	if err := Require(Tool{Bin: "sh"}); err != nil {
		t.Errorf("Require(sh) = %v; want nil", err)
	}
	if err := Require(Tool{Bin: "no-such-xyz", Optional: true}); err != nil {
		t.Errorf("optional missing program: err = %v; want nil", err)
	}

	err := Require(
		Tool{Bin: "sh"},
		Tool{Bin: "no-such-a", Install: "pacman -S a"},
		Tool{Bin: "no-such-b", Install: "npm i -g b"},
	)
	if err == nil {
		t.Fatal("missing required programs: want error")
	}
	for _, want := range []string{"no-such-a: pacman -S a", "no-such-b: npm i -g b"} {
		if !strings.Contains(err.Error(), want) {
			t.Errorf("error has no %q:\n%v", want, err)
		}
	}
	if strings.Contains(err.Error(), "sh:") {
		t.Errorf("installed program in error:\n%v", err)
	}
}
