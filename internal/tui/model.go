package tui

import (
	"fmt"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/Axawys/lxprofiler/internal/detect"
)

// animDuration — длительность анимации заполнения полосок при запуске.
const animDuration = 3 * time.Second

// animTickMsg — тик анимации (~30 кадров/с).
type animTickMsg time.Time

func animTick() tea.Cmd {
	return tea.Tick(time.Second/30, func(t time.Time) tea.Msg {
		return animTickMsg(t)
	})
}

var (
	boldStyle  = lipgloss.NewStyle().Bold(true)
	dimStyle   = lipgloss.NewStyle().Faint(true)
	greenStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("2"))
	cyanStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("6"))
	redStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("1"))
	yellowStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("3"))
)

type Mode int

const (
	ListMode Mode = iota
	CompassMode
	StatsMode
	FetchMode
)

// modeCount — число режимов для циклического переключения ←/→.
const modeCount = 4

type Model struct {
	selected int
	mode     Mode
	results  []detect.ArchetypeResult
	width    int
	height   int
	reqW     int // минимальная ширина терминала, чтобы всё поместилось
	reqH     int // минимальная высота терминала

	animating    bool      // идёт ли анимация заполнения полосок
	animStart    time.Time // момент первого тика (для расчёта прогресса)
	animProgress float64   // 0..1 — доля заполнения полосок
}

func NewModel(results []detect.ArchetypeResult, animate bool) Model {
	m := Model{selected: 0, results: results, animating: animate}
	m.reqW, m.reqH = computeRequiredSize(results)
	return m
}

// computeRequiredSize считает минимальный размер терминала, при котором ни один
// из режимов (список / координаты / статистика) не обрезается. Считается один
// раз при создании модели: набор классов фиксирован, так что размер не меняется.
func computeRequiredSize(results []detect.ArchetypeResult) (int, int) {
	// Ширина строки списка: "▶ " + метка + добивка + "  100%  " + бар(20).
	maxLabel := 0
	for _, r := range results {
		if w := lipgloss.Width(r.Label); w > maxLabel {
			maxLabel = w
		}
	}
	reqW := 2 + maxLabel + 2 + 3 + 1 + 2 + 20
	if reqW < 48 { // пол для разделителя/шапки
		reqW = 48
	}

	// Координаты и статистика по ширине не переносятся — берём их натуральную
	// ширину, отрисовав на заведомо большом «холсте». Заголовок-линейка
	// (titleRule) тянется во всю ширину холста — она адаптивна и не должна
	// задавать минимально требуемую ширину, поэтому строки длиной с холст
	// (== canvas) при измерении пропускаем.
	const canvas = 400
	big := Model{results: results, width: canvas, height: canvas}
	big.mode = CompassMode
	compassView := renderCompass(big)
	big.mode = StatsMode
	statsView := renderStats(big)
	big.mode = FetchMode
	fetchView := renderFetch(big)
	for _, v := range []string{compassView, statsView, fetchView} {
		if w := maxContentWidth(v, canvas); w > reqW {
			reqW = w
		}
	}

	// Высота: максимум по режимам. У списка высота зависит от переноса описания
	// и «что повлияло», а те переносятся по ширине reqW — берём худший класс.
	reqH := lineCount(compassView)
	for _, v := range []string{statsView, fetchView} {
		if h := lineCount(v); h > reqH {
			reqH = h
		}
	}
	lm := Model{results: results, mode: ListMode, width: reqW, height: 400}
	if len(results) == 0 {
		if h := lineCount(renderList(lm)); h > reqH {
			reqH = h
		}
	}
	for sel := range results {
		lm.selected = sel
		if h := lineCount(renderList(lm)); h > reqH {
			reqH = h
		}
	}
	return reqW, reqH
}

func maxLineWidth(s string) int {
	max := 0
	for _, line := range strings.Split(s, "\n") {
		if w := lipgloss.Width(line); w > max {
			max = w
		}
	}
	return max
}

// maxContentWidth — как maxLineWidth, но пропускает строки во всю ширину холста
// (адаптивные заголовки-линейки), чтобы они не завышали требуемую ширину.
func maxContentWidth(s string, canvas int) int {
	max := 0
	for _, line := range strings.Split(s, "\n") {
		w := lipgloss.Width(line)
		if w >= canvas {
			continue
		}
		if w > max {
			max = w
		}
	}
	return max
}

func lineCount(s string) int { return strings.Count(s, "\n") + 1 }

func (m Model) Init() tea.Cmd {
	if m.animating {
		return animTick()
	}
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case animTickMsg:
		if !m.animating {
			return m, nil
		}
		if m.animStart.IsZero() {
			m.animStart = time.Time(msg)
		}
		m.animProgress = float64(time.Time(msg).Sub(m.animStart)) / float64(animDuration)
		if m.animProgress >= 1 {
			m.animProgress = 1
			m.animating = false
			return m, nil
		}
		return m, animTick()
	case tea.KeyMsg:
		key := msg.String()
		// Пока идёт анимация, доступен только выход по q — всё остальное
		// (навигация, смена режима, промотка) недоступно до её конца.
		if m.animating {
			switch key {
			case "q", "Q", "й", "Й":
				return m, tea.Quit
			}
			return m, nil
		}
		switch key {
		case "q", "Q", "й", "Й":
			return m, tea.Quit
		case "j", "о", "down":
			if m.selected < len(m.results)-1 {
				m.selected++
			}
		case "k", "л", "up":
			if m.selected > 0 {
				m.selected--
			}
		case "g", "п":
			m.selected = 0
		case "G", "П":
			m.selected = len(m.results) - 1
		case "m", "M", "ь", "Ь", "l", "д", "right":
			m.mode = (m.mode + 1) % modeCount
		case "h", "р", "left":
			m.mode = (m.mode + modeCount - 1) % modeCount
		}
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
	}
	return m, nil
}

func (m Model) View() string {
	// Пока не пришёл первый WindowSizeMsg — не мигаем сообщением «0×0».
	if m.width == 0 && m.height == 0 {
		return ""
	}
	if m.width < m.reqW || m.height < m.reqH {
		return tooSmallView(m)
	}
	switch m.mode {
	case CompassMode:
		return renderCompass(m)
	case StatsMode:
		return renderStats(m)
	case FetchMode:
		return renderFetch(m)
	default:
		return renderList(m)
	}
}

// tooSmallView показывается, когда терминал меньше требуемого. При ресайзе окна
// bubbletea перерисует View, и как только размер дорастёт — покажется утилита.
func tooSmallView(m Model) string {
	curW := greenStyle
	if m.width < m.reqW {
		curW = redStyle
	}
	curH := greenStyle
	if m.height < m.reqH {
		curH = redStyle
	}
	cur := curW.Render(fmt.Sprintf("%d", m.width)) +
		dimStyle.Render("×") + curH.Render(fmt.Sprintf("%d", m.height))

	lines := []string{
		boldStyle.Render("🐧 Окно слишком маленькое"),
		"",
		fmt.Sprintf("Нужно: %s", greenStyle.Render(fmt.Sprintf("%d×%d", m.reqW, m.reqH))),
		fmt.Sprintf("Сейчас: %s", cur),
		"",
		dimStyle.Render("Увеличьте окно · q — выход"),
	}
	block := strings.Join(lines, "\n")
	return lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center, block)
}

func makeBar(pct int, width int) string {
	if pct < 0 {
		return makeBrokenBar(width)
	}
	filled := pct * width / 100
	return strings.Repeat("█", filled) + strings.Repeat("░", width-filled)
}

// makeBrokenBar рисует «сломанную» полоску ширины width рун (не байт!).
// Прежняя версия резала многобайтовые руны по байтам (pattern[:width-1]),
// из-за чего бар секретных классов ломался в кракозябры.
func makeBrokenBar(width int) string {
	if width <= 0 {
		return ""
	}
	glyphs := []rune("█▒ █░▓ ░▒█ ░ ▓█░ ▓")
	var b strings.Builder
	for i := 0; i < width; i++ {
		b.WriteRune(glyphs[i%len(glyphs)])
	}
	return b.String()
}

// growBrokenBar рисует «сломанную» полоску, заполненную слева на долю prog
// (для анимации секретных классов); остаток — пробелы.
func growBrokenBar(width int, prog float64) string {
	if width <= 0 {
		return ""
	}
	n := int(float64(width) * prog)
	if n < 0 {
		n = 0
	}
	if n > width {
		n = width
	}
	full := []rune(makeBrokenBar(width))
	return string(full[:n]) + strings.Repeat(" ", width-n)
}

func maskLabel(label string) string {
	return strings.Repeat("?", len([]rune(label)))
}

// padRight добивает строку пробелами справа до видимой ширины w (по lipgloss.Width,
// чтобы кириллица и emoji считались корректно). Уже длинную строку не трогает.
func padRight(s string, w int) string {
	if d := w - lipgloss.Width(s); d > 0 {
		return s + strings.Repeat(" ", d)
	}
	return s
}

// padLeft добивает строку пробелами слева до видимой ширины w.
func padLeft(s string, w int) string {
	if d := w - lipgloss.Width(s); d > 0 {
		return strings.Repeat(" ", d) + s
	}
	return s
}

// titleRule рисует заголовок секции и добивает строку линией до правого края:
//
//	  📊 Заголовок ───────────────────────────────
func titleRule(title string, width int) string {
	head := "  " + title + " "
	rem := width - lipgloss.Width(head)
	if rem < 1 {
		rem = 1
	}
	return "  " + boldStyle.Render(title) + " " + dimStyle.Render(strings.Repeat("─", rem))
}

func wrapText(text string, width int) string {
	if width <= 0 {
		width = 72
	}
	words := strings.Fields(text)
	var lines []string
	current := ""
	for _, word := range words {
		// Ширина считается по видимым символам (lipgloss.Width), а не байтам —
		// иначе кириллица (2 байта/символ) переносится вдвое раньше, текст не
		// доходит до правой рамки и занимает лишние строки.
		switch {
		case current == "":
			current = word
		case lipgloss.Width(current)+1+lipgloss.Width(word) > width:
			lines = append(lines, current)
			current = word
		default:
			current += " " + word
		}
	}
	if current != "" {
		lines = append(lines, current)
	}
	return strings.Join(lines, "\n")
}
