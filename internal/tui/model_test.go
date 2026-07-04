package tui

import (
	"strings"
	"testing"
	"unicode/utf8"

	"github.com/Axawys/lxprofiler/internal/detect"
)

func TestMakeBrokenBarRuneWidth(t *testing.T) {
	for _, w := range []int{1, 5, 20, 33} {
		bar := makeBrokenBar(w)
		if n := utf8.RuneCountInString(bar); n != w {
			t.Errorf("makeBrokenBar(%d): got %d runes, want %d (%q)", w, n, w, bar)
		}
		if strings.ContainsRune(bar, utf8.RuneError) {
			t.Errorf("makeBrokenBar(%d): contains broken rune (byte-slicing bug): %q", w, bar)
		}
	}
	if makeBrokenBar(0) != "" {
		t.Errorf("makeBrokenBar(0) should be empty")
	}
}

func TestMakeBar(t *testing.T) {
	bar := makeBar(50, 20)
	if n := utf8.RuneCountInString(bar); n != 20 {
		t.Errorf("makeBar(50,20): got %d runes, want 20", n)
	}
	if got := strings.Count(bar, "█"); got != 10 {
		t.Errorf("makeBar(50,20): got %d filled cells, want 10", got)
	}
	// negative pct routes to the broken bar (used for masked mystery classes)
	if makeBar(-1, 20) != makeBrokenBar(20) {
		t.Errorf("makeBar(-1,20) should equal makeBrokenBar(20)")
	}
}

func TestBarAnimation(t *testing.T) {
	res := []detect.ArchetypeResult{
		{Key: "programmer", Label: "Программист", NormScore: 100, Reason: "x"},
	}
	m := Model{results: res, mode: ListMode, width: 80, height: 400, animating: true}

	m.animProgress = 0.0
	c0 := strings.Count(renderList(m), "█")
	m.animProgress = 0.5
	c5 := strings.Count(renderList(m), "█")
	m.animProgress = 1.0
	c10 := strings.Count(renderList(m), "█")

	if !(c0 < c5 && c5 < c10) {
		t.Errorf("bars should grow with progress: prog0=%d prog0.5=%d prog1=%d", c0, c5, c10)
	}
	if c0 != 0 {
		t.Errorf("at progress 0 the bar should be empty, got %d filled cells", c0)
	}
}

func TestRaceOrder(t *testing.T) {
	// Финальный порядок (по очкам): Яблоко(100) > Апельсин(50) > Банан(10).
	// Алфавитный (стартовый): Апельсин < Банан < Яблоко.
	res := []detect.ArchetypeResult{
		{Key: "a", Label: "Яблоко", NormScore: 100},
		{Key: "b", Label: "Апельсин", NormScore: 50},
		{Key: "c", Label: "Банан", NormScore: 10},
	}
	labels := func(rs []detect.ArchetypeResult) []string {
		out := make([]string, len(rs))
		for i, r := range rs {
			out[i] = r.Label
		}
		return out
	}
	m := Model{results: res, mode: ListMode, animating: true}

	m.animProgress = 0
	start := labels(m.orderedResults())
	if !equalStrings(start, []string{"Апельсин", "Банан", "Яблоко"}) {
		t.Errorf("start order = %v, want alphabetical [Апельсин Банан Яблоко]", start)
	}

	m.animProgress = 1
	end := labels(m.orderedResults())
	if !equalStrings(end, []string{"Яблоко", "Апельсин", "Банан"}) {
		t.Errorf("end order = %v, want score order [Яблоко Апельсин Банан]", end)
	}
}

func equalStrings(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

func sampleResults() []detect.ArchetypeResult {
	return []detect.ArchetypeResult{
		{Key: "programmer", Label: "Программист", NormScore: 100, Reason: "go, git, make"},
		{Key: "gamer", Label: "Геймер", NormScore: 40, Reason: "steam"},
	}
}

func TestRequiredSizePositive(t *testing.T) {
	w, h, fw, fh := computeRequiredSize(sampleResults())
	if w < 48 {
		t.Errorf("reqW = %d, want >= 48 (floor)", w)
	}
	if h < 10 {
		t.Errorf("reqH = %d, want >= 10", h)
	}
	for i := range fw {
		if fw[i] < 1 || fh[i] < 1 {
			t.Errorf("fetchReq[%d] = %dx%d, want positive", i, fw[i], fh[i])
		}
	}
}

func TestViewGate(t *testing.T) {
	m := NewModel(sampleResults(), false)

	// меньше требуемого → сообщение о размере
	m.width, m.height = 20, 8
	if !strings.Contains(m.View(), "слишком") {
		t.Errorf("small terminal should show the too-small message")
	}

	// достаточно большого — рендерится утилита, не сообщение
	m.width, m.height = m.reqW+5, m.reqH+5
	v := m.View()
	if strings.Contains(v, "слишком") {
		t.Errorf("large enough terminal should NOT show the too-small message")
	}
	// в режиме списка всегда есть полоска заполнения (символ бара)
	if !strings.Contains(v, "█") {
		t.Errorf("large enough terminal should render the list (progress bars)")
	}

	// нулевой размер (до первого WindowSizeMsg) — пусто, без мигания
	m.width, m.height = 0, 0
	if m.View() != "" {
		t.Errorf("zero size should render empty, got %q", m.View())
	}
}
