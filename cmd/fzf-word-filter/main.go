// fzf-word-filter печатает строки файла, подходящие под запрос (см. fzf.Matches).
// Вызывается из fzf через --bind change:reload(...).
//
// Usage: fzf-word-filter <query> <file>
package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/seriiserii825/go-libs/fzf"
)

func main() {
	if len(os.Args) < 3 {
		fmt.Fprintln(os.Stderr, "usage: fzf-word-filter <query> <file>")
		os.Exit(2)
	}
	query, path := os.Args[1], os.Args[2]

	data, err := os.ReadFile(path)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	for _, line := range strings.Split(string(data), "\n") {
		if fzf.Matches(query, line) {
			fmt.Println(line)
		}
	}
}
