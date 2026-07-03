# Сборка Go-версии lxprofiler

## Требования

- Go 1.21+
- Git

## Установка зависимостей

```bash
go mod tidy
```

## Сборка

```bash
go build -o lxprofile ./cmd/lxprofile
```

## Запуск

```bash
# Статический режим (для пайпов)
./lxprofile --static

# Интерактивный режим (в терминале)
./lxprofile
```

## Кросс-компиляция

```bash
# Linux amd64
GOOS=linux GOARCH=amd64 go build -o lxprofile-linux-amd64 ./cmd/lxprofile

# Linux arm64
GOOS=linux GOARCH=arm64 go build -o lxprofile-linux-arm64 ./cmd/lxprofile

# macOS amd64
GOOS=darwin GOARCH=amd64 go build -o lxprofile-darwin-amd64 ./cmd/lxprofile

# macOS arm64 (Apple Silicon)
GOOS=darwin GOARCH=arm64 go build -o lxprofile-darwin-arm64 ./cmd/lxprofile
```

## Структура проекта

```
cmd/lxprofile/main.go      # CLI entry point
internal/
  data/data.go              # Архетипы и описания
  detect/
    helpers.go              # Вспомогательные функции (Add, Has)
    detect.go               # Анализ системы
  tui/
    model.go                # Bubble Tea модель
    list.go                 # Режим «Список»
    compass.go              # Режим «Компас»
    stats.go                # Режим «Статистика»
go.mod
go.sum
```

## Управление (интерактивный режим)

| Клавиша | Действие |
|---------|----------|
| ↑, k | листать вверх |
| ↓, j | листать вниз |
| →, l, m | следующий режим |
| ←, h | предыдущий режим |
| g | к первому архетипу |
| G | к последнему архетипу |
| q | выход |

## Режимы

1. **Список** — архетипы с прогресс-барами и деталями
2. **Компас** — линуксоидный компас (новаторы ↔ традиции, контроль ↔ удобство)
3. **Статистика** — топ команд из истории

## Лицензия

GPLv2
