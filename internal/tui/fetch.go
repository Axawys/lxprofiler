package tui

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/charmbracelet/lipgloss"
)

// fetchInfo — данные для режима «суперфетч» (аналог fastfetch, но уже
// сконфигурированный): секции OS / Hardware / Software.
type fetchInfo struct {
	User, Host string

	// OS
	Distro      string // PRETTY_NAME из os-release
	DistroID    string // ID из os-release — выбор логотипа
	Kernel      string
	Uptime      string
	InstallDate string
	AgeDays     int
	HasInstall  bool

	// Hardware
	CPU        string
	Cores      int
	GPU        string
	Resolution string // "1920x1080@60Hz" (или без @Гц, если частота неизвестна)
	RAMUsed    string
	RAMTotal   string
	RAMPct     int
	DiskPct    int
	DiskUsed   string
	DiskTotal  string
	HasDisk    bool

	// Software
	Shell    string
	DEWM     string
	Packages int
	FlatpakN int
	SnapN    int
}

var (
	fetchOnce   sync.Once
	fetchCached fetchInfo
)

// computeFetch собирает системную инфу один раз (df/stat/pkg не на каждый кадр).
func computeFetch() fetchInfo {
	fetchOnce.Do(func() { fetchCached = gatherFetch() })
	return fetchCached
}

func gatherFetch() fetchInfo {
	fi := fetchInfo{User: "user", Host: "linux", Distro: "Linux", DistroID: "linux"}
	if u := os.Getenv("USER"); u != "" {
		fi.User = u
	} else if u := os.Getenv("LOGNAME"); u != "" {
		fi.User = u
	}
	if h, err := os.Hostname(); err == nil && h != "" {
		fi.Host = h
	}

	// ── OS ─────────────────────────────────────────────
	if d, err := os.ReadFile("/etc/os-release"); err == nil {
		for _, line := range strings.Split(string(d), "\n") {
			switch {
			case strings.HasPrefix(line, "PRETTY_NAME="):
				fi.Distro = trimOSValue(line[len("PRETTY_NAME="):])
			case strings.HasPrefix(line, "ID="):
				fi.DistroID = strings.ToLower(trimOSValue(line[len("ID="):]))
			}
		}
	}
	if d, err := os.ReadFile("/proc/sys/kernel/osrelease"); err == nil {
		fi.Kernel = strings.TrimSpace(string(d))
	}
	if d, err := os.ReadFile("/proc/uptime"); err == nil {
		if f := strings.Fields(string(d)); len(f) > 0 {
			sec := 0
			fmt.Sscanf(f[0], "%d", &sec)
			fi.Uptime = fmtUptime(sec)
		}
	}
	// Дата установки: время рождения корня, иначе mtime /etc/machine-id.
	installEpoch := 0
	if out, err := exec.Command("stat", "-c", "%W", "/").Output(); err == nil {
		fmt.Sscanf(strings.TrimSpace(string(out)), "%d", &installEpoch)
	}
	if installEpoch == 0 {
		if out, err := exec.Command("stat", "-c", "%Y", "/etc/machine-id").Output(); err == nil {
			fmt.Sscanf(strings.TrimSpace(string(out)), "%d", &installEpoch)
		}
	}
	if installEpoch > 0 {
		fi.InstallDate = time.Unix(int64(installEpoch), 0).Format("2006-01-02")
		fi.AgeDays = (int(time.Now().Unix()) - installEpoch) / 86400
		fi.HasInstall = true
	}

	// ── Hardware ───────────────────────────────────────
	if d, err := os.ReadFile("/proc/cpuinfo"); err == nil {
		for _, line := range strings.Split(string(d), "\n") {
			if strings.HasPrefix(line, "processor") {
				fi.Cores++
			}
			if fi.CPU == "" && strings.HasPrefix(line, "model name") {
				if i := strings.IndexByte(line, ':'); i >= 0 {
					fi.CPU = cleanCPU(strings.TrimSpace(line[i+1:]))
				}
			}
		}
	}
	var memTotal, memAvail int
	if d, err := os.ReadFile("/proc/meminfo"); err == nil {
		for _, line := range strings.Split(string(d), "\n") {
			f := strings.Fields(line)
			if len(f) >= 2 {
				switch f[0] {
				case "MemTotal:":
					fmt.Sscanf(f[1], "%d", &memTotal)
				case "MemAvailable:":
					fmt.Sscanf(f[1], "%d", &memAvail)
				}
			}
		}
	}
	if memTotal > 0 {
		used := memTotal - memAvail
		fi.RAMTotal = humanGiB(memTotal)
		fi.RAMUsed = humanGiB(used)
		fi.RAMPct = used * 100 / memTotal
	}
	fi.GPU = gpuModel()
	fi.Resolution = displayMode()
	// Диск на корне: df -h /, берём последнюю строку (на случай переноса).
	if out, err := exec.Command("df", "-h", "/").Output(); err == nil {
		lines := strings.Split(strings.TrimSpace(string(out)), "\n")
		if len(lines) >= 2 {
			f := strings.Fields(lines[len(lines)-1])
			if len(f) >= 5 {
				fi.DiskTotal = f[1]
				fi.DiskUsed = f[2]
				if p, err := strconv.Atoi(strings.TrimSuffix(f[4], "%")); err == nil {
					fi.DiskPct = p
					fi.HasDisk = true
				}
			}
		}
	}

	// ── Software ───────────────────────────────────────
	fi.Shell = filepath.Base(os.Getenv("SHELL"))
	if fi.Shell == "" || fi.Shell == "." {
		fi.Shell = "—"
	}
	fi.DEWM = firstNonEmpty(os.Getenv("XDG_CURRENT_DESKTOP"), os.Getenv("DESKTOP_SESSION"))
	if fi.DEWM == "" {
		if os.Getenv("XDG_SESSION_TYPE") == "tty" {
			fi.DEWM = "TTY"
		} else {
			fi.DEWM = "—"
		}
	}
	fi.Packages = pkgCount()
	if cmdExists("flatpak") {
		if out, err := exec.Command("flatpak", "list", "--app").Output(); err == nil {
			fi.FlatpakN = countLines(string(out))
		}
	}
	if cmdExists("snap") {
		if out, err := exec.Command("snap", "list").Output(); err == nil {
			if n := countLines(string(out)) - 1; n > 0 { // минус строка-заголовок
				fi.SnapN = n
			}
		}
	}

	return fi
}

// countLines считает непустые строки вывода команды.
func countLines(s string) int {
	n := 0
	for _, line := range strings.Split(s, "\n") {
		if strings.TrimSpace(line) != "" {
			n++
		}
	}
	return n
}

// gpuModel достаёт название видеокарты через lspci (0 — если lspci нет).
func gpuModel() string {
	if !cmdExists("lspci") {
		return ""
	}
	out, err := exec.Command("lspci").Output()
	if err != nil {
		return ""
	}
	for _, line := range strings.Split(string(out), "\n") {
		if strings.Contains(line, "VGA compatible controller") ||
			strings.Contains(line, "3D controller") ||
			strings.Contains(line, "Display controller") {
			if i := strings.LastIndex(line, ": "); i >= 0 {
				return cleanGPU(line[i+2:])
			}
		}
	}
	return ""
}

// cleanGPU сокращает строку lspci до узнаваемого имени карты.
func cleanGPU(s string) string {
	s = strings.TrimSpace(s)
	if i := strings.Index(s, " (rev "); i >= 0 {
		s = s[:i]
	}
	// Народное имя карты обычно в последних скобках: [GeForce RTX 3070].
	if l := strings.LastIndex(s, "["); l >= 0 {
		if r := strings.Index(s[l:], "]"); r > 1 {
			return strings.TrimSpace(s[l+1 : l+r])
		}
	}
	for _, t := range []string{" Corporation", ", Inc.", " Inc."} {
		s = strings.ReplaceAll(s, t, "")
	}
	return strings.Join(strings.Fields(s), " ")
}

// displayMode возвращает разрешение и частоту активного дисплея: сначала через
// xrandr (X11/XWayland), затем через wlr-randr (wlroots). "" — если не удалось.
func displayMode() string {
	if cmdExists("xrandr") {
		if out, err := exec.Command("xrandr", "--current").Output(); err == nil {
			// строка вида: "   1920x1080     60.00*+  59.94 ..."
			re := regexp.MustCompile(`(\d+x\d+)\s+([\d.]+)\*`)
			if m := re.FindStringSubmatch(string(out)); m != nil {
				return fmt.Sprintf("%s@%sHz", m[1], trimHz(m[2]))
			}
		}
	}
	if cmdExists("wlr-randr") {
		if out, err := exec.Command("wlr-randr").Output(); err == nil {
			// строка вида: "    1920x1080 px, 60.000000 Hz (current)"
			re := regexp.MustCompile(`(\d+x\d+) px, ([\d.]+) Hz.*current`)
			if m := re.FindStringSubmatch(string(out)); m != nil {
				return fmt.Sprintf("%s@%sHz", m[1], trimHz(m[2]))
			}
		}
	}
	return ""
}

// trimHz округляет частоту "60.000000"/"59.94" до целого числа герц.
func trimHz(s string) string {
	if f, err := strconv.ParseFloat(s, 64); err == nil {
		return strconv.Itoa(int(f + 0.5))
	}
	return s
}

func trimOSValue(s string) string {
	return strings.Trim(strings.TrimSpace(s), `"`)
}

func firstNonEmpty(vals ...string) string {
	for _, v := range vals {
		if v != "" {
			return v
		}
	}
	return ""
}

func cmdExists(name string) bool { _, err := exec.LookPath(name); return err == nil }

// cleanCPU убирает маркетинговый мусор из «model name»: ®/™, частоту «@ 2.6GHz»
// и слова, дублирующие показанное число ядер (Processor / CPU / «8-Core»).
func cleanCPU(s string) string {
	for _, t := range []string{"(R)", "(r)", "(TM)", "(tm)", "(C)"} {
		s = strings.ReplaceAll(s, t, "")
	}
	if i := strings.Index(s, "@"); i >= 0 {
		s = s[:i]
	}
	var out []string
	for _, w := range strings.Fields(s) {
		lw := strings.ToLower(w)
		if lw == "processor" || lw == "cpu" || strings.HasSuffix(lw, "-core") {
			continue
		}
		out = append(out, w)
	}
	return strings.Join(out, " ")
}

// pkgCount возвращает число установленных пакетов для основного пакетного
// менеджера системы (0 — если менеджер не найден или запрос не удался).
func pkgCount() int {
	managers := []struct {
		bin  string
		args []string
	}{
		{"pacman", []string{"-Qq"}},
		{"dpkg-query", []string{"-f", ".\n", "-W"}},
		{"rpm", []string{"-qa"}},
		{"apk", []string{"info"}},
		{"xbps-query", []string{"-l"}},
	}
	for _, mgr := range managers {
		if !cmdExists(mgr.bin) {
			continue
		}
		out, err := exec.Command(mgr.bin, mgr.args...).Output()
		if err != nil {
			return 0
		}
		s := strings.TrimRight(string(out), "\n")
		if s == "" {
			return 0
		}
		return strings.Count(s, "\n") + 1
	}
	return 0
}

func fmtUptime(sec int) string {
	d := sec / 86400
	h := (sec % 86400) / 3600
	mnt := (sec % 3600) / 60
	var parts []string
	if d > 0 {
		parts = append(parts, fmt.Sprintf("%dд", d))
	}
	if h > 0 {
		parts = append(parts, fmt.Sprintf("%dч", h))
	}
	if mnt > 0 || len(parts) == 0 {
		parts = append(parts, fmt.Sprintf("%dм", mnt))
	}
	return strings.Join(parts, " ")
}

func humanGiB(kb int) string {
	return fmt.Sprintf("%.1f GiB", float64(kb)/1048576)
}

// pluralRu выбирает форму русского слова по числу (1 ядро / 2 ядра / 5 ядер).
func pluralRu(n int, one, few, many string) string {
	n10, n100 := n%10, n%100
	switch {
	case n10 == 1 && n100 != 11:
		return one
	case n10 >= 2 && n10 <= 4 && (n100 < 12 || n100 > 14):
		return few
	default:
		return many
	}
}

func humanAge(days int) string {
	if days >= 365 {
		y, d := days/365, days%365
		yw := fmt.Sprintf("%d %s", y, pluralRu(y, "год", "года", "лет"))
		if d == 0 {
			return yw
		}
		return fmt.Sprintf("%s %d дн.", yw, d)
	}
	return fmt.Sprintf("%d дн.", days)
}

func splitLogo(s string) []string {
	return strings.Split(strings.Trim(s, "\n"), "\n")
}

// Логотипы дистрибутивов (ASCII). Ключ — ID из os-release.
var distroLogos = map[string][]string{
	"arch": splitLogo(`
      /\
     /  \
    /\   \
   /  \   \
  / /\ \   \
 /_/  \_\   \
/_/    \_\___\`),
	"fedora": splitLogo(`
     ____
    / __ \_
   / /  \ \\
  | | () | |
  | |___/ /
   \     /
    \___/`),
	"ubuntu": splitLogo(`
      _
    _(_)_
   (_)'(_)
  (_) _ (_)
   (_)_(_)
     (_)`),
}

// genericLogo — запасной пингвин для неизвестных дистрибутивов.
var genericLogo = splitLogo(`
   .--.
  |o_o |
  |:_/ |
 //   \ \
(|     | )
/'\_  _/'\
\___)(___/`)

var logoStyles = map[string]lipgloss.Style{
	"arch":   lipgloss.NewStyle().Foreground(lipgloss.Color("12")),
	"fedora": lipgloss.NewStyle().Foreground(lipgloss.Color("4")),
	"ubuntu": lipgloss.NewStyle().Foreground(lipgloss.Color("208")),
}

var sectionStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("6")).Bold(true)

// fetchKV — строка «метка · значение» с выравниванием метки на ширину w
// (с отступом слева — для краткого вида без рамок).
func fetchKV(label, value string, w int) string {
	return "  " + dimStyle.Render(padRight(label, w)) + value
}

// kvRow — как fetchKV, но без левого отступа (отступ даёт рамка секции).
func kvRow(label, value string, w int) string {
	return dimStyle.Render(padRight(label, w)) + value
}

// boxSection рисует секцию в рамке с заголовком на верхней грани:
//
//	╭─ OS ───────────────╮
//	│ Distro   Arch Linux │
//	╰─────────────────────╯
//
// contentW — общая внутренняя ширина (одинаковая у всех секций, чтобы рамки
// были ровными). Рамка тусклая, заголовок — акцентный.
func boxSection(title string, rows []string, contentW int) []string {
	t := lipgloss.Width(title)
	if contentW < t+2 {
		contentW = t + 2
	}
	frame := dimStyle
	top := frame.Render("╭─ ") + sectionStyle.Render(title) +
		frame.Render(" "+strings.Repeat("─", contentW-t-1)+"╮")
	out := []string{top}
	for _, r := range rows {
		out = append(out, frame.Render("│ ")+padRight(r, contentW)+frame.Render(" │"))
	}
	out = append(out, frame.Render("╰"+strings.Repeat("─", contentW+2)+"╯"))
	return out
}

// contentWidth — максимальная видимая ширина среди строк секций и их заголовков
// (заголовку нужно место title+2 на верхней грани).
func contentWidth(titles []string, groups ...[]string) int {
	w := 0
	for _, g := range groups {
		for _, r := range g {
			if x := lipgloss.Width(r); x > w {
				w = x
			}
		}
	}
	for _, tt := range titles {
		if x := lipgloss.Width(tt) + 2; x > w {
			w = x
		}
	}
	return w
}

// usageBar — полоска заполнения (диск/ОЗУ): закрашенная часть по проценту
// (зелёный / жёлтый ≥70% / красный ≥90%), остаток — тусклый.
func usageBar(pct, width int) string {
	if pct < 0 {
		pct = 0
	}
	if pct > 100 {
		pct = 100
	}
	filled := pct * width / 100
	col := greenStyle
	switch {
	case pct >= 90:
		col = redStyle
	case pct >= 70:
		col = yellowStyle
	}
	return col.Render(strings.Repeat("█", filled)) + dimStyle.Render(strings.Repeat("░", width-filled))
}

func fetchHeader(fi fetchInfo) []string {
	return []string{
		boldStyle.Render(fi.User) + dimStyle.Render("@") + boldStyle.Render(fi.Host),
		dimStyle.Render(strings.Repeat("─", lipgloss.Width(fi.User)+1+lipgloss.Width(fi.Host))),
	}
}

// fetchInfoFull — подробный вид: секции OS / Hardware / Software в рамках.
func fetchInfoFull(fi fetchInfo) []string {
	kv := func(l, v string) string { return kvRow(l, v, 13) }

	osRows := []string{kv("Distro", fi.Distro), kv("Kernel", fi.Kernel)}
	if fi.Uptime != "" {
		osRows = append(osRows, kv("Uptime", fi.Uptime))
	}
	if fi.HasInstall {
		osRows = append(osRows,
			kv("Установлена", fi.InstallDate),
			kv("Возраст", humanAge(fi.AgeDays)))
	}

	var hwRows []string
	if fi.CPU != "" {
		hwRows = append(hwRows, kv("CPU", fmt.Sprintf("%s (%d %s)",
			fi.CPU, fi.Cores, pluralRu(fi.Cores, "ядро", "ядра", "ядер"))))
	}
	if fi.GPU != "" {
		hwRows = append(hwRows, kv("GPU", fi.GPU))
	}
	if fi.Resolution != "" {
		hwRows = append(hwRows, kv("Display", fi.Resolution))
	}
	if fi.RAMTotal != "" {
		hwRows = append(hwRows, kv("RAM", fmt.Sprintf("%s %d%%  %s / %s",
			usageBar(fi.RAMPct, 10), fi.RAMPct, fi.RAMUsed, fi.RAMTotal)))
	}
	if fi.HasDisk {
		hwRows = append(hwRows, kv("Disk", fmt.Sprintf("%s %d%%  %s / %s",
			usageBar(fi.DiskPct, 10), fi.DiskPct, fi.DiskUsed, fi.DiskTotal)))
	}

	swRows := []string{kv("Shell", fi.Shell), kv("DE/WM", fi.DEWM)}
	if fi.Packages > 0 {
		swRows = append(swRows, kv("Пакетов", strconv.Itoa(fi.Packages)))
	}
	if fi.FlatpakN > 0 {
		swRows = append(swRows, kv("Flatpak", strconv.Itoa(fi.FlatpakN)))
	}
	if fi.SnapN > 0 {
		swRows = append(swRows, kv("Snap", strconv.Itoa(fi.SnapN)))
	}

	cw := contentWidth([]string{"OS", "Hardware", "Software"}, osRows, hwRows, swRows)
	info := fetchHeader(fi)
	info = append(info, boxSection("OS", osRows, cw)...)
	info = append(info, boxSection("Hardware", hwRows, cw)...)
	info = append(info, boxSection("Software", swRows, cw)...)
	return info
}

// fetchInfoMinimal — краткий вид: без секций, только основное с барами.
func fetchInfoMinimal(fi fetchInfo) []string {
	info := fetchHeader(fi)
	kv := func(l, v string) string { return fetchKV(l, v, 8) }

	info = append(info, kv("Distro", fi.Distro))
	info = append(info, kv("Kernel", fi.Kernel))
	if fi.Uptime != "" {
		info = append(info, kv("Uptime", fi.Uptime))
	}
	if fi.CPU != "" {
		info = append(info, kv("CPU", fi.CPU))
	}
	if fi.RAMTotal != "" {
		info = append(info, kv("RAM", fmt.Sprintf("%s %d%%", usageBar(fi.RAMPct, 10), fi.RAMPct)))
	}
	if fi.HasDisk {
		info = append(info, kv("Disk", fmt.Sprintf("%s %d%%", usageBar(fi.DiskPct, 10), fi.DiskPct)))
	}
	return info
}

func renderFetch(m Model) string {
	fi := computeFetch()

	logo := genericLogo
	if l, ok := distroLogos[fi.DistroID]; ok {
		logo = l
	}
	style := dimStyle
	if s, ok := logoStyles[fi.DistroID]; ok {
		style = s
	}
	logoW := 0
	for _, l := range logo {
		if w := lipgloss.Width(l); w > logoW {
			logoW = w
		}
	}

	info := fetchInfoMinimal(fi)
	if m.fetchFull {
		info = fetchInfoFull(fi)
	}

	rows := len(logo)
	if len(info) > rows {
		rows = len(info)
	}
	var sb strings.Builder
	sb.WriteString("\n")
	for i := 0; i < rows; i++ {
		l, r := "", ""
		if i < len(logo) {
			l = logo[i]
		}
		if i < len(info) {
			r = info[i]
		}
		sb.WriteString("  " + style.Render(padRight(l, logoW)) + "   " + r + "\n")
	}
	sb.WriteString("\n")
	sb.WriteString(dimStyle.Render("  ↑↓ — кратко/подробно · ←→ — режим · q — выход"))
	return sb.String()
}
