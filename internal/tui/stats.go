package tui

import (
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"
)

type StatsResult struct {
	TotalCmds    int
	UniqueCmds   int
	TopCmds      []CmdCount
	FFCount      int
	SudoCount    int
	UpdCount     int
	RMRFCount    int
	TypoCount    int
	Editors      map[string]int // редактор → сколько раз встретился в истории
	SpanDays     int
	HasHistory   bool
	BrowserCache map[string]int
}

type CmdCount struct {
	Cmd   string
	Count int
}

func cacheKB(dirs ...string) int {
	total := 0
	for _, d := range dirs {
		if info, err := os.Stat(d); err == nil && info.IsDir() {
			out, err := exec.Command("du", "-sk", d).Output()
			if err == nil {
				fields := strings.Fields(string(out))
				if len(fields) > 0 {
					if kb, err := strconv.Atoi(fields[0]); err == nil {
						total += kb
					}
				}
			}
		}
	}
	return total
}

func humanKB(kb int) string {
	if kb >= 1048576 {
		return fmt.Sprintf("%d.%dG", kb/1048576, (kb%1048576)*10/1048576)
	} else if kb >= 1024 {
		return fmt.Sprintf("%dM", kb/1024)
	}
	return fmt.Sprintf("%dK", kb)
}

func humanInterval(secs int) string {
	if secs <= 0 {
		return "—"
	} else if secs >= 86400 {
		return fmt.Sprintf("раз в %d дн.", secs/86400)
	} else if secs >= 3600 {
		return fmt.Sprintf("раз в %d ч.", secs/3600)
	} else if secs >= 60 {
		return fmt.Sprintf("раз в %d мин.", secs/60)
	}
	return fmt.Sprintf("раз в %d сек.", secs)
}

func freq(count, spanSec int) string {
	if count <= 0 || spanSec <= 0 {
		return "—"
	}
	return humanInterval(spanSec / count)
}

var (
	statsOnce   sync.Once
	statsCached StatsResult
)

// computeStats кеширует результат: сбор статистики (в т.ч. du по кешам браузеров)
// выполняется один раз, а не на каждый кадр отрисовки и расчёт размеров.
func computeStats() StatsResult {
	statsOnce.Do(func() { statsCached = computeStatsRaw() })
	return statsCached
}

// typos — распространённые опечатки команд (учитываются по первому слову).
var typos = map[string]bool{
	"sl": true, "gti": true, "claer": true, "clera": true, "grpe": true,
	"grep-": true, "gerp": true, "mkdri": true, "pythno": true, "cd..": true,
	"cd...": true, "ecoh": true, "sudp": true, "suod": true, "vin": true,
	"whcih": true, "cta": true, "nvin": true, "lls": true, "sl-": true,
}

// updatePatterns — подстроки команд обновления для разных дистрибутивов.
var updatePatterns = []string{
	"pacman -s", "yay -s", "paru -s", "apt update", "apt upgrade",
	"apt-get update", "apt-get upgrade", "dnf update", "dnf upgrade",
	"zypper up", "zypper dup", "emerge -u", "emerge --sync", "nixos-rebuild",
	"nix-channel --update", "flatpak update", "xbps-install -su", "apk upgrade",
}

// fetchTools — базовые команды-«фечи» (вывод инфы о системе).
var fetchTools = map[string]bool{
	"fastfetch": true, "neofetch": true, "screenfetch": true, "pfetch": true,
	"hyfetch": true, "macchina": true, "nerdfetch": true, "ufetch": true,
	"paleofetch": true, "cpufetch": true,
}

// editorNames — команда запуска редактора → отображаемое имя. Явные варианты
// одного редактора (vi/gvim → vim, emacsclient → emacs) сливаются вместе.
var editorNames = map[string]string{
	"vim": "vim", "vi": "vim", "gvim": "vim", "vimx": "vim", "nvi": "vim",
	"nvim": "nvim", "neovim": "nvim", "nvim-qt": "nvim",
	"nano": "nano",
	"emacs": "emacs", "emacsclient": "emacs",
	"micro": "micro",
	"hx": "helix", "helix": "helix",
	"kak": "kakoune",
	"ne": "ne", "joe": "joe", "jed": "jed", "mg": "mg", "vis": "vis", "ed": "ed",
	"pico": "pico", "mcedit": "mcedit",
	"code": "vscode", "code-insiders": "vscode", "codium": "vscodium",
	"subl": "sublime", "zed": "zed", "zeditor": "zed",
	"gedit": "gedit", "kate": "kate", "kwrite": "kwrite",
	"mousepad": "mousepad", "leafpad": "leafpad", "geany": "geany",
}

// aliasRe разбирает объявления alias/abbr: имя и значение (bash/zsh/fish).
var aliasRe = regexp.MustCompile(`(?i)^\s*(?:alias|abbr)(?:\s+-\S+)*\s+([\w.-]+)\s*=?\s*(.*)$`)

// fetchAliasNames возвращает имена алиасов/abbr, разрешающихся в fetch-инструмент.
func fetchAliasNames(configText string) []string {
	var names []string
	for _, line := range strings.Split(configText, "\n") {
		m := aliasRe.FindStringSubmatch(line)
		if m == nil {
			continue
		}
		name := strings.ToLower(strings.TrimSpace(m[1]))
		value := strings.ToLower(strings.Trim(strings.TrimSpace(m[2]), `'"`))
		for _, w := range strings.FieldsFunc(value, func(r rune) bool {
			return strings.ContainsRune(" \t|&;\"'", r)
		}) {
			if fetchTools[w] {
				names = append(names, name)
				break
			}
		}
	}
	return names
}

// fetchAliases читает конфиги шелла и собирает алиасы, ведущие на fetch-инструмент.
func fetchAliases() map[string]bool {
	set := map[string]bool{}
	home := os.Getenv("HOME")
	if home == "" {
		home = "/root"
	}
	for _, f := range []string{
		"/.bashrc", "/.bash_aliases", "/.zshrc", "/.zsh_aliases", "/.aliases",
		"/.profile", "/.config/fish/config.fish",
	} {
		data, err := os.ReadFile(home + f)
		if err != nil {
			continue
		}
		for _, name := range fetchAliasNames(string(data)) {
			set[name] = true
		}
	}
	return set
}

// extractCommand достаёт саму команду из строки истории любого шелла:
// zsh extended (": <ts>:<dur>;cmd"), fish ("- cmd: cmd"), bash/zsh plain ("cmd").
// Служебные строки (fish when:/paths:, bash-таймстампы #<ts>) отбрасываются.
func extractCommand(line string) string {
	line = strings.TrimRight(line, "\r")
	if strings.HasPrefix(line, "#") {
		if _, err := strconv.Atoi(strings.TrimSpace(line[1:])); err == nil {
			return "" // bash HISTTIMEFORMAT: строка "#1609459200"
		}
	}
	if strings.HasPrefix(line, ": ") {
		if i := strings.IndexByte(line, ';'); i >= 0 {
			return strings.TrimSpace(line[i+1:]) // zsh extended history
		}
	}
	if strings.HasPrefix(line, "- cmd: ") {
		return strings.TrimSpace(line[len("- cmd: "):]) // fish
	}
	if strings.HasPrefix(line, "  ") || strings.HasPrefix(line, "- ") {
		return "" // служебные строки fish
	}
	return strings.TrimSpace(line)
}

// tallyCommands считает все счётчики истории по списку команд (первое слово
// команды, sudo/doas «прозрачны»). Вынесено отдельно, чтобы логика была
// тестируемой без чтения файлов.
func tallyCommands(commands []string, fetchAlias map[string]bool, result *StatsResult) {
	result.TotalCmds = len(commands)
	if result.TotalCmds == 0 {
		return
	}
	result.HasHistory = true
	if result.Editors == nil {
		result.Editors = map[string]int{}
	}

	cmdCount := map[string]int{}
	seen := map[string]bool{}
	for _, c := range commands {
		fields := strings.Fields(c)
		if len(fields) == 0 {
			continue
		}
		head := strings.ToLower(fields[0])
		if head == "sudo" || head == "doas" {
			result.SudoCount++
		}
		word := head
		if (head == "sudo" || head == "doas") && len(fields) > 1 {
			word = strings.ToLower(fields[1])
		}
		cmdCount[word]++
		seen[word] = true
		if disp, ok := editorNames[word]; ok {
			result.Editors[disp]++
		}
		if typos[word] {
			result.TypoCount++
		}
		// fetch-инструмент напрямую или через алиас (ff → fastfetch)
		if fetchTools[word] || fetchAlias[word] {
			result.FFCount++
		}
	}
	result.UniqueCmds = len(seen)

	for cmd, count := range cmdCount {
		result.TopCmds = append(result.TopCmds, CmdCount{Cmd: cmd, Count: count})
	}
	sort.Slice(result.TopCmds, func(i, j int) bool {
		return result.TopCmds[i].Count > result.TopCmds[j].Count
	})
	if len(result.TopCmds) > 3 {
		result.TopCmds = result.TopCmds[:3]
	}

	// Подстрочные счётчики — по тексту команд (без метаданных истории).
	cmdText := strings.ToLower(strings.Join(commands, "\n"))
	result.RMRFCount = strings.Count(cmdText, "rm -rf")
	for _, p := range updatePatterns {
		result.UpdCount += strings.Count(cmdText, p)
	}
}

func computeStatsRaw() StatsResult {
	home := os.Getenv("HOME")
	if home == "" {
		home = "/root"
	}

	result := StatsResult{
		BrowserCache: make(map[string]int),
	}

	var raw string
	for _, f := range []string{
		home + "/.bash_history",
		home + "/.zsh_history",
		home + "/.local/share/fish/fish_history",
	} {
		if d, err := os.ReadFile(f); err == nil {
			raw += "\n" + string(d)
		}
	}

	if raw == "" {
		return result
	}

	// Достаём реальные команды из истории (с учётом форматов zsh/fish/bash)
	// и считаем по ним все счётчики.
	var commands []string
	for _, line := range strings.Split(raw, "\n") {
		if c := extractCommand(line); c != "" {
			commands = append(commands, c)
		}
	}
	tallyCommands(commands, fetchAliases(), &result)
	if !result.HasHistory {
		return result
	}

	// Time span
	var timestamps []int
	re := regexp.MustCompile(`^: ([0-9]+):`)
	for _, line := range strings.Split(raw, "\n") {
		if matches := re.FindStringSubmatch(line); len(matches) > 1 {
			if ts, err := strconv.Atoi(matches[1]); err == nil {
				timestamps = append(timestamps, ts)
			}
		}
	}
	if len(timestamps) >= 2 {
		sort.Ints(timestamps)
		result.SpanDays = (timestamps[len(timestamps)-1] - timestamps[0]) / 86400
	} else {
		// Fallback to install time
		installEpoch := 0
		if out, err := exec.Command("stat", "-c", "%W", "/").Output(); err == nil {
			fmt.Sscanf(strings.TrimSpace(string(out)), "%d", &installEpoch)
		}
		if installEpoch > 0 {
			result.SpanDays = (int(time.Now().Unix()) - installEpoch) / 86400
		}
	}

	// Browser cache
	cache := home + "/.cache"
	varApp := home + "/.var/app"
	result.BrowserCache["Firefox"] = cacheKB(cache+"/mozilla", varApp+"/org.mozilla.firefox/cache")
	result.BrowserCache["Chromium"] = cacheKB(cache+"/chromium", varApp+"/org.chromium.Chromium/cache")
	result.BrowserCache["Chrome"] = cacheKB(cache+"/google-chrome", varApp+"/com.google.Chrome/cache")
	result.BrowserCache["Brave"] = cacheKB(cache+"/BraveSoftware", varApp+"/com.brave.Browser/cache")
	result.BrowserCache["Vivaldi"] = cacheKB(cache+"/vivaldi", varApp+"/com.vivaldi.Vivaldi/cache")
	result.BrowserCache["Opera"] = cacheKB(cache+"/opera", varApp+"/com.opera.Opera/cache")
	result.BrowserCache["Edge"] = cacheKB(cache+"/microsoft-edge", varApp+"/com.microsoft.Edge/cache")
	result.BrowserCache["LibreWolf"] = cacheKB(cache+"/librewolf", varApp+"/io.gitlab.librewolf-community/cache")
	result.BrowserCache["Waterfox"] = cacheKB(cache+"/waterfox", varApp+"/net.waterfox.waterfox/cache")
	result.BrowserCache["Zen"] = cacheKB(cache+"/zen", varApp+"/app.zen_browser.zen/cache")
	result.BrowserCache["Yandex"] = cacheKB(cache+"/yandex-browser", varApp+"/ru.yandex.Browser/cache")
	result.BrowserCache["qutebrowser"] = cacheKB(cache+"/qutebrowser", varApp+"/org.qutebrowser.qutebrowser/cache")

	return result
}

func ffQuip(spanSec, ffCount int) string {
	if ffCount == 0 {
		return "ни разу, аскет"
	}
	if spanSec <= 0 {
		return "частоту не определить"
	}
	i := spanSec / ffCount
	switch {
	case i >= 2592000:
		return "очень редко"
	case i >= 604800:
		return "иногда любуешься системой"
	case i >= 86400:
		return "почти ежедневный ритуал"
	case i >= 3600:
		return "по нескольку раз в день"
	default:
		return "это уже зависимость))"
	}
}

func updQuip(spanSec, updCount int) string {
	if updCount == 0 {
		return "ни разу — смело"
	}
	if spanSec <= 0 {
		return "частоту не определить"
	}
	i := spanSec / updCount
	switch {
	case i >= 2592000:
		return "обновляешься редко — стабильность важнее"
	case i >= 604800:
		return "апдейт по выходным"
	case i >= 86400:
		return "держишь систему свежей"
	default:
		return "апдейт — это медитация"
	}
}

func rmrfQuip(spanSec, rmrfCount int) string {
	if rmrfCount == 0 {
		return "аккуратно"
	}
	if spanSec <= 0 {
		return "частоту не определить"
	}
	i := spanSec / rmrfCount
	switch {
	case i >= 2592000:
		return "редко, но метко"
	case i >= 604800:
		return "бывает"
	case i >= 86400:
		return "живёшь опасно"
	default:
		return "как ты ещё жив?"
	}
}

func sudoQuip(spanSec, sudoCount int) string {
	if sudoCount == 0 {
		return "живёшь без рута"
	}
	if spanSec <= 0 {
		return "частоту не определить"
	}
	i := spanSec / sudoCount
	switch {
	case i >= 86400:
		return "рут по праздникам"
	case i >= 3600:
		return "уверенно у руля"
	default:
		return "практически root"
	}
}

func renderStats(m Model) string {
	return pinTabBar(statsBody(m), m)
}

// statsBody собирает забавную статистику панелями: шапка со сводкой, «Топ
// команд», «Активность» и «Битвы». Панель вкладок добавляет pinTabBar.
func statsBody(m Model) []string {
	s := computeStats()
	if !s.HasHistory {
		return []string{
			dimStyle.Render("  История команд пуста или недоступна."),
			dimStyle.Render("  Подсказка: включи HISTTIMEFORMAT — и время будет точнее."),
		}
	}

	spanSec := s.SpanDays * 86400

	// Топ команд — ранг, имя (выровнено), счётчик, частота.
	var top []string
	for i, c := range s.TopCmds {
		name := padRight(c.Cmd, 12)
		if i == 0 {
			name = greenStyle.Render(name)
		} else {
			name = boldStyle.Render(name)
		}
		top = append(top, fmt.Sprintf("%d. %s %s  %s",
			i+1, name,
			boldStyle.Render(padLeft(fmt.Sprintf("%d×", c.Count), 5)),
			dimStyle.Render(freq(c.Count, spanSec))))
	}

	// Активность — метка · счётчик · частота · комментарий.
	var act []string
	metric := func(label string, count int, freqTxt, quip string) {
		act = append(act, fmt.Sprintf("%s %s %s %s",
			padRight(label, 11),
			boldStyle.Render(padLeft(fmt.Sprintf("%d×", count), 5)),
			dimStyle.Render(padRight(freqTxt, 13)),
			dimStyle.Render("· "+quip)))
	}
	metric("fastfetch", s.FFCount, freq(s.FFCount, spanSec), ffQuip(spanSec, s.FFCount))
	metric("обновления", s.UpdCount, freq(s.UpdCount, spanSec), updQuip(spanSec, s.UpdCount))
	metric("sudo/doas", s.SudoCount, freq(s.SudoCount, spanSec), sudoQuip(spanSec, s.SudoCount))
	metric("rm -rf", s.RMRFCount, freq(s.RMRFCount, spanSec), rmrfQuip(spanSec, s.RMRFCount))
	metric("опечатки", s.TypoCount, "", "sl, gti, claer, cd..…")

	// Битвы — редакторы и браузеры.
	var battles []string
	var edItems []battleItem
	for name, cnt := range s.Editors {
		edItems = append(edItems, battleItem{name, cnt, strconv.Itoa(cnt)})
	}
	if line := battleLine("Редактор-война", edItems); line != "" {
		battles = append(battles, line)
	} else {
		battles = append(battles, boldStyle.Render(padRight("Редактор-война", 15))+
			"  "+dimStyle.Render("все мимо — GUI?"))
	}
	var brItems []battleItem
	for name, kb := range s.BrowserCache {
		brItems = append(brItems, battleItem{name, kb, humanKB(kb)})
	}
	if line := battleLine("Битва браузеров", brItems); line != "" {
		battles = append(battles, line)
	}

	// Единая внутренняя ширина панелей — по самому широкому содержимому.
	inner := contentWidth([]string{"Топ команд", "Активность", "Битвы"}, top, act, battles)
	if max := m.width - 6; inner > max {
		inner = max
	}

	content := []string{
		fmt.Sprintf("  %s команд · %s уникальных · %s",
			boldStyle.Render(fmt.Sprintf("%d", s.TotalCmds)),
			boldStyle.Render(fmt.Sprintf("%d", s.UniqueCmds)),
			dimStyle.Render(fmt.Sprintf("охват ~%d дн.", s.SpanDays))),
		"",
	}
	content = append(content, panel("Топ команд", top, inner)...)
	content = append(content, panel("Активность", act, inner)...)
	content = append(content, panel("Битвы", battles, inner)...)
	return content
}

type battleItem struct {
	name  string
	val   int
	label string // что показать рядом с именем (счётчик или размер)
}

// battleLine формирует строку «войны»: ненулевые участники по убыванию (при
// равенстве — по имени, чтобы порядок был стабилен), лидер подсвечен. Пусто,
// если участников нет.
func battleLine(title string, items []battleItem) string {
	var live []battleItem
	max := 0
	for _, it := range items {
		if it.val > 0 {
			live = append(live, it)
			if it.val > max {
				max = it.val
			}
		}
	}
	if len(live) == 0 {
		return ""
	}
	sort.Slice(live, func(i, j int) bool {
		if live[i].val != live[j].val {
			return live[i].val > live[j].val
		}
		return live[i].name < live[j].name
	})

	var parts []string
	ties, leader := 0, ""
	for _, it := range live {
		parts = append(parts, it.name+" "+boldStyle.Render(it.label))
		if it.val == max {
			ties++
			leader = it.name
		}
	}
	winner := "→ " + leader
	if ties > 1 {
		winner = "→ ничья"
	}
	return fmt.Sprintf("%s  %s  %s",
		boldStyle.Render(padRight(title, 15)),
		strings.Join(parts, dimStyle.Render(" · ")),
		greenStyle.Render(winner))
}
