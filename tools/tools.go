package tools

import (
	"fmt"
	"os"
	"strings"

	"github.com/seriiserii825/go-libs/command"
	"github.com/seriiserii825/go-libs/msg"
)

// Tool — внешняя программа, которую вызывает скрипт.
type Tool struct {
	Bin      string
	Install  string // команда установки для сообщения об ошибке
	Optional bool   // без неё скрипт работает, только хуже
	Without  string // что будет без Optional-программы
}

// Программы, которые вызывают пакеты go-libs.
var (
	Fzf       = Tool{Bin: "fzf", Install: "sudo pacman -S fzf"}
	FzfFilter = Tool{
		Bin:      "fzf-word-filter",
		Install:  "go install github.com/seriiserii825/go-libs/cmd/fzf-word-filter@latest (or copy the binary to ~/.local/bin)",
		Optional: true,
		Without:  "fzf uses its own search",
	}
	Notify = Tool{
		Bin:      "notify-send",
		Install:  "sudo pacman -S libnotify",
		Optional: true,
		Without:  "no desktop notifications",
	}
)

// Clipboard — wl-copy в Wayland, xclip в X11 (как в пакете clipboard).
func Clipboard() Tool {
	if os.Getenv("WAYLAND_DISPLAY") != "" {
		return Tool{Bin: "wl-copy", Install: "sudo pacman -S wl-clipboard"}
	}
	return Tool{Bin: "xclip", Install: "sudo pacman -S xclip"}
}

// Require проверяет, что программы установлены. Нет обязательных — ошибка
// со списком команд установки; нет необязательных — только предупреждение.
func Require(tools ...Tool) error {
	var missing []string
	for _, t := range tools {
		if command.Exists(t.Bin) {
			continue
		}
		if t.Optional {
			msg.Warning(fmt.Sprintf("%s is not installed (%s). Install: %s", t.Bin, t.Without, t.Install))
			continue
		}
		missing = append(missing, fmt.Sprintf("  %s: %s", t.Bin, t.Install))
	}
	if len(missing) > 0 {
		return fmt.Errorf("required programs are not installed in your system:\n%s", strings.Join(missing, "\n"))
	}
	return nil
}
