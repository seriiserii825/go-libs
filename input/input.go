package input

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"

	"github.com/seriiserii825/go-libs/color"
)

// Один общий reader на весь stdin: если создавать новый bufio на каждый вызов,
// он может забрать в буфер следующие строки, и они потеряются.
var (
	reader           = bufio.NewReader(os.Stdin)
	out    io.Writer = os.Stdout
)

// Line печатает prompt и возвращает введённую строку без пробелов по краям.
// io.EOF — если ввод закончился (Ctrl-D).
func Line(prompt string) (string, error) {
	fmt.Fprint(out, prompt)
	line, err := reader.ReadString('\n')
	if errors.Is(err, io.EOF) && line != "" {
		err = nil // последняя строка без \n
	}
	if err != nil {
		fmt.Fprintln(out) // после Ctrl-D курсор остаётся на строке с prompt
		return "", err
	}
	return strings.TrimSpace(line), nil
}

// Int спрашивает целое число, пока не введут правильное.
func Int(prompt string) (int, error) {
	for {
		s, err := Line(prompt)
		if err != nil {
			return 0, err
		}
		n, err := strconv.Atoi(s)
		if err == nil {
			return n, nil
		}
		invalid("Please enter a valid integer.")
	}
}

// PositiveInt спрашивает целое число больше 0.
func PositiveInt(prompt string) (int, error) {
	for {
		n, err := Int(prompt)
		if err != nil {
			return 0, err
		}
		if n > 0 {
			return n, nil
		}
		invalid("Please enter a number greater than 0.")
	}
}

// Float спрашивает число, можно дробное (1.5).
func Float(prompt string) (float64, error) {
	for {
		s, err := Line(prompt)
		if err != nil {
			return 0, err
		}
		f, err := strconv.ParseFloat(s, 64)
		if err == nil {
			return f, nil
		}
		invalid("Please enter a valid number.")
	}
}

// String спрашивает текст; allowEmpty=false — пока не введут непустой.
func String(prompt string, allowEmpty bool) (string, error) {
	for {
		s, err := Line(prompt)
		if err != nil {
			return "", err
		}
		if s != "" || allowEmpty {
			return s, nil
		}
		invalid("Input cannot be empty.")
	}
}

// Bool спрашивает y/yes или n/no, пока не введут одно из них.
func Bool(prompt string) (bool, error) {
	for {
		s, err := Line(prompt)
		if err != nil {
			return false, err
		}
		switch strings.ToLower(s) {
		case "y", "yes":
			return true, nil
		case "n", "no":
			return false, nil
		}
		invalid("Please enter 'y' or 'n'.")
	}
}

// Confirm задаёт вопрос да/нет, по умолчанию да: Enter или y — true,
// n или no — false.
func Confirm(question string) (bool, error) {
	s, err := Line(question + " (Y/n, Y - press Enter): ")
	if err != nil {
		return false, err
	}
	switch strings.ToLower(s) {
	case "n", "no":
		return false, nil
	}
	return true, nil
}

func invalid(text string) {
	fmt.Fprintln(out, color.Wrap(color.Red, "Invalid input. "+text))
}
