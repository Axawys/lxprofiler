package tui

import (
	"fmt"
	"sort"
	"strings"

	"github.com/charmbracelet/lipgloss"

	"github.com/Axawys/lxprofiler/internal/data"
	"github.com/Axawys/lxprofiler/internal/detect"
)

// smoothstep — плавная кривая 0→1 (медленно на краях, быстрее в середине).
func smoothstep(x float64) float64 {
	if x <= 0 {
		return 0
	}
	if x >= 1 {
		return 1
	}
	return x * x * (3 - 2*x)
}

// orderedResults возвращает классы в текущем порядке отрисовки. Вне анимации —
// финальный порядок (по очкам). Во время анимации это «гонка»: старт всегда в
// одном и том же порядке (по алфавиту), а по мере заполнения классы с бо́льшим
// счётом поднимаются к своим финальным местам. Таинственные классы всегда внизу.
func (m Model) orderedResults() []detect.ArchetypeResult {
	if !m.animating {
		return m.results
	}

	var normal, mystery []detect.ArchetypeResult
	for _, r := range m.results {
		if data.Mystery[r.Key] {
			mystery = append(mystery, r)
		} else {
			normal = append(normal, r)
		}
	}

	// Финальные позиции — m.results уже отсортирован по очкам.
	finalIdx := make(map[string]int, len(normal))
	for i, r := range normal {
		finalIdx[r.Key] = i
	}
	// Стартовые позиции — фиксированный порядок (по метке), одинаковый каждый раз.
	fixed := make([]detect.ArchetypeResult, len(normal))
	copy(fixed, normal)
	sort.Slice(fixed, func(i, j int) bool { return fixed[i].Label < fixed[j].Label })
	startIdx := make(map[string]int, len(fixed))
	for i, r := range fixed {
		startIdx[r.Key] = i
	}

	// Интерполяция позиции старт→финал по прогрессу.
	e := smoothstep(m.animProgress)
	pos := func(r detect.ArchetypeResult) float64 {
		return float64(startIdx[r.Key])*(1-e) + float64(finalIdx[r.Key])*e
	}
	ordered := make([]detect.ArchetypeResult, len(normal))
	copy(ordered, normal)
	sort.SliceStable(ordered, func(i, j int) bool { return pos(ordered[i]) < pos(ordered[j]) })

	return append(ordered, mystery...)
}

func renderList(m Model) string {
	// Во время анимации — только «гонка» полосок без рамок и панелей (они
	// прыгали бы при переупорядочивании классов).
	if m.animating {
		return renderRacing(m)
	}
	return pinTabBar(listBody(m), m)
}

// listBody собирает основной режим: панель со списком архетипов и панель с
// деталями выбранного. Панель вкладок добавляет pinTabBar.
func listBody(m Model) []string {
	rows := m.orderedResults()
	maxLen := 0
	for _, r := range rows {
		if l := len([]rune(r.Label)); l > maxLen {
			maxLen = l
		}
	}
	inner := m.width - 6 // 2 отступ панели + рамка «│ » и « │»
	// Ширина полоски: остаток строки после «▶ » + метка + добивка + «  100%  ».
	barW := inner - (maxLen + 10)
	if barW < 12 {
		barW = 12
	}

	bars := make([]string, len(rows))
	for i, r := range rows {
		label := r.Label
		pct := r.NormScore
		selected := i == m.selected
		if data.Mystery[r.Key] && !selected {
			label = maskLabel(label)
			pct = -1
		}
		var bar, pf string
		if pct < 0 {
			bar = makeBrokenBar(barW)
			pf = "???"
		} else {
			bar = makeBar(pct, barW)
			pf = fmt.Sprintf("%3d", pct)
		}
		pad := maxLen - len([]rune(label))
		marker := "  "
		if selected {
			marker = "▶ "
		}
		text := fmt.Sprintf("%s%s%*s  %s%%  %s", marker, label, pad, "", pf, bar)
		if selected {
			bars[i] = greenStyle.Render(text)
		} else {
			bars[i] = dimStyle.Render(text)
		}
	}
	content := panel("Архетипы", bars, inner)

	if m.selected < len(rows) {
		r := rows[m.selected]
		var d []string
		d = append(d, wrapLines(data.Describe(r.Key, r.NormScore), inner, lipgloss.NewStyle())...)
		d = append(d, "")
		d = append(d, boldStyle.Render("Что повлияло:"))
		d = append(d, wrapLines(r.Reason, inner, dimStyle)...)
		content = append(content, panel(fmt.Sprintf("%s — %d%%", r.Label, r.NormScore), d, inner)...)
	}
	return content
}

// renderRacing — экран анимации: полоски классов «догоняют» свои значения.
func renderRacing(m Model) string {
	var sb strings.Builder
	rows := m.orderedResults()
	maxLen := 0
	for _, r := range rows {
		if l := len([]rune(r.Label)); l > maxLen {
			maxLen = l
		}
	}
	barW := m.width - (maxLen + 10)
	if barW < 20 {
		barW = 20
	}
	prog := m.animProgress
	for _, r := range rows {
		label := r.Label
		pct := r.NormScore
		if data.Mystery[r.Key] {
			label = maskLabel(label)
			pct = -1
		}
		var bar, pf string
		if pct < 0 {
			bar = growBrokenBar(barW, prog)
			pf = "???"
		} else {
			dp := int(float64(pct) * prog)
			bar = makeBar(dp, barW)
			pf = fmt.Sprintf("%3d", dp)
		}
		pad := maxLen - len([]rune(label))
		sb.WriteString(dimStyle.Render(fmt.Sprintf("  %s%*s  %s%%  %s", label, pad, "", pf, bar)))
		sb.WriteString("\n")
	}
	sb.WriteString("\n")
	sb.WriteString(dimStyle.Render("  q — выход"))
	return sb.String()
}
