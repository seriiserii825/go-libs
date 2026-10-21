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
