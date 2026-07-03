package tui

import (
	"fmt"
	"strings"

	"github.com/Axawys/lxprofiler/internal/data"
)

func renderList(m Model) string {
	var sb strings.Builder

	maxLen := 0
	for _, r := range m.results {
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

	for i, r := range m.results {
		label := r.Label
		pct := r.NormScore

		if data.Mystery[r.Key] && i != m.selected {
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
		if i == m.selected {
			line = greenStyle.Render(fmt.Sprintf("▶ %s%*s  %s%%  %s", label, pad, "", pf, bar))
		} else {
			line = dimStyle.Render(fmt.Sprintf("  %s%*s  %s%%  %s", label, pad, "", pf, bar))
		}
		sb.WriteString(line)
		sb.WriteString("\n")
	}

	sb.WriteString("\n")
	sb.WriteString(dimStyle.Render(strings.Repeat("─", m.width)))
	sb.WriteString("\n")

	if m.selected < len(m.results) {
		r := m.results[m.selected]
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
