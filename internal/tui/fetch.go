package tui

import (
	"fmt"
	"math"
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
	CPU       string
	Cores     int
	GPU       string
	Display   string // "2560x1440, 27\", 180Hz, 108ppi" (части опускаются, если неизвестны)
	RAMUsed   string
	RAMTotal  string
	RAMPct    int
	DiskPct   int
	DiskUsed  string
	DiskTotal string
	HasDisk   bool

	// Software
	Shell    string
	DEWM     string
	PkgKind  string // тип системного пакетного менеджера (rpm/deb/pacman/…)
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
				// Без уточнения в скобках: «Fedora Linux 44 (Workstation Edition)»
				// → «Fedora Linux 44».
				name := trimOSValue(line[len("PRETTY_NAME="):])
				if i := strings.Index(name, " ("); i >= 0 {
					name = name[:i]
				}
				fi.Distro = name
			case strings.HasPrefix(line, "ID="):
				fi.DistroID = strings.ToLower(trimOSValue(line[len("ID="):]))
			}
		}
	}
	if d, err := os.ReadFile("/proc/sys/kernel/osrelease"); err == nil {
		// Только версия ядра: «7.0.13-200.fc44.x86_64» → «7.0.13».
		k := strings.TrimSpace(string(d))
		if i := strings.IndexByte(k, '-'); i >= 0 {
			k = k[:i]
		}
		fi.Kernel = k
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
	fi.Display = displayInfo()
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
	fi.Shell = shellInfo()
	fi.DEWM = deWM()
	fi.PkgKind, fi.Packages = pkgInfo()
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

var pciIDRe = regexp.MustCompile(`\[([0-9a-fA-F]{4}):([0-9a-fA-F]{4})\]`)
var pciRevRe = regexp.MustCompile(`\(rev ([0-9a-fA-F]+)\)`)

// gpuModel достаёт название видеокарты через lspci -nn. Для AMD берёт
// маркетинговое имя из amdgpu.ids (как fastfetch) по PCI-ID и ревизии — иначе
// lspci даёт общее имя чипа «Radeon RX 470/480/570/…». "" — если lspci нет.
func gpuModel() string {
	if !cmdExists("lspci") {
		return ""
	}
	out, err := exec.Command("lspci", "-nn").Output()
	if err != nil {
		return ""
	}
	for _, line := range strings.Split(string(out), "\n") {
		if strings.Contains(line, "VGA compatible controller") ||
			strings.Contains(line, "3D controller") ||
			strings.Contains(line, "Display controller") {
			return resolveGPU(line)
		}
	}
	return ""
}

// resolveGPU превращает строку lspci -nn в имя карты. Для AMD (vendor 1002) ищет
// маркетинговое имя в amdgpu.ids по device-ID и ревизии; иначе чистит имя из lspci.
func resolveGPU(line string) string {
	ids := pciIDRe.FindStringSubmatch(line)
	if ids != nil && strings.EqualFold(ids[1], "1002") {
		rev := ""
		if m := pciRevRe.FindStringSubmatch(line); m != nil {
			rev = m[1]
		}
		if name := amdMarketing(ids[2], rev); name != "" {
			// В amdgpu.ids имя обычно уже с «AMD»; не дублируем.
			if !strings.HasPrefix(name, "AMD") {
				name = "AMD " + name
			}
			return name
		}
	}
	name := line
	if i := strings.Index(name, "]: "); i >= 0 {
		name = name[i+3:]
	}
	// Убираем "[1002:67df]" и "(rev ef)".
	name = pciIDRe.ReplaceAllString(name, "")
	name = pciRevRe.ReplaceAllString(name, "")
	return cleanGPU(name)
}

// amdMarketing ищет в amdgpu.ids (поставляется с libdrm/Mesa) маркетинговое имя
// карты по device-ID и ревизии — тот же источник, что у fastfetch.
func amdMarketing(devID, revID string) string {
	var path string
	for _, p := range []string{"/usr/share/libdrm/amdgpu.ids", "/usr/share/drm/amdgpu.ids"} {
		if _, err := os.Stat(p); err == nil {
			path = p
			break
		}
	}
	if path == "" {
		return ""
	}
	data, err := os.ReadFile(path)
	if err != nil {
		return ""
	}
	return amdLookup(string(data), devID, revID)
}

// amdLookup ищет имя карты в содержимом amdgpu.ids по device-ID и ревизии.
// Формат строки: «DEVID,<таб>REVID,<таб>Название» (плюс шапка-версия и #-коммент).
func amdLookup(data, devID, revID string) string {
	devID, revID = strings.ToUpper(devID), strings.ToUpper(revID)
	for _, line := range strings.Split(data, "\n") {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		parts := strings.SplitN(line, ",", 3) // DEVID, REVID, NAME
		if len(parts) < 3 {
			continue
		}
		if strings.ToUpper(strings.TrimSpace(parts[0])) == devID &&
			strings.ToUpper(strings.TrimSpace(parts[1])) == revID {
			return strings.TrimSpace(parts[2])
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

// displayInfo формирует строку дисплея вида «2560x1440, 27", 180Hz, 108ppi».
// Диагональ и ppi считаются из физического размера (мм) и разрешения; если
// физический размер неизвестен, эти части опускаются.
func displayInfo() string {
	pxW, pxH, hz, mmW, mmH := probeDisplay()
	if pxW == 0 || pxH == 0 {
		return ""
	}
	parts := []string{fmt.Sprintf("%dx%d", pxW, pxH)}
	diagIn := 0.0
	if mmW > 0 && mmH > 0 {
		diagIn = math.Hypot(float64(mmW), float64(mmH)) / 25.4
	}
	if diagIn > 0 {
		parts = append(parts, fmt.Sprintf("%.0f\"", diagIn))
	}
	if hz > 0 {
		parts = append(parts, fmt.Sprintf("%dHz", hz))
	}
	if diagIn > 0 {
		ppi := math.Hypot(float64(pxW), float64(pxH)) / diagIn
		parts = append(parts, fmt.Sprintf("%.0fppi", ppi))
	}
	return strings.Join(parts, ", ")
}

// probeDisplay возвращает разрешение (px), частоту (Гц) и физический размер (мм)
// активного дисплея: сначала через xrandr (X11/XWayland), затем wlr-randr.
func probeDisplay() (pxW, pxH, hz, mmW, mmH int) {
	if cmdExists("xrandr") {
		if out, err := exec.Command("xrandr", "--current").Output(); err == nil {
			s := string(out)
			// Строка подключённого выхода: "…2560x1440+0+0 (…) 597mm x 336mm".
			conn := regexp.MustCompile(`(\d+)x(\d+)\+\d+\+\d+.*?(\d+)mm x (\d+)mm`)
			if m := conn.FindStringSubmatch(s); m != nil {
				pxW, pxH = atoiSafe(m[1]), atoiSafe(m[2])
				mmW, mmH = atoiSafe(m[3]), atoiSafe(m[4])
			}
			// Активный режим помечен «*» — из него берём частоту.
			ref := regexp.MustCompile(`(\d+)x(\d+)\s+([\d.]+)\*`)
			if m := ref.FindStringSubmatch(s); m != nil {
				if pxW == 0 {
					pxW, pxH = atoiSafe(m[1]), atoiSafe(m[2])
				}
				hz = roundHz(m[3])
			}
			if pxW != 0 {
				return
			}
		}
	}
	if cmdExists("wlr-randr") {
		if out, err := exec.Command("wlr-randr").Output(); err == nil {
			s := string(out)
			if m := regexp.MustCompile(`(\d+)x(\d+) px, ([\d.]+) Hz.*current`).FindStringSubmatch(s); m != nil {
				pxW, pxH, hz = atoiSafe(m[1]), atoiSafe(m[2]), roundHz(m[3])
			}
			if m := regexp.MustCompile(`Physical size: (\d+)x(\d+) mm`).FindStringSubmatch(s); m != nil {
				mmW, mmH = atoiSafe(m[1]), atoiSafe(m[2])
			}
		}
	}
	return
}

func atoiSafe(s string) int { n, _ := strconv.Atoi(s); return n }

// roundHz округляет частоту "60.000000"/"59.94" до целого числа герц.
func roundHz(s string) int {
	if f, err := strconv.ParseFloat(s, 64); err == nil {
		return int(f + 0.5)
	}
	return 0
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

// pkgInfo определяет тип системного пакетного менеджера и число установленных
// пакетов (например "rpm", 2451). Пусто — если менеджер не найден.
func pkgInfo() (string, int) {
	managers := []struct {
		bin, kind string
		args      []string
	}{
		{"pacman", "pacman", []string{"-Qq"}},
		{"dpkg-query", "deb", []string{"-f", ".\n", "-W"}},
		{"rpm", "rpm", []string{"-qa"}},
		{"apk", "apk", []string{"info"}},
		{"xbps-query", "xbps", []string{"-l"}},
		{"qlist", "emerge", []string{"-I"}},
		{"eopkg", "eopkg", []string{"list-installed"}},
	}
	for _, mgr := range managers {
		if !cmdExists(mgr.bin) {
			continue
		}
		if out, err := exec.Command(mgr.bin, mgr.args...).Output(); err == nil {
			if n := countLines(string(out)); n > 0 {
				return mgr.kind, n
			}
		}
	}
	// Gentoo без portage-utils: считаем каталоги установленных пакетов.
	if n := countPortagePkgs(); n > 0 {
		return "emerge", n
	}
	return "", 0
}

// countPortagePkgs считает установленные пакеты Gentoo по /var/db/pkg/<cat>/<pkg>.
func countPortagePkgs() int {
	cats, err := os.ReadDir("/var/db/pkg")
	if err != nil {
		return 0
	}
	n := 0
	for _, c := range cats {
		if !c.IsDir() {
			continue
		}
		if pkgs, err := os.ReadDir("/var/db/pkg/" + c.Name()); err == nil {
			for _, p := range pkgs {
				if p.IsDir() {
					n++
				}
			}
		}
	}
	return n
}

var verRe = regexp.MustCompile(`\d+(?:\.\d+)+`)

// versionFrom запускает `bin args…` и вытаскивает первую версию вида X.Y[.Z].
func versionFrom(bin string, args ...string) string {
	if !cmdExists(bin) {
		return ""
	}
	out, err := exec.Command(bin, args...).CombinedOutput()
	if err != nil && len(out) == 0 {
		return ""
	}
	return verRe.FindString(string(out))
}

// shellInfo возвращает имя shell с версией: «zsh 5.9», «bash 5.2.21».
func shellInfo() string {
	sh := filepath.Base(os.Getenv("SHELL"))
	if sh == "" || sh == "." {
		return "—"
	}
	if v := versionFrom(sh, "--version"); v != "" {
		return sh + " " + v
	}
	return sh
}

// deWM возвращает окружение рабочего стола с версией: «GNOME 50.2», «KDE Plasma 6.1».
func deWM() string {
	de := firstNonEmpty(os.Getenv("XDG_CURRENT_DESKTOP"), os.Getenv("DESKTOP_SESSION"))
	if de == "" {
		if os.Getenv("XDG_SESSION_TYPE") == "tty" {
			return "TTY"
		}
		return "—"
	}
	// XDG_CURRENT_DESKTOP бывает составным: «ubuntu:GNOME» → «GNOME».
	if i := strings.LastIndex(de, ":"); i >= 0 {
		de = de[i+1:]
	}
	if v := deVersion(de); v != "" {
		return de + " " + v
	}
	return de
}

// deVersion пытается узнать версию окружения через его CLI (--version).
func deVersion(de string) string {
	switch low := strings.ToLower(de); {
	case strings.Contains(low, "gnome"):
		return versionFrom("gnome-shell", "--version")
	case strings.Contains(low, "kde") || strings.Contains(low, "plasma"):
		return versionFrom("plasmashell", "--version")
	case strings.Contains(low, "xfce"):
		return versionFrom("xfce4-session", "--version")
	case strings.Contains(low, "cinnamon"):
		return versionFrom("cinnamon", "--version")
	case strings.Contains(low, "mate"):
		return versionFrom("mate-session", "--version")
	case strings.Contains(low, "budgie"):
		return versionFrom("budgie-desktop", "--version")
	case strings.Contains(low, "lxqt"):
		return versionFrom("lxqt-session", "--version")
	case strings.Contains(low, "deepin"):
		return versionFrom("startdde", "--version")
	case strings.Contains(low, "hyprland"):
		return versionFrom("Hyprland", "--version")
	case strings.Contains(low, "sway"):
		return versionFrom("sway", "--version")
	case strings.Contains(low, "i3"):
		return versionFrom("i3", "--version")
	case strings.Contains(low, "river"):
		return versionFrom("river", "--version")
	}
	return ""
}

func fmtUptime(sec int) string {
	d := sec / 86400
	h := (sec % 86400) / 3600
	mnt := (sec % 3600) / 60
	var parts []string
	if d > 0 {
		parts = append(parts, fmt.Sprintf("%dd", d))
	}
	if h > 0 {
		parts = append(parts, fmt.Sprintf("%dh", h))
	}
	if mnt > 0 || len(parts) == 0 {
		parts = append(parts, fmt.Sprintf("%dm", mnt))
	}
	return strings.Join(parts, " ")
}

func humanGiB(kb int) string {
	return fmt.Sprintf("%.1f GiB", float64(kb)/1048576)
}

// plural выбирает единственное/множественное число (английский).
func plural(n int, one, many string) string {
	if n == 1 {
		return one
	}
	return many
}

func humanAge(days int) string {
	if days >= 365 {
		y, d := days/365, days%365
		yw := fmt.Sprintf("%d %s", y, plural(y, "year", "years"))
		if d == 0 {
			return yw
		}
		return fmt.Sprintf("%s %d %s", yw, d, plural(d, "day", "days"))
	}
	return fmt.Sprintf("%d %s", days, plural(days, "day", "days"))
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
//	╭─ OS ────────────────╮
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
	kv := func(l, v string) string { return kvRow(l, v, 7) }

	osRows := []string{kv("Distro", fi.Distro), kv("Kernel", fi.Kernel)}
	if fi.Uptime != "" {
		osRows = append(osRows, kv("Uptime", fi.Uptime))
	}
	if fi.HasInstall {
		osRows = append(osRows,
			kv("Birth", fi.InstallDate),
			kv("Age", humanAge(fi.AgeDays)))
	}

	var hwRows []string
	if fi.CPU != "" {
		hwRows = append(hwRows, kv("CPU", fmt.Sprintf("%s (%d %s)",
			fi.CPU, fi.Cores, plural(fi.Cores, "core", "cores"))))
	}
	if fi.GPU != "" {
		hwRows = append(hwRows, kv("GPU", fi.GPU))
	}
	if fi.Display != "" {
		hwRows = append(hwRows, kv("Disp", fi.Display))
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
	// Пакеты: тип системного менеджера + количество, затем flatpak/snap.
	var pkgs []string
	if fi.Packages > 0 {
		kind := fi.PkgKind
		if kind == "" {
			kind = "pkgs"
		}
		pkgs = append(pkgs, fmt.Sprintf("%s %d", kind, fi.Packages))
	}
	if fi.FlatpakN > 0 {
		pkgs = append(pkgs, fmt.Sprintf("flatpak %d", fi.FlatpakN))
	}
	if fi.SnapN > 0 {
		pkgs = append(pkgs, fmt.Sprintf("snap %d", fi.SnapN))
	}
	if len(pkgs) > 0 {
		swRows = append(swRows, kv("Pkgs", strings.Join(pkgs, " · ")))
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
		info = append(info, kv("RAM", fmt.Sprintf("%s  %s of %s",
			usageBar(fi.RAMPct, 10), fi.RAMUsed, fi.RAMTotal)))
	}
	if fi.HasDisk {
		info = append(info, kv("Disk", fmt.Sprintf("%s  %s of %s",
			usageBar(fi.DiskPct, 10), fi.DiskUsed, fi.DiskTotal)))
	}
	return info
}

func renderFetch(m Model) string {
	fi := computeFetch()

	key := logoKey(fi.DistroID)
	// В кратком виде — компактный логотип; в подробном инфы много и окно
	// большое, поэтому уместен крупный логотип.
	logo := distroLogo(fi.DistroID, m.fetchFull)
	// Видимая ширина считается без кодов цвета $N.
	logoW := 0
	for _, l := range logo {
		if w := lipgloss.Width(stripLogoCodes(l)); w > logoW {
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
		// Раскраска по кодам $N + добивка пробелами до ширины колонки логотипа.
		pad := strings.Repeat(" ", logoW-lipgloss.Width(stripLogoCodes(l)))
		sb.WriteString("  " + colorLogoLine(l, key) + pad + "   " + r + "\n")
	}
	return sb.String()
}
