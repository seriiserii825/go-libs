package fzf

import (
	"regexp"
	"strings"
)

const minChunk = 3

var ansiRe = regexp.MustCompile(`\x1b\[[0-9;]*m`)

// Matches сообщает, подходит ли строка под запрос: каждое слово запроса
// должно быть подстрокой строки или цепочкой подстрок по порядку, где каждый
// кусок, кроме последнего, не короче minChunk символов.
// "resdow" находит "Restore from Downloads" (res + dow), "rfd" — нет.
func Matches(query, line string) bool {
	text := strings.ToLower(ansiRe.ReplaceAllString(line, ""))
	for _, term := range strings.Fields(strings.ToLower(query)) {
		if !chain([]rune(term), text) {
			return false
		}
	}
	return true
}

// term режется по рунам, чтобы не разрезать кириллицу посреди символа.
func chain(term []rune, text string) bool {
	if strings.Contains(text, string(term)) {
		return true
	}
	for n := len(term) - 1; n >= minChunk; n-- {
		head := string(term[:n])
		i := strings.Index(text, head)
		if i != -1 && chain(term[n:], text[i+len(head):]) {
			return true
		}
	}
	return false
}
