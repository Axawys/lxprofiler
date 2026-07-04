package tui

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"sync"

	"github.com/charmbracelet/lipgloss"
)

// fetchInfo — данные для режима «суперфетч» (аналог fastfetch, но уже
// сконфигурированный). Пока минимальный набор: дистрибутив, ядро, диск.
type fetchInfo struct {
	User, Host string
	Distro     string // PRETTY_NAME из os-release
	DistroID   string // ID из os-release — выбор логотипа
	Kernel     string
	DiskPct    int
	DiskUsed   string
	DiskTotal  string
	HasDisk    bool
}

var (
	fetchOnce   sync.Once
	fetchCached fetchInfo
)

// computeFetch собирает системную инфу один раз (df вызывается не на каждый кадр).
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
	return fi
}

func trimOSValue(s string) string {
	return strings.Trim(strings.TrimSpace(s), `"`)
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

func fetchRow(label, value string) string {
	return cyanStyle.Render(padRight(label, 7)) + value
}

// diskBar — полоска заполнения диска: закрашенная часть по проценту (зелёный/
// жёлтый/красный), остаток — тусклый.
func diskBar(pct, width int) string {
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

	title := boldStyle.Render(fi.User) + dimStyle.Render("@") + boldStyle.Render(fi.Host)
	info := []string{
		title,
		dimStyle.Render(strings.Repeat("─", lipgloss.Width(fi.User)+1+lipgloss.Width(fi.Host))),
		fetchRow("OS", fi.Distro),
		fetchRow("Kernel", fi.Kernel),
	}
	if fi.HasDisk {
		info = append(info, fetchRow("Disk", fmt.Sprintf("%s %d%%  %s / %s",
			diskBar(fi.DiskPct, 10), fi.DiskPct, fi.DiskUsed, fi.DiskTotal)))
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
	sb.WriteString(dimStyle.Render("  ↑↓ — листать · ←→ — режим · q — выход"))
	return sb.String()
}
