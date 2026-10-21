# go-libs

Общие пакеты для Go-скриптов. Подключаются в проекты как Go-модуль по версии
(git-тег), без git submodule.

```
module github.com/seriiserii825/go-libs
```

## Пакеты

| Пакет | Что делает |
|---|---|
| `color` | ANSI-цвета и `Wrap` |
| `msg` | цветные сообщения `Info` / `Success` / `Error` / `Warning` |
| `fzf` | выбор через fzf: `One`, `Multi`, `Run`, фильтр `Matches` |
| `table` | таблица с рамкой `Show` и сетка по ширине терминала `Grid` |
| `menu` | меню: fzf-меню с номерами, выбор индекса, ввод номера |
| `input` | ввод с проверкой: `Int`, `Float`, `String`, `Bool`, `Confirm` |
| `clipboard` | буфер обмена `Read` / `Write` / `WriteSecret` (Wayland и X11) |
| `notify` | уведомление на рабочем столе `Send` |
| `command` | запуск программ: `Run`, `Output`, `Shell`, `Exists`, `Quote` |
| `cli` | подкоманды с `--help` и fzf-меню без аргументов |
| `cmd/fzf-word-filter` | бинарник-фильтр, который fzf вызывает при вводе |

## Установка в проект

```
go get github.com/seriiserii825/go-libs@latest
go install github.com/seriiserii825/go-libs/cmd/fzf-word-filter@latest
```

- `go get` добавит `require` в `go.mod` и создаст `go.sum` — оба коммитить.
- `go install` кладёт `fzf-word-filter` в `~/go/bin` — эта папка должна быть
  в `PATH`. Ставится один раз на машину. Без него `fzf` работает со своим
  обычным поиском.

## Использование

### color

```go
import "github.com/seriiserii825/go-libs/color"

fmt.Println(color.Wrap(color.Green, "ok"))
fmt.Println(color.Red + "error" + color.Reset)
```

Цвета: `Red`, `Green`, `Yellow`, `Blue`, `Magenta`, `Cyan`, `Reset`.
`color.Strip(s)` убирает ANSI-цвета из строки.

### msg

```go
import "github.com/seriiserii825/go-libs/msg"

msg.Info("start")
msg.Success("done")
msg.Warning("skip")
msg.Error("failed")
```

### fzf

```go
import "github.com/seriiserii825/go-libs/fzf"

// один вариант
s, ok, err := fzf.One([]string{"Restore from Downloads", "Backup"})
if err != nil {
	log.Fatal(err)
}
if !ok {
	return // Esc / ничего не выбрано
}

// несколько: <Tab> отметить, ctrl-a выбрать все, ctrl-r инвертировать
items, err := fzf.Multi([]string{"a", "b", "c"})

// свои флаги fzf
out, code, err := fzf.Run(lines, []string{"--reverse", "--print-query"})
```

Коды выхода `Run`: `0` — выбор сделан, `1` — нет совпадений, `130` — Esc /
Ctrl-C. Это не ошибки: `err` только если fzf не запустился.

#### Поиск

Каждое слово запроса должно быть подстрокой строки или цепочкой подстрок по
порядку, где каждый кусок, кроме последнего, не короче 3 символов. Регистр и
ANSI-цвета игнорируются, кириллица поддерживается.

- `resdow` → `Restore from Downloads` (res + dow)
- `down rest` → `Restore from Downloads` (слова в любом порядке)
- `rfd` → не найдёт (побуквенного fuzzy нет)

### table

```go
import "github.com/seriiserii825/go-libs/table"

table.Show("Сайты", []string{"#", "Name", "Status"}, [][]string{
	{"0", "Restore from Downloads", "ok"},
	{"1", "Бэкап", "fail"},
}, map[int]string{1: color.Red}) // цвет строки по индексу, или nil

table.Grid(items, color.Blue) // "" — без цвета
```

```
Сайты
╭───┬────────────────────────┬────────╮
│ # │ Name                   │ Status │
├───┼────────────────────────┼────────┤
│ 0 │ Restore from Downloads │ ok     │
│ 1 │ Бэкап                  │ fail   │
╰───┴────────────────────────┴────────╯
пункт-00  пункт-04  пункт-08  пункт-12
пункт-01  пункт-05  пункт-09  пункт-13
```

`Grid` раскладывает элементы сверху вниз по колонкам; число колонок зависит от
ширины терминала (если вывод не в терминал — 80).

### menu

```go
import "github.com/seriiserii825/go-libs/menu"

// fzf-меню с номерами: 01) Deploy, 02) Backup, 00) Exit
n, ok, err := menu.SelectFzfMenu(
	[]string{color.Wrap(color.Green, "Deploy"), "Backup"},
	"Что делаем?",
	color.Wrap(color.Red, "Exit"), // "" — без пункта Exit
)
switch {
case err != nil:
	log.Fatal(err)
case !ok: // Esc / Ctrl-C
case n == 0: // Exit
default: // n = 1..len(labels)
}

// индекс выбранного варианта в options
i, ok, err := menu.SelectFzf(options)

// индексы нескольких вариантов
idx, err := menu.SelectFzfMulti(options)

// после table.Show: спросить номер 0..count-1 в консоли
i, err := menu.ChooseOption(len(rows)) // io.EOF на Ctrl-D
```

В `SelectFzfMenu` можно фильтровать по тексту или набрать номер (`2`, `02`) и
нажать Enter — номер в запросе важнее подсвеченной строки.

### input

```go
import "github.com/seriiserii825/go-libs/input"

n, err := input.Int("Width: ")             // спрашивает, пока не введут число
n, err := input.PositiveInt("Count: ")     // > 0
f, err := input.Float("Ratio: ")           // 1.5
s, err := input.String("Name: ", false)    // true — можно пустую
b, err := input.Bool("Delete? (y/n): ")    // y/yes или n/no
ok, err := input.Confirm("Continue?")      // Enter = да, n = нет
s, err := input.Line("> ")                 // строка как есть, без проверки
```

`err` — только `io.EOF` (Ctrl-D) или ошибка чтения; неверный ввод
переспрашивается сам. Все функции читают через один общий буфер stdin — не
смешивать с собственным `bufio.Scanner(os.Stdin)` в проекте, иначе часть ввода
может потеряться.

### clipboard

```go
import "github.com/seriiserii825/go-libs/clipboard"

s, err := clipboard.Read()
err = clipboard.Write("16:9")      // + уведомление с текстом
err = clipboard.WriteSecret(pass)  // + уведомление без текста
```

Wayland (`WAYLAND_DISPLAY`) — `wl-copy` / `wl-paste`, X11 (`DISPLAY`) —
`xclip`. Пакеты: `wl-clipboard`, `xclip`, `libnotify` (`notify-send`).
Ошибка уведомления игнорируется — главное, что текст скопирован.

### notify

```go
import "github.com/seriiserii825/go-libs/notify"

err := notify.Send("Backup", "Done")
```

### command

```go
import "github.com/seriiserii825/go-libs/command"

// вывод прямо в терминал (git clone с прогрессом, интерактивные программы)
err := command.Run("git", "clone", "--mirror", url, dir)

// перехватить stdout; в ошибке будет stderr команды
branch, err := command.Output("git", "branch", "--show-current")

// через sh -c: пайпы, &&, *, ~ (значения — через Quote)
err = command.Shell("npm ls -g | grep ttf2woff")
err = command.Shell("cd " + command.Quote(dir) + " && npm install")

if !command.Exists("woff2_compress") { ... } // как which
code := command.ExitCode(err)                 // 0, код выхода или -1
```

`Run` и `Output` передают аргументы без shell — кавычки и `Quote` не нужны,
пробелы в путях безопасны. `Shell` — только когда нужны возможности shell.

### cli

```go
import "github.com/seriiserii825/go-libs/cli"

var app = cli.App{
	Name: "go-scripts",
	Commands: []cli.Command{
		{
			Name:  "sort-files",               // go-scripts sort-files
			Short: "Sort files",               // в списке команд и в меню
			Long:  "Copy every file ...",      // go-scripts sort-files --help
			Color: color.Blue,                 // цвет в меню
			Run:   scripts.SortFiles,          // func() error
		},
	},
}

func main() {
	if err := app.Run(os.Args[1:]); err != nil && !errors.Is(err, io.EOF) {
		msg.Error(err.Error())
		os.Exit(1)
	}
}
```

```
go-scripts                     # fzf-меню со всеми командами
go-scripts --help              # список команд (-h, help тоже)
go-scripts sort-files          # запустить команду сразу, без меню
go-scripts sort-files --help   # инструкция к команде (Long)
```

Неизвестная команда или лишние аргументы — ошибка с подсказкой `--help`.

## Разработка

### Локально — go work

В проекте, где нужна правленая версия библиотеки:

```
go work init . ../go-libs
echo go.work >> .gitignore
echo go.work.sum >> .gitignore
```

Пока есть `go.work`, используется **локальная папка** `../go-libs` — правки
видны сразу, теги и версия из `go.mod` игнорируются.

Проверить работу именно по тегу — временно выключить `go.work`:

```
GOWORK=off go run .
```

### Проверки

```
go vet ./...
go test ./...
go install ./cmd/fzf-word-filter   # если менялся фильтр
```

### Новая версия

```
git add . && git commit -m "..."
git push
git tag v0.2.0
git push --tags
```

- patch (`v0.2.1`) — фикс без изменения API.
- minor (`v0.3.0`) — новые функции или пакеты.
- major (`v2.0.0`) — ломающие изменения; путь модуля меняется на
  `github.com/seriiserii825/go-libs/v2`.
- Запушенный тег **не перемещать и не пересоздавать** — Go кэширует checksum
  (`sum.golang.org`), будет `checksum mismatch`. Ошибся — делай следующий тег.

### Обновить версию в проекте

```
GOWORK=off go get github.com/seriiserii825/go-libs@v0.2.0
GOWORK=off go mod tidy
```

Затем закоммитить `go.mod` и `go.sum`. Если менялся фильтр — переустановить:

```
go install github.com/seriiserii825/go-libs/cmd/fzf-word-filter@v0.2.0
```

### Если новый тег не находится

Модули качаются через `proxy.golang.org`, новый тег может появиться с
задержкой. Взять сразу с GitHub:

```
GOPROXY=direct GOWORK=off go get github.com/seriiserii825/go-libs@v0.2.0
```

Навсегда для своих репо (обязательно, если репо приватный):

```
go env -w GOPRIVATE=github.com/seriiserii825/*
```

## Именование пакетов

- Нельзя `select` — ключевое слово Go.
- Не стоит `print` — совпадает со встроенной функцией.
- Не делать `utils` / `common` — пакет называется по тому, что он делает.
