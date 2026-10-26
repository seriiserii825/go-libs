package menu

import (
	"errors"
	"fmt"
	"slices"
	"strconv"
	"strings"

	"github.com/seriiserii825/go-libs/color"
	"github.com/seriiserii825/go-libs/fzf"
	"github.com/seriiserii825/go-libs/input"
	"github.com/seriiserii825/go-libs/msg"
)

// ChooseOption спрашивает номер от 0 до count-1, пока не введут правильный.
// io.EOF — если ввод закончился (Ctrl-D).
func ChooseOption(count int) (int, error) {
	if count <= 0 {
		return 0, errors.New("menu: no options to choose from")
	}
	for {
		n, err := input.Int(fmt.Sprintf("Please enter a number between 0 and %d: ", count-1))
		if err != nil {
			return 0, err
		}
		if n >= 0 && n < count {
			return n, nil
		}
		fmt.Println(color.Wrap(color.Red, fmt.Sprintf("Invalid input. Enter a number between 0 and %d.", count-1)))
	}
}

// SelectFzf — выбор одного варианта, возвращает его индекс в options.
// Заголовком (--header) служит текущий путь по меню (BreadcrumbPath), если
// он не пуст. false, если ничего не выбрано.
func SelectFzf(options []string) (int, bool, error) {
	args := []string{"--height", "50%", "--reverse", "--no-mouse"}
	args = append(args, headerArgs()...)
	out, code, err := fzf.Run(options, args)
	if err != nil || code != 0 {
		return 0, false, err
	}
	s := strings.TrimRight(out, "\n")
	if s == "" {
		return 0, false, nil
	}
	i := slices.Index(options, s)
	return i, i != -1, nil
}

// SelectFzfMulti — выбор нескольких вариантов, возвращает их индексы в
// options. Заголовком (--header) служит текущий путь по меню
// (BreadcrumbPath), если он не пуст.
func SelectFzfMulti(options []string) ([]int, error) {
	msg.Info("<Tab> select, <Ctrl-A> select all, <Ctrl-R> invert selection, <Enter> confirm.")
	args := []string{
		"--multi",
		"--height", "50%",
		"--reverse",
		"--no-mouse",
		"--bind", "ctrl-a:select-all,ctrl-r:toggle-all",
	}
	args = append(args, headerArgs()...)
	out, code, err := fzf.Run(options, args)
	if err != nil || code != 0 {
		return nil, err
	}
	out = strings.TrimRight(out, "\n")
	if out == "" {
		return nil, nil
	}
	var indexes []int
	for _, s := range strings.Split(out, "\n") {
		if i := slices.Index(options, s); i != -1 {
			indexes = append(indexes, i)
		}
	}
	return indexes, nil
}

// headerArgs — "--header <путь>", если текущий путь по меню не пуст, иначе
// пусто.
func headerArgs() []string {
	if path := BreadcrumbPath(); path != "" {
		return []string{"--header", path}
	}
	return nil
}

// breadcrumbs — стек названий вложенных меню, в которых сейчас находится
// пользователь (см. PushBreadcrumb).
var breadcrumbs []string

// PushBreadcrumb добавляет name в конец текущего пути по меню. Вызвать
// возвращённую функцию (обычно через defer), когда меню name закрывается,
// чтобы убрать его из пути.
func PushBreadcrumb(name string) func() {
	breadcrumbs = append(breadcrumbs, name)
	return func() {
		if len(breadcrumbs) > 0 {
			breadcrumbs = breadcrumbs[:len(breadcrumbs)-1]
		}
	}
}

// BreadcrumbPath — текущий путь по вложенным меню через " > "
// ("Main Menu > Page checks > SEO"), "" если PushBreadcrumb ещё не вызывали.
func BreadcrumbPath() string {
	return strings.Join(breadcrumbs, " > ")
}

// SelectFzfMenu — цветное fzf-меню с номерами: "01) First", "02) Second",
// "00) Exit".
// labels могут быть цветными (color.Wrap). Можно фильтровать по тексту или
// набрать номер и нажать Enter.
// Заголовком меню (--header) служит полный путь по меню (BreadcrumbPath),
// если он не пуст, иначе — title.
// Возвращает 1..len(labels) для пункта, 0 для exitLabel (exitLabel "" — пункта
// нет), false на Esc/Ctrl-C.
func SelectFzfMenu(labels []string, title, exitLabel string) (int, bool, error) {
	lines := make([]string, 0, len(labels)+1)
	for i, label := range labels {
		lines = append(lines, index(i+1)+" "+label)
	}
	if exitLabel != "" {
		lines = append(lines, index(0)+" "+exitLabel)
	}

	args := []string{
		"--ansi",
		"--reverse",
		"--no-mouse",
		"--no-sort",
		"--height", "50%",
		"--print-query",
	}
	if h := headerArgs(); h != nil {
		args = append(args, h...)
	} else if title != "" {
		args = append(args, "--header", title)
	}

	out, code, err := fzf.Run(lines, args)
	if err != nil {
		return 0, false, err
	}
	if code != 0 && code != 1 {
		return 0, false, nil
	}
	n, ok := parseMenuOutput(out, len(labels), exitLabel != "")
	return n, ok, nil
}

func index(n int) string {
	return color.Wrap(color.Magenta, fmt.Sprintf("%02d)", n))
}

// parseMenuOutput разбирает вывод fzf --print-query: первая строка — запрос,
// вторая — выбранная строка. Номер в запросе важнее выбранной строки.
func parseMenuOutput(out string, maxIndex int, hasExit bool) (int, bool) {
	query, selected, _ := strings.Cut(out, "\n")
	selected, _, _ = strings.Cut(selected, "\n")

	if n, err := strconv.Atoi(strings.TrimSpace(query)); err == nil && n >= 0 && n <= maxIndex {
		if n != 0 || hasExit {
			return n, true
		}
	}

	num, _, found := strings.Cut(color.Strip(selected), ")")
	if !found {
		return 0, false
	}
	n, err := strconv.Atoi(strings.TrimSpace(num))
	if err != nil {
		return 0, false
	}
	return n, true
}
