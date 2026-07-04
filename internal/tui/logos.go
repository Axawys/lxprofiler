package tui

import (
	"embed"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// Логотипы дистрибутивов лежат в каталоге logos/ в формате <distro>.txt
// (крупный, для краткого вида) и <distro>_small.txt (компактный, для
// подробного вида). Ключ <distro> — ID из os-release (см. logoKey).
// generic.txt / generic_small.txt — запасной вариант для неизвестных систем.
//
//go:embed logos
var logoFS embed.FS

func splitLogo(s string) []string {
	return strings.Split(strings.Trim(s, "\n"), "\n")
}

// loadLogo читает logos/<name>.txt из встроенной ФС; ok=false, если файла нет.
func loadLogo(name string) ([]string, bool) {
	b, err := logoFS.ReadFile("logos/" + name + ".txt")
	if err != nil {
		return nil, false
	}
	return splitLogo(string(b)), true
}

// distroLogo возвращает логотип по ID из os-release: крупный (big=true) или
// компактный. Сначала пробуется точный вариант дистрибутива (opensuse-tumbleweed
// → opensuse_tumbleweed), затем сведённый ключ (opensuse), иначе — generic.
func distroLogo(id string, big bool) []string {
	suffix := "_small"
	if big {
		suffix = ""
	}
	for _, key := range []string{fileKey(id), logoKey(id)} {
		if logo, ok := loadLogo(key + suffix); ok {
			return logo
		}
	}
	logo, _ := loadLogo("generic" + suffix)
	return logo
}

// fileKey переводит ID из os-release в имя файла логотипа: дефисы в ID
// (opensuse-tumbleweed) соответствуют подчёркиваниям в имени файла.
func fileKey(id string) string {
	return strings.ReplaceAll(id, "-", "_")
}

// logoKey приводит ID из os-release к ключу логотипа: openSUSE выпускается как
// opensuse-leap / opensuse-tumbleweed — сводим к общему «opensuse».
func logoKey(id string) string {
	if strings.HasPrefix(id, "opensuse") {
		return "opensuse"
	}
	return id
}

// logoPalettes — цветовая схема каждого дистрибутива под его официальный бренд.
// В тексте логотипа (формат fastfetch) код $N переключает цвет на N-й элемент
// схемы: $1 — первый цвет, $2 — второй и т.д. Текст до первого кода красится
// первым цветом. Цвета — truecolor hex бренда.
var logoPalettes = map[string][]lipgloss.Style{
	// Fedora: синий контур + белая «f».
	"fedora":    palette("#3C6EB4", "#FFFFFF"),
	"arch":      palette("#1793D1", "#0F94D2"),
	"ubuntu":    palette("#E95420", "#FFFFFF", "#772953"),
	"kubuntu":   palette("#0F7CDD", "#FFFFFF"),
	"debian":    palette("#D70A53", "#FFFFFF"),
	"linuxmint": palette("#87CF3E", "#FFFFFF"),
	"opensuse":  palette("#73BA25", "#FFFFFF"),
	"alpine":    palette("#0D597F", "#FFFFFF"),
	"gentoo":    palette("#9C71B9", "#62548E"),
	"nixos":     palette("#7EBAE4", "#5277C3"),
	"calculate": palette("#00A2A2", "#FFFFFF"),
	"altlinux":  palette("#E00000", "#FFFFFF"),
}

// palette собирает набор стилей из hex-цветов.
func palette(hexes ...string) []lipgloss.Style {
	styles := make([]lipgloss.Style, len(hexes))
	for i, h := range hexes {
		styles[i] = lipgloss.NewStyle().Foreground(lipgloss.Color(h))
	}
	return styles
}

// logoBase — базовый цвет строки логотипа: первый цвет схемы дистрибутива
// (текст до первого кода $N), либо тусклый, если схемы нет.
func logoBase(key string) lipgloss.Style {
	if p := logoPalettes[key]; len(p) > 0 {
		return p[0]
	}
	return dimStyle
}

// stripLogoCodes убирает коды цвета $N — для расчёта видимой ширины строки.
func stripLogoCodes(line string) string {
	var sb strings.Builder
	for i := 0; i < len(line); i++ {
		if line[i] == '$' && i+1 < len(line) && line[i+1] >= '0' && line[i+1] <= '9' {
			i++ // пропустить цифру
			continue
		}
		sb.WriteByte(line[i])
	}
	return sb.String()
}

// colorLogoLine раскрашивает строку логотипа по схеме key: код $N переключает
// цвет на N-й элемент схемы, остальной текст красится текущим цветом.
func colorLogoLine(line string, key string) string {
	pal := logoPalettes[key]
	base := logoBase(key)
	cur := base
	var out, seg strings.Builder
	flush := func() {
		if seg.Len() > 0 {
			out.WriteString(cur.Render(seg.String()))
			seg.Reset()
		}
	}
	for i := 0; i < len(line); i++ {
		if line[i] == '$' && i+1 < len(line) && line[i+1] >= '0' && line[i+1] <= '9' {
			flush()
			n := int(line[i+1] - '0')
			if n >= 1 && n <= len(pal) {
				cur = pal[n-1]
			} else {
				cur = base
			}
			i++
			continue
		}
		seg.WriteByte(line[i])
	}
	flush()
	return out.String()
}
