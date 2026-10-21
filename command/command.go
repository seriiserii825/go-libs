package command

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

// Run запускает команду с выводом прямо в терминал (как subprocess.run).
// Аргументы передаются как есть, без shell — кавычки не нужны.
func Run(name string, args ...string) error {
	cmd := exec.Command(name, args...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("%s: %w", line(name, args), err)
	}
	return nil
}

// Output запускает команду и возвращает её stdout без пробелов по краям.
// При ошибке в тексте ошибки будет stderr команды.
func Output(name string, args ...string) (string, error) {
	var stderr bytes.Buffer
	cmd := exec.Command(name, args...)
	cmd.Stderr = &stderr
	out, err := cmd.Output()
	if err != nil {
		if msg := strings.TrimSpace(stderr.String()); msg != "" {
			return "", fmt.Errorf("%s: %w\n%s", line(name, args), err, msg)
		}
		return "", fmt.Errorf("%s: %w", line(name, args), err)
	}
	return strings.TrimSpace(string(out)), nil
}

// Shell запускает строку через sh -c с выводом в терминал (как os.system):
// работают пайпы, &&, *, ~. Подставляемые значения оборачивай в Quote.
func Shell(script string) error {
	return Run("sh", "-c", script)
}

// Exists сообщает, есть ли программа в PATH (как which).
func Exists(name string) bool {
	_, err := exec.LookPath(name)
	return err == nil
}

// ExitCode возвращает код выхода из ошибки Run/Output/Shell, -1 — если
// команда не запустилась, 0 — если ошибки нет.
func ExitCode(err error) int {
	if err == nil {
		return 0
	}
	var exitErr *exec.ExitError
	if errors.As(err, &exitErr) {
		return exitErr.ExitCode()
	}
	return -1
}

// Quote экранирует строку для shell: Shell("rm -rf " + Quote(dir)).
func Quote(s string) string {
	return "'" + strings.ReplaceAll(s, "'", `'\''`) + "'"
}

func line(name string, args []string) string {
	return strings.Join(append([]string{name}, args...), " ")
}
