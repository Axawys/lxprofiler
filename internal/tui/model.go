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
	StatsMode
	FetchMode
)

// modeCount — число режимов для циклического переключения ←/→.
const modeCount = 3

// tabBar — нижняя строка списка/статистики: точки-вкладки и подсказка. Порядок
// точек слева направо: суперфетч, основной список, статистика; активная вкладка
// подсвечена. В суперфетче панель не показывается (там своя вёрстка).
func tabBar(cur Mode) string {
	order := []Mode{FetchMode, ListMode, StatsMode}
	dots := make([]string, len(order))
	for i, mo := range order {
		if mo == cur {
			dots[i] = cyanStyle.Render("●")
		} else {
			dots[i] = dimStyle.Render("○")
		}
	}
	return "  " + strings.Join(dots, " ") +
		dimStyle.Render("   ↑↓ — листать · ←→ — режим · q — выход")
}

// pinTabBar пришпиливает панель вкладок к нижней строке экрана: контент сверху,
// добивка пустыми строками, вкладки — последней строкой. Так точки-вкладки
// всегда на одной высоте при переключении режимов и не «скачут».
func pinTabBar(content []string, m Model) string {
	for len(content) < m.height-1 {
		content = append(content, "")
	}
	content = append(content, tabBar(m.mode))
	return strings.Join(content, "\n")
}

// panel рисует рамку-секцию (boxSection) с отступом слева на 2 — как остальная
// вёрстка режимов.
func panel(title string, rows []string, inner int) []string {
	out := boxSection(title, rows, inner)
	for i := range out {
		out[i] = "  " + out[i]
	}
	return out
}

// wrapLines переносит текст по ширине width и красит каждую строку стилем st.
func wrapLines(text string, width int, st lipgloss.Style) []string {
	var out []string
	for _, ln := range strings.Split(wrapText(text, width), "\n") {
		out = append(out, st.Render(ln))
	}
	return out
}

type Model struct {
	selected int
	mode     Mode
	results  []detect.ArchetypeResult
	width    int
	height   int
	reqW     int // минимальная ширина для списка/статистики
	reqH     int // минимальная высота для них же

	// Суперфетч со своим требованием по размеру — оно применяется только в
	// этом режиме и зависит от под-режима: [0] — краткий, [1] — полный.
	fetchReqW [2]int
	fetchReqH [2]int
	fetchFull bool // полный (true) или краткий (false) вид суперфетча

	animating    bool      // идёт ли анимация заполнения полосок
	animStart    time.Time // момент первого тика (для расчёта прогресса)
	animProgress float64   // 0..1 — доля заполнения полосок
}

func NewModel(results []detect.ArchetypeResult, animate bool) Model {
	m := Model{selected: 0, results: results, animating: animate, fetchFull: false}
	m.reqW, m.reqH = computeRequiredSize(results)
	return m
}

// computeRequiredSize считает минимальный размер терминала для базовых режимов
// (список/статистика) — максимум по ним. Требование суперфетча считается
// отдельно и лениво (см. ensureFetchReq), чтобы тяжёлый сбор системной инфы не
// замедлял запуск. Считается один раз при создании модели.
func computeRequiredSize(results []detect.ArchetypeResult) (reqW, reqH int) {
	// Строка списка внутри панели: отступ+рамка (6) + «▶ » (2) + метка + добивка
	// + «  » (2) + «100» (3) + «%» (1) + «  » (2) + бар (20).
	maxLabel := 0
	for _, r := range results {
		if w := lipgloss.Width(r.Label); w > maxLabel {
			maxLabel = w
		}
	}
	reqW = 6 + 2 + maxLabel + 2 + 3 + 1 + 2 + 20
	if reqW < 48 {
		reqW = 48
	}

	// Статистика: панели шире, чем список меток, — берём их натуральную ширину,
	// собрав тело на заведомо большом «холсте». Строки во всю ширину холста
	// (== canvas) пропускаем: они не задают минимально требуемую ширину.
	const canvas = 400
	sBody := statsBody(Model{results: results, mode: StatsMode, width: canvas, height: canvas})
	for _, l := range sBody {
		if w := lipgloss.Width(l); w < canvas && w > reqW {
			reqW = w
		}
	}

	// Высота: максимум по режимам плюс строка панели вкладок (её добавит
	// pinTabBar). У списка высота зависит от переноса описания/«что повлияло»,
	// а те переносятся по ширине reqW — берём худший класс.
	reqH = len(sBody) + 1
	lm := Model{results: results, mode: ListMode, width: reqW, height: canvas}
	worst := func() {
		if h := len(listBody(lm)) + 1; h > reqH {
			reqH = h
		}
	}
	if len(results) == 0 {
		worst()
	}
	for sel := range results {
		lm.selected = sel
		worst()
	}
	return
}

// ensureFetchReq лениво считает требование размера для текущего под-режима
// суперфетча (краткий/полный). Тяжёлый сбор системной инфы (df, lspci, счётчики
// пакетов, flatpak/snap…) откладывается до первого входа в суперфетч, чтобы не
// замедлять запуск; результат кешируется в модели и в fetchOnce.
func (m Model) ensureFetchReq() Model {
	i := 0
	if m.fetchFull {
		i = 1
	}
	if m.fetchReqW[i] == 0 {
		const canvas = 400
		fv := renderFetch(Model{mode: FetchMode, width: canvas, height: canvas, fetchFull: m.fetchFull})
		m.fetchReqW[i] = maxContentWidth(fv, canvas)
		m.fetchReqH[i] = lineCount(fv)
	}
	return m
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
			// В суперфетче ↑↓ циклично переключают детализацию (краткий/полный),
			// в остальных режимах — листают список.
			if m.mode == FetchMode {
				m.fetchFull = !m.fetchFull
			} else if m.selected < len(m.results)-1 {
				m.selected++
			}
		case "k", "л", "up":
			if m.mode == FetchMode {
				m.fetchFull = !m.fetchFull
			} else if m.selected > 0 {
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
	// Требование суперфетча считаем лениво — только когда пользователь реально
	// в этом режиме (первый вход запускает сбор системной инфы).
	if m.mode == FetchMode {
		m = m.ensureFetchReq()
	}
	return m, nil
}

func (m Model) View() string {
	// Пока не пришёл первый WindowSizeMsg — не мигаем сообщением «0×0».
	if m.width == 0 && m.height == 0 {
		return ""
	}
	// Требование по размеру — базовое для всех режимов, кроме суперфетча:
	// у него своё и зависит от под-режима (краткий/полный).
	reqW, reqH := m.reqW, m.reqH
	if m.mode == FetchMode {
		i := 0
		if m.fetchFull {
			i = 1
		}
		reqW, reqH = m.fetchReqW[i], m.fetchReqH[i]
	}
	if m.width < reqW || m.height < reqH {
		return tooSmallView(m, reqW, reqH)
	}
	switch m.mode {
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
func tooSmallView(m Model, reqW, reqH int) string {
	curW := greenStyle
	if m.width < reqW {
		curW = redStyle
	}
	curH := greenStyle
	if m.height < reqH {
		curH = redStyle
	}
	cur := curW.Render(fmt.Sprintf("%d", m.width)) +
		dimStyle.Render("×") + curH.Render(fmt.Sprintf("%d", m.height))

	lines := []string{
		boldStyle.Render("🐧 Окно слишком маленькое"),
		"",
		fmt.Sprintf("Нужно: %s", greenStyle.Render(fmt.Sprintf("%d×%d", reqW, reqH))),
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
