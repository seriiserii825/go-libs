package cli

import (
	"fmt"
	"io"
	"os"
	"slices"
	"strings"

	"github.com/seriiserii825/go-libs/color"
	"github.com/seriiserii825/go-libs/menu"
)

var out io.Writer = os.Stdout

// Command — подкоманда: go-scripts <Name>.
type Command struct {
	Name  string // "sort-files"
	Short string // одна строка: в списке команд и в меню
	Long  string // подробная инструкция для <Name> --help
	Color string // цвет пункта в меню, "" — без цвета
	Run   func() error
}

// App — программа с подкомандами:
//
//	app                    — fzf-меню со всеми командами
//	app --help             — список команд
//	app <command>          — запустить команду
//	app <command> --help   — инструкция к команде
type App struct {
	Name     string
	Commands []Command
}

// Run разбирает аргументы (без имени программы, то есть os.Args[1:]).
func (a App) Run(args []string) error {
	if len(args) == 0 {
		return a.menu()
	}
	if isHelp(args[0]) {
		a.printHelp()
		return nil
	}

	cmd, ok := a.find(args[0])
	if !ok {
		return fmt.Errorf("unknown command %q, see '%s --help'", args[0], a.Name)
	}
	if slices.ContainsFunc(args[1:], isHelp) {
		a.printCommandHelp(cmd)
		return nil
	}
	if len(args) > 1 {
		return fmt.Errorf("%s: unexpected arguments %q, see '%s %s --help'", cmd.Name, args[1:], a.Name, cmd.Name)
	}
	return cmd.Run()
}

func (a App) menu() error {
	labels := make([]string, len(a.Commands))
	for i, cmd := range a.Commands {
		labels[i] = cmd.Short
		if cmd.Color != "" {
			labels[i] = color.Wrap(cmd.Color, cmd.Short)
		}
	}
	n, ok, err := menu.SelectFzfMenu(labels, a.Name, color.Wrap(color.Red, "Exit"))
	if err != nil || !ok || n == 0 {
		return err
	}
	return a.Commands[n-1].Run()
}

func (a App) printHelp() {
	fmt.Fprintf(out, "Usage (go-libs v0.9.2):\n")
	fmt.Fprintf(out, "  %s                    interactive menu\n", a.Name)
	fmt.Fprintf(out, "  %s <command>          run command\n", a.Name)
	fmt.Fprintf(out, "  %s <command> --help   command help\n", a.Name)
	fmt.Fprintf(out, "\nCommands:\n")

	width := 0
	for _, cmd := range a.Commands {
		width = max(width, len(cmd.Name))
	}
	for _, cmd := range a.Commands {
		fmt.Fprintf(out, "  %-*s  %s\n", width, cmd.Name, cmd.Short)
	}
}

func (a App) printCommandHelp(cmd Command) {
	fmt.Fprintf(out, "%s %s — %s\n", a.Name, cmd.Name, cmd.Short)
	if long := strings.TrimSpace(cmd.Long); long != "" {
		fmt.Fprintf(out, "\n%s\n", long)
	}
	fmt.Fprintf(out, "\nUsage:\n  %s %s\n", a.Name, cmd.Name)
}

func (a App) find(name string) (Command, bool) {
	for _, cmd := range a.Commands {
		if cmd.Name == name {
			return cmd, true
		}
	}
	return Command{}, false
}

func isHelp(arg string) bool {
	return arg == "-h" || arg == "--help" || arg == "help"
}
