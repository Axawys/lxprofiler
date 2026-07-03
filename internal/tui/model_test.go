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
		if !strings.HasSuffix(bar, "?") {
			t.Errorf("makeBrokenBar(%d): should end with '?', got %q", w, bar)
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

func sampleResults() []detect.ArchetypeResult {
	return []detect.ArchetypeResult{
		{Key: "programmer", Label: "Программист", NormScore: 100, Reason: "go, git, make"},
		{Key: "gamer", Label: "Геймер", NormScore: 40, Reason: "steam"},
	}
}

func TestRequiredSizePositive(t *testing.T) {
	w, h := computeRequiredSize(sampleResults())
	if w < 48 {
		t.Errorf("reqW = %d, want >= 48 (floor)", w)
	}
	if h < 10 {
		t.Errorf("reqH = %d, want >= 10", h)
	}
}

func TestViewGate(t *testing.T) {
	m := NewModel(sampleResults())

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
	if !strings.Contains(v, "Профиль архетипов") {
		t.Errorf("large enough terminal should render the list")
	}

	// нулевой размер (до первого WindowSizeMsg) — пусто, без мигания
	m.width, m.height = 0, 0
	if m.View() != "" {
		t.Errorf("zero size should render empty, got %q", m.View())
	}
}
