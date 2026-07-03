package tui

import (
	"fmt"
	"strings"

	"github.com/Axawys/lxprofiler/internal/detect"
)

var VX = map[string]int{
	"devops": -20, "programmer": -10, "sysadmin": 50, "minimalist": 0,
	"old_hacker": 70, "ricer": -30, "gamer": -10, "anonymous": 0,
	"pentester": 10, "import_substituted": 50, "fresh_witness": -90, "atomic": -60,
	"creative": -25,
}

var VY = map[string]int{
	"devops": -20, "programmer": 30, "sysadmin": 40, "minimalist": 60,
	"old_hacker": 70, "ricer": 70, "gamer": -60, "anonymous": 40,
	"pentester": 50, "import_substituted": -10, "fresh_witness": -10, "atomic": -50,
	"creative": -15,
}

type CompassResult struct {
	CX, CY   int
	Quadrant string
}

func ComputeCompass() CompassResult {
	sumX, sumY, sumW := 0, 0, 0
	for key, w := range detect.Score {
		if w <= 0 {
			continue
		}
		vx, ok := VX[key]
		if !ok {
			continue
		}
		vy := VY[key]
		sumX += w * vx
		sumY += w * vy
		sumW += w
	}
	if sumW == 0 {
		sumW = 1
	}
	cx := sumX / sumW
	cy := sumY / sumW
	return CompassResult{CX: cx, CY: cy, Quadrant: getQuadrant(cx, cy)}
}

func getQuadrant(cx, cy int) string {
	h, v := "C", "C"
	if cx <= -15 {
		h = "N"
	} else if cx >= 15 {
		h = "T"
	}
	if cy >= 15 {
		v = "U"
	} else if cy <= -15 {
		v = "D"
	}
	switch h + v {
	case "NU":
		return "Лаборатория — DIY-новатор (Arch/NixOS/tiling)"
	case "TU":
		return "Цитадель Unix — всё руками, старая школа"
	case "ND":
		return "Гладкое будущее — новое и из коробки (atomic/Bazzite)"
	case "TD":
		return "Тёплая гавань — стабильно и удобно (Ubuntu/Mint)"
	case "CU":
		return "Инженер-середняк — по взглядам центрист, но всё руками"
	case "CD":
		return "Прагматик — посередине по взглядам, ценит удобство"
	case "NC":
		return "Новатор-центрист — за свежее, баланс DIY и удобства"
	case "TC":
		return "Традиционалист-центрист — проверенное, баланс DIY и удобства"
	default:
		return "Центрист — сбалансированный линуксоид"
	}
}

func renderCompass(m Model) string {
	var sb strings.Builder
	compass := ComputeCompass()

	gw := 49
	gh := 15
	if gw > m.width-4 {
		gw = m.width - 4
	}
	if gw < 21 {
		gw = 21
	}
	if gw%2 == 0 {
		gw--
	}
	if gh > m.height-14 {
		gh = m.height - 14
	}
	if gh < 7 {
		gh = 7
	}
	if gh%2 == 0 {
		gh--
	}

	ccol := (gw - 1) / 2
	crow := (gh - 1) / 2

	ac := (compass.CX + 100) * (gw - 1) / 200
	ar := (100 - compass.CY) * (gh - 1) / 200
	if ac < 0 {
		ac = 0
	}
	if ac > gw-1 {
		ac = gw - 1
	}
	if ar < 0 {
		ar = 0
	}
	if ar > gh-1 {
		ar = gh - 1
	}

	sx := fmt.Sprintf("%+d", compass.CX)
	sy := fmt.Sprintf("%+d", compass.CY)
	coordLabel := cyanStyle.Render(fmt.Sprintf("%s;%s", sx, sy))

	sb.WriteString("   " + boldStyle.Render("▲ КОНТРОЛЬ (всё руками)"))
	sb.WriteString("\n")

	for r := 0; r < gh; r++ {
		line := "   "
		for c := 0; c < gw; c++ {
			var cell string
			if r == ar && c == ac {
				cell = cyanStyle.Render("●")
			} else if r == crow && c == ccol {
				cell = dimStyle.Render("┼")
			} else if r == crow {
				cell = dimStyle.Render("─")
			} else if c == ccol {
				cell = dimStyle.Render("│")
			} else {
				cell = " "
			}
			line += cell
		}
		if r == ar {
			line += "  " + coordLabel
		}
		sb.WriteString(line)
		sb.WriteString("\n")
	}

	sb.WriteString("   " + boldStyle.Render("▼ УДОБСТВО (из коробки)"))
	sb.WriteString("\n")
	sb.WriteString("   " + dimStyle.Render(fmt.Sprintf("◄ новаторы%*sтрадиции ►", gw-20, "")))
	sb.WriteString("\n\n")

	sb.WriteString(fmt.Sprintf("  %s %s %s", cyanStyle.Render("●"), boldStyle.Render("ты:"), compass.Quadrant))
	sb.WriteString("\n")
	sb.WriteString(fmt.Sprintf("  %s", dimStyle.Render(fmt.Sprintf("   координаты (новат↔трад ; контроль↔удоб): %s;%s", sx, sy))))
	sb.WriteString("\n\n")
	sb.WriteString(dimStyle.Render("  ↑↓ — листать · ←→ — режим · q — выход"))

	return sb.String()
}
