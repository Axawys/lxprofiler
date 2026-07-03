package tui

import (
	"fmt"
	"sort"
	"strings"

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
	var sb strings.Builder
	rows := m.orderedResults()

	maxLen := 0
	for _, r := range rows {
		if len([]rune(r.Label)) > maxLen {
			maxLen = len([]rune(r.Label))
		}
	}

	// Полоска тянется до правого края: ширина окна минус префикс строки
	// ("▶ " + метка + добивка + "  100%  "). Минимум 20, чтобы не схлопывалась.
	barW := m.width - (maxLen + 10)
	if barW < 20 {
		barW = 20
	}

	// Прогресс анимации: полоски и проценты заполняются от 0 до финала.
	prog := 1.0
	if m.animating {
		prog = m.animProgress
	}

	for i, r := range rows {
		label := r.Label
		pct := r.NormScore
		selected := !m.animating && i == m.selected

		if data.Mystery[r.Key] && !selected {
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

		var line string
		if selected {
			line = greenStyle.Render(fmt.Sprintf("▶ %s%*s  %s%%  %s", label, pad, "", pf, bar))
		} else {
			line = dimStyle.Render(fmt.Sprintf("  %s%*s  %s%%  %s", label, pad, "", pf, bar))
		}
		sb.WriteString(line)
		sb.WriteString("\n")
	}

	// Во время анимации — только гонка полосок и подсказка (панель деталей
	// прыгала бы при переупорядочивании).
	if m.animating {
		sb.WriteString("\n")
		sb.WriteString(dimStyle.Render("  q — выход"))
		return sb.String()
	}

	sb.WriteString("\n")
	sb.WriteString(dimStyle.Render(strings.Repeat("─", m.width)))
	sb.WriteString("\n")

	if m.selected < len(rows) {
		r := rows[m.selected]
		sb.WriteString(boldStyle.Render(fmt.Sprintf("▶ %s — %d%%", r.Label, r.NormScore)))
		sb.WriteString("\n")

		desc := data.Describe(r.Key, r.NormScore)
		sb.WriteString("  " + wrapText(desc, m.width-4))
		sb.WriteString("\n")

		sb.WriteString(boldStyle.Render("  Что повлияло:"))
		sb.WriteString("\n")
		sb.WriteString("  " + dimStyle.Render(wrapText(r.Reason, m.width-4)))
		sb.WriteString("\n")
	}

	sb.WriteString("\n")
	sb.WriteString(dimStyle.Render("  ↑↓ — листать · ←→ — режим · q — выход"))

	return sb.String()
}
