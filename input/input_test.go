package input

import (
	"bufio"
	"errors"
	"io"
	"strings"
	"testing"
)

func feed(t *testing.T, text string) {
	t.Helper()
	reader = bufio.NewReader(strings.NewReader(text))
	out = io.Discard
}

func TestInt(t *testing.T) {
	feed(t, "abc\n 1.5\n  42 \n")
	if n, err := Int(""); n != 42 || err != nil {
		t.Errorf("Int = %d, %v; want 42, nil", n, err)
	}
}

func TestPositiveInt(t *testing.T) {
	feed(t, "-3\n0\n7\n")
	if n, err := PositiveInt(""); n != 7 || err != nil {
		t.Errorf("PositiveInt = %d, %v; want 7, nil", n, err)
	}
}

func TestFloat(t *testing.T) {
	feed(t, "x\n1.5\n")
	if f, err := Float(""); f != 1.5 || err != nil {
		t.Errorf("Float = %v, %v; want 1.5, nil", f, err)
	}
}

func TestString(t *testing.T) {
	feed(t, "  \n hello \n")
	if s, err := String("", false); s != "hello" || err != nil {
		t.Errorf("String = %q, %v; want hello, nil", s, err)
	}
	feed(t, "\n")
	if s, err := String("", true); s != "" || err != nil {
		t.Errorf("String(allowEmpty) = %q, %v; want empty, nil", s, err)
	}
}

func TestBool(t *testing.T) {
	feed(t, "maybe\nYES\nn\n")
	if b, err := Bool(""); !b || err != nil {
		t.Errorf("Bool = %v, %v; want true, nil", b, err)
	}
	if b, err := Bool(""); b || err != nil {
		t.Errorf("Bool = %v, %v; want false, nil", b, err)
	}
}

func TestConfirm(t *testing.T) {
	tests := map[string]bool{"\n": true, "y\n": true, "Y\n": true, "n\n": false, "No\n": false}
	for in, want := range tests {
		feed(t, in)
		if b, err := Confirm("ok?"); b != want || err != nil {
			t.Errorf("Confirm(%q) = %v, %v; want %v, nil", in, b, err, want)
		}
	}
}

func TestEOF(t *testing.T) {
	feed(t, "")
	if _, err := Int(""); !errors.Is(err, io.EOF) {
		t.Errorf("Int on empty input: err = %v; want io.EOF", err)
	}
	feed(t, "5")
	if n, err := Int(""); n != 5 || err != nil {
		t.Errorf("Int without trailing newline = %d, %v; want 5, nil", n, err)
	}
}
