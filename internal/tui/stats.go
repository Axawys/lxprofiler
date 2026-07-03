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
	VimCount     int
	NvimCount    int
	NanoCount    int
	EmacsCount   int
	MicroCount   int
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
		switch word {
		case "vim", "vi":
			result.VimCount++
		case "nvim":
			result.NvimCount++
		case "nano":
			result.NanoCount++
		case "emacs":
			result.EmacsCount++
		case "micro":
			result.MicroCount++
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

func editorWin(vim, nvim, nano, emacs, micro int) string {
	counts := []struct {
		name  string
		count int
	}{
		{"vim", vim}, {"nvim", nvim}, {"nano", nano}, {"emacs", emacs}, {"micro", micro},
	}
	total := 0
	bestn := -1
	for _, c := range counts {
		total += c.count
		if c.count > bestn {
			bestn = c.count
		}
	}
	if total == 0 {
		return "все мимо — GUI?"
	}
	ties := 0
	best := ""
	for _, c := range counts {
		if c.count == bestn {
			ties++
			best = c.name
		}
	}
	if ties > 1 {
		return "ничья"
	}
	return "победил " + best
}

func browserBattleLine(browserCache map[string]int) string {
	type browser struct {
		name string
		kb   int
	}
	var browsers []browser
	max := 0
	for name, kb := range browserCache {
		if kb > 0 {
			browsers = append(browsers, browser{name: name, kb: kb})
			if kb > max {
				max = kb
			}
		}
	}
	if max == 0 {
		return ""
	}
	sort.Slice(browsers, func(i, j int) bool {
		return browsers[i].kb > browsers[j].kb
	})
	var parts []string
	ties := 0
	win := ""
	for _, b := range browsers {
		parts = append(parts, fmt.Sprintf("%s %s", b.name, humanKB(b.kb)))
		if b.kb == max {
			ties++
			win = b.name
		}
	}
	line := strings.Join(parts, " : ")
	if ties > 1 {
		line += "  → ничья"
	} else {
		line += "  → чаще всех " + win
	}
	return line
}

func renderStats(m Model) string {
	var sb strings.Builder

	s := computeStats()

	if !s.HasHistory {
		sb.WriteString(dimStyle.Render("  История команд пуста или недоступна."))
		sb.WriteString("\n")
		sb.WriteString(dimStyle.Render("  Подсказка: включи HISTTIMEFORMAT — и время будет точнее."))
		sb.WriteString("\n\n")
		sb.WriteString(dimStyle.Render("  ↑↓ — листать · ←→ — режим · q — выход"))
		return sb.String()
	}

	spanSec := s.SpanDays * 86400
	sb.WriteString(fmt.Sprintf("  В истории %s команд, %s уникальных%s",
		boldStyle.Render(fmt.Sprintf("%d", s.TotalCmds)),
		boldStyle.Render(fmt.Sprintf("%d", s.UniqueCmds)),
		dimStyle.Render(fmt.Sprintf(" (охват ~%d дн.)", s.SpanDays))))
	sb.WriteString("\n\n")

	sb.WriteString("  Любимые команды:")
	sb.WriteString("\n")
	if len(s.TopCmds) > 0 {
		sb.WriteString(fmt.Sprintf("    1. %s — %s× %s\n",
			greenStyle.Render(s.TopCmds[0].Cmd),
			boldStyle.Render(fmt.Sprintf("%d", s.TopCmds[0].Count)),
			dimStyle.Render(fmt.Sprintf("(%s)", freq(s.TopCmds[0].Count, spanSec)))))
	}
	if len(s.TopCmds) > 1 {
		sb.WriteString(fmt.Sprintf("    2. %s — %s× %s\n",
			boldStyle.Render(s.TopCmds[1].Cmd),
			boldStyle.Render(fmt.Sprintf("%d", s.TopCmds[1].Count)),
			dimStyle.Render(fmt.Sprintf("(%s)", freq(s.TopCmds[1].Count, spanSec)))))
	}
	if len(s.TopCmds) > 2 {
		sb.WriteString(fmt.Sprintf("    3. %s — %s× %s\n",
			boldStyle.Render(s.TopCmds[2].Cmd),
			boldStyle.Render(fmt.Sprintf("%d", s.TopCmds[2].Count)),
			dimStyle.Render(fmt.Sprintf("(%s)", freq(s.TopCmds[2].Count, spanSec)))))
	}

	sb.WriteString(fmt.Sprintf("  fastfetch/neofetch: %s× %s — %s\n",
		boldStyle.Render(fmt.Sprintf("%d", s.FFCount)),
		dimStyle.Render(fmt.Sprintf("(%s)", freq(s.FFCount, spanSec))),
		dimStyle.Render(ffQuip(spanSec, s.FFCount))))

	sb.WriteString(fmt.Sprintf("  Обновления: %s× %s — %s\n",
		boldStyle.Render(fmt.Sprintf("%d", s.UpdCount)),
		dimStyle.Render(fmt.Sprintf("(%s)", freq(s.UpdCount, spanSec))),
		dimStyle.Render(updQuip(spanSec, s.UpdCount))))

	sb.WriteString(fmt.Sprintf("  sudo/doas: %s× %s — %s\n",
		boldStyle.Render(fmt.Sprintf("%d", s.SudoCount)),
		dimStyle.Render(fmt.Sprintf("(%s)", freq(s.SudoCount, spanSec))),
		dimStyle.Render(sudoQuip(spanSec, s.SudoCount))))

	sb.WriteString(fmt.Sprintf("  rm -rf: %s× — %s\n",
		boldStyle.Render(fmt.Sprintf("%d", s.RMRFCount)),
		dimStyle.Render(rmrfQuip(spanSec, s.RMRFCount))))

	sb.WriteString(fmt.Sprintf("  Опечаток поймано: %s%s\n",
		boldStyle.Render(fmt.Sprintf("%d", s.TypoCount)),
		dimStyle.Render(" (sl, gti, claer, cd..…)")))

	sb.WriteString(fmt.Sprintf("  Редактор-война: %s %d : %s %d : %s %d : %s %d : %s %d  %s\n",
		dimStyle.Render("vim"), s.VimCount,
		dimStyle.Render("nvim"), s.NvimCount,
		dimStyle.Render("nano"), s.NanoCount,
		dimStyle.Render("emacs"), s.EmacsCount,
		dimStyle.Render("micro"), s.MicroCount,
		dimStyle.Render("→ "+editorWin(s.VimCount, s.NvimCount, s.NanoCount, s.EmacsCount, s.MicroCount))))

	bb := browserBattleLine(s.BrowserCache)
	if bb != "" {
		sb.WriteString(fmt.Sprintf("  Битва браузеров: %s\n", bb))
	}

	sb.WriteString("\n")
	sb.WriteString(dimStyle.Render("  ↑↓ — листать · ←→ — режим · q — выход"))

	return sb.String()
}
