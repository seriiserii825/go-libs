package table

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/charmbracelet/lipgloss"
	ltable "github.com/charmbracelet/lipgloss/table"
	"golang.org/x/term"

	"github.com/seriiserii825/go-libs/color"
)

var (
	titleStyle  = lipgloss.NewStyle().Bold(true)
	headerStyle = lipgloss.NewStyle().Bold(true).Padding(0, 1)
	cellStyle   = lipgloss.NewStyle().Padding(0, 1)
)

// Show печатает таблицу с рамкой.
// rowColors — цвет строки по индексу (с нуля), например {0: color.Green};
// можно передать nil.
func Show(title string, columns []string, rows [][]string, rowColors map[int]string) {
	indexedColumns := append([]string{"#"}, columns...)
	colored := make([][]string, len(rows))
	for i, row := range rows {
		indexed := append([]string{strconv.Itoa(i + 1)}, row...)
		c, ok := rowColors[i]
		if !ok {
			colored[i] = indexed
			continue
		}
		colored[i] = make([]string, len(indexed))
		for j, cell := range indexed {
			colored[i][j] = color.Wrap(c, cell)
		}
	}

	t := ltable.New().
		Border(lipgloss.RoundedBorder()).
		Headers(indexedColumns...).
		Rows(colored...).
		StyleFunc(func(row, col int) lipgloss.Style {
			if row == ltable.HeaderRow {
				return headerStyle
			}
			return cellStyle
		})

	if width, _, err := term.GetSize(int(os.Stdout.Fd())); err == nil && width > 0 {
		if lipgloss.Width(t.String()) > width {
			t = t.Width(width - 1)
		}
	}

	if title != "" {
		fmt.Println(titleStyle.Render(title))
	}
	fmt.Println(t)
}

// Grid печатает элементы в несколько колонок по ширине терминала.
// Элементы идут сверху вниз, потом в следующую колонку.
// c — цвет из пакета color, "" — без цвета.
func Grid(items []string, c string) {
	if len(items) == 0 {
		return
	}

	width, _, err := term.GetSize(int(os.Stdout.Fd()))
	if err != nil {
		width = 80
	}

	maxLen := 0
	for _, item := range items {
		maxLen = max(maxLen, lipgloss.Width(item))
	}
	maxLen += 2

	cols := max(1, width/maxLen)
	rows := (len(items) + cols - 1) / cols

	for r := range rows {
		var line strings.Builder
		for col := range cols {
			i := col*rows + r
			if i >= len(items) {
				break
			}
			item := items[i]
			pad := strings.Repeat(" ", maxLen-lipgloss.Width(item))
			if c != "" {
				item = color.Wrap(c, item)
			}
			line.WriteString(item + pad)
		}
		fmt.Println(strings.TrimRight(line.String(), " "))
	}
}
