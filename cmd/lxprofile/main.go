package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/Axawys/lxprofiler/internal/data"
	"github.com/Axawys/lxprofiler/internal/detect"
	"github.com/Axawys/lxprofiler/internal/tui"
)

// Version — версия сборки. По умолчанию для локальной сборки; на релизе
// подставляется линкером: go build -ldflags "-X main.Version=X.Y.Z".
var Version = "5.11.1"

func main() {
	args := os.Args[1:]

	// Слитная форма короткой команды: lxu == lx -u, lxc == lx -c и т.д.
	// (симлинки lxu/lxs/… создаёт ensureShortCommands). Имя вызова → флаг.
	switch filepath.Base(os.Args[0]) {
	case "lxu":
		args = append([]string{"-u"}, args...)
	case "lxs":
		args = append([]string{"-s"}, args...)
	case "lxv":
		args = append([]string{"-v"}, args...)
	case "lxh":
		args = append([]string{"-h"}, args...)
	case "lxc":
		args = append([]string{"-c"}, args...)
	case "lxrm":
		args = append([]string{"--rm"}, args...)
	case "lxa":
		args = append([]string{"-a"}, args...)
	case "lxf":
		args = append([]string{"-f"}, args...)
	case "lxfm":
		args = append([]string{"-fm"}, args...)
	}

	forceStatic := false
	forceAnim := false

	if len(args) > 0 {
		// Флаги принимаются с дефисом (-u/--update), без дефиса (u/update)
		// и слитно с командой lx (lxu) — формы равнозначны.
		switch strings.TrimLeft(args[0], "-") {
		case "v", "version":
			printVersion()
			return
		case "h", "help":
			printHelp()
			return
		case "u", "update":
			ensureShortCommands(true)
			if err := doUpdate(argAt(args, 1)); err != nil {
				fmt.Fprintln(os.Stderr, err)
				os.Exit(1)
			}
			return
		case "c", "changes":
			doChanges(argAt(args, 1))
			return
		case "f", "fetch":
			printFetch(true)
			return
		case "fm", "fetch-min":
			printFetch(false)
			return
		case "rm", "remove":
			if err := doRemove(); err != nil {
				os.Exit(1)
			}
			return
		case "s", "static":
			forceStatic = true
		case "a", "animate":
			forceAnim = true
		default:
			fmt.Fprintf(os.Stderr, "Неизвестный аргумент: %s\n\n", args[0])
			printHelp()
			os.Exit(1)
		}
	}

	// Короткие команды: разовая проверка/создание (закешировано маркером).
	ensureShortCommands(false)

	detect.Detect()
	results := detect.Normalize()

	if forceStatic || !isatty() {
		printStatic(results)
		return
	}

	// Анимация полосок: принудительно по флагу -a либо один раз при первом
	// интерактивном запуске (закешировано маркером).
	animate := forceAnim || firstRunAnimation()

	// Проверка обновлений стартует в фоне — чтобы утилита открывалась сразу,
	// а уведомление показывалось уже после закрытия просмотра.
	updCh := startBackgroundCheck()

	// Прогрев суперфетча в фоне: пока пользователь смотрит список, тяжёлый сбор
	// (подсчёт пакетов и т.п.) уже идёт — к открытию суперфетча кеш готов.
	go tui.WarmSuperfetch()

	model := tui.NewModel(results, animate)
	p := tea.NewProgram(model, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Ошибка: %v\n", err)
		os.Exit(1)
	}

	// После закрытия просмотра — результат фоновой проверки обновлений.
	finishBackgroundCheck(updCh)
}

// argAt возвращает args[i] или "" — для необязательных аргументов флага.
func argAt(args []string, i int) string {
	if i < len(args) {
		return args[i]
	}
	return ""
}

// isatty → true, только если и ввод, и вывод — терминал: при `lx > файл`
// или `lx | less` нужен статический режим, а не escape-коды TUI в пайпе.
func isatty() bool {
	for _, f := range []*os.File{os.Stdin, os.Stdout} {
		fi, err := f.Stat()
		if err != nil || fi.Mode()&os.ModeCharDevice == 0 {
			return false
		}
	}
	return true
}

func printVersion() {
	fmt.Printf("lxprofile %s\n", Version)
}

// printFetch печатает суперфетч статично, как fastfetch: full — подробный вид,
// иначе краткий. Профилирование архетипов при этом не запускается.
func printFetch(full bool) {
	fmt.Println(tui.RenderSuperfetch(full))
}

func printStatic(results []detect.ArchetypeResult) {
	fmt.Println()
	fmt.Println("  \033[1m🐧 l x p r o f i l e\033[0m")
	fmt.Println("  \033[2m╶─────────────────╴\033[0m")
	fmt.Println()
	for _, r := range results {
		desc := data.Describe(r.Key, r.NormScore)
		if strings.TrimSpace(desc) == "" {
			continue
		}
		color := "\033[2m"
		if r.NormScore >= 80 {
			color = "\033[32m"
		} else if r.NormScore >= 50 {
			color = "\033[33m"
		}
		fmt.Printf("%s%s\033[0m\n", color, desc)
	}
	fmt.Println()
	fmt.Println("\033[1m\U0001f50d Найдено в системе:\033[0m")
	fmt.Println("\033[2m  " + collectReasons(results) + "\033[0m")
	fmt.Println()
}

func collectReasons(results []detect.ArchetypeResult) string {
	seen := map[string]bool{}
	var reasons []string
	for _, r := range results {
		for _, reason := range strings.Split(r.Reason, ",") {
			reason = strings.TrimSpace(reason)
			if reason != "" && !seen[reason] {
				seen[reason] = true
				reasons = append(reasons, reason)
			}
		}
	}
	return strings.Join(reasons, " · ")
}

func printHelp() {
	fmt.Print(`lxprofile — Linux Psychological Profiler

ИСПОЛЬЗОВАНИЕ:
  lxprofile [ОПЦИЯ]

  Без опций запускает профайлер: в терминале — интерактивно,
  в пайпе или не-TTY — статической сводкой.

ОПЦИИ:
  -s, --static          статическая сводка вместо интерактивного режима
  -a, --animate         запустить с анимацией заполнения полосок
  -f, --fetch           суперфетч статично (подробный вид), как fastfetch
      --fm, --fetch-min суперфетч статично (краткий вид)
  -u, --update [ВЕР]    обновить с GitHub (или откатиться на версию ВЕР)
  -c, --changes [ВЕР]   changelog: указанной версии или 5 последних
      --rm, --remove    удалить lxprofile из системы
  -v, --version         показать версию
  -h, --help            показать эту справку

  Флаги можно писать без дефиса (lx u, lx static) и слитно с короткой командой
  lx (lxu = lx u = lx -u; так же lxs, lxa, lxf, lxfm, lxv, lxh, lxc, lxrm).

УПРАВЛЕНИЕ (интерактивный режим):
  ↑, k                  листать вверх
  ↓, j                  листать вниз
  →, l                  следующий режим (список / статистика / суперфетч)
  ←, h                  предыдущий режим
  m                     следующий режим
  g                     к первому архетипу
  G                     к последнему архетипу
  q                     выход

Разработано Axawys.
GitHub: https://github.com/Axawys/lxprofiler
`)
}
