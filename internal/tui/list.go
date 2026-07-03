package tui

import (
	"fmt"
	"strings"

	"github.com/Axawys/lxprofiler/internal/data"
)

func renderList(m Model) string {
	var sb strings.Builder

	sb.WriteString(boldStyle.Render("  🐧 Профиль архетипов"))
	sb.WriteString("\n\n")

	maxLen := 0
	for _, r := range m.results {
		if len([]rune(r.Label)) > maxLen {
			maxLen = len([]rune(r.Label))
		}
	}

	for i, r := range m.results {
		label := r.Label
		pct := r.NormScore

		if data.Mystery[r.Key] && i != m.selected {
			label = maskLabel(label)
			pct = -1
		}

		bar := makeBar(pct, 20)
		pad := maxLen - len([]rune(label))

		var pf string
		if pct < 0 {
			pf = "???"
		} else {
			pf = fmt.Sprintf("%3d", pct)
		}

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
	sb.WriteString(dimStyle.Render("  ────────────────────────────────────────────"))
	sb.WriteString("\n")

	if m.selected < len(m.results) {
		r := m.results[m.selected]
		sb.WriteString(boldStyle.Render(fmt.Sprintf("▶ %s — %d%%", r.Label, r.NormScore)))
		sb.WriteString("\n")

		desc := data.Describe(r.Key, r.NormScore)
		sb.WriteString("  " + wrapText(desc, m.width-4))
		sb.WriteString("\n\n")

		sb.WriteString(boldStyle.Render("  Что повлияло:"))
		sb.WriteString("\n")
		sb.WriteString("  " + dimStyle.Render(wrapText(r.Reason, m.width-4)))
		sb.WriteString("\n")
	}

	sb.WriteString("\n")
	sb.WriteString(dimStyle.Render("  ↑↓ — листать · ←→ — режим · q — выход"))

	return sb.String()
}
