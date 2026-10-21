package fzf

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"slices"
	"strings"

	"github.com/seriiserii825/go-libs/command"
)

// Run запускает fzf со строками lines и флагами args.
// Если в PATH есть fzf-word-filter, фильтрация идёт через него, а не через
// собственный поиск fzf. Возвращает вывод fzf и код выхода:
// 0 — выбор сделан, 1 — нет совпадений, 130 — Esc/Ctrl-C.
func Run(lines, args []string) (string, int, error) {
	input := strings.Join(lines, "\n")

	if filter, err := exec.LookPath("fzf-word-filter"); err == nil {
		f, err := os.CreateTemp("", "*.fzf")
		if err != nil {
			return "", 0, err
		}
		defer os.Remove(f.Name())
		_, err = f.WriteString(input)
		f.Close()
		if err != nil {
			return "", 0, err
		}
		reload := fmt.Sprintf("%s {q} %s || true", command.Quote(filter), command.Quote(f.Name()))
		args = slices.Concat(args, []string{"--disabled", "--bind", "change:reload(" + reload + ")"})
	}

	cmd := exec.Command("fzf", args...)
	cmd.Stdin = strings.NewReader(input)
	cmd.Stderr = os.Stderr
	out, err := cmd.Output()
	if err != nil {
		var exitErr *exec.ExitError
		if errors.As(err, &exitErr) {
			return string(out), exitErr.ExitCode(), nil
		}
		return "", 0, err
	}
	return string(out), 0, nil
}

// One — выбор одного варианта. false, если ничего не выбрано.
func One(options []string) (string, bool, error) {
	out, code, err := Run(options, []string{"--height", "50%", "--reverse", "--no-mouse"})
	if err != nil || code != 0 {
		return "", false, err
	}
	selected := strings.TrimRight(out, "\n")
	return selected, selected != "", nil
}

// Multi — выбор нескольких вариантов: <Tab> отметить, ctrl-a выбрать все,
// ctrl-r инвертировать.
func Multi(options []string) ([]string, error) {
	out, code, err := Run(options, []string{
		"--multi",
		"--height", "50%",
		"--reverse",
		"--no-mouse",
		"--bind", "ctrl-a:select-all,ctrl-r:toggle-all",
	})
	if err != nil || code != 0 {
		return nil, err
	}
	out = strings.TrimRight(out, "\n")
	if out == "" {
		return nil, nil
	}
	return strings.Split(out, "\n"), nil
}
