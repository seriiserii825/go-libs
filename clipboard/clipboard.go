package clipboard

import (
	"errors"
	"os"
	"os/exec"
	"strings"

	"github.com/seriiserii825/go-libs/notify"
)

var errNoDisplay = errors.New("clipboard: neither WAYLAND_DISPLAY nor DISPLAY is set")

// Read возвращает текст из буфера обмена без пробелов по краям.
func Read() (string, error) {
	var cmd *exec.Cmd
	switch {
	case os.Getenv("WAYLAND_DISPLAY") != "":
		cmd = exec.Command("wl-paste", "--no-newline")
	case os.Getenv("DISPLAY") != "":
		cmd = exec.Command("xclip", "-o", "-selection", "clipboard")
	default:
		return "", errNoDisplay
	}
	cmd.Stderr = os.Stderr
	out, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(out)), nil
}

// Write копирует текст в буфер обмена и показывает его в уведомлении.
func Write(text string) error {
	text = strings.TrimSpace(text)
	if err := write(text); err != nil {
		return err
	}
	notify.Send("Buffer", text)
	return nil
}

// WriteSecret — как Write, но текст не показывается в уведомлении
// (для паролей и токенов).
func WriteSecret(text string) error {
	if err := write(strings.TrimSpace(text)); err != nil {
		return err
	}
	notify.Send("Buffer", "Copied to clipboard (hidden)")
	return nil
}

func write(text string) error {
	var cmd *exec.Cmd
	switch {
	case os.Getenv("WAYLAND_DISPLAY") != "":
		cmd = exec.Command("wl-copy")
	case os.Getenv("DISPLAY") != "":
		cmd = exec.Command("xclip", "-selection", "clipboard")
	default:
		return errNoDisplay
	}
	// wl-copy и xclip остаются в фоне, чтобы отдавать буфер. Поэтому stdout
	// не перехватываем (Output() ждал бы, пока фоновый процесс его закроет),
	// а stderr отдаём напрямую.
	cmd.Stdin = strings.NewReader(text)
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
