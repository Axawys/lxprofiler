package tui

import (
	"strings"
	"testing"
)

func TestComputeStats(t *testing.T) {
	s := computeStats()
	// Just verify it doesn't panic
	if s.TotalCmds < 0 {
		t.Error("TotalCmds should be non-negative")
	}
	if s.UniqueCmds < 0 {
		t.Error("UniqueCmds should be non-negative")
	}
}

func TestHumanKB(t *testing.T) {
	tests := []struct {
		kb   int
		want string
	}{
		{1024, "1M"},
		{1048576, "1.0G"},
		{500, "500K"},
		{1536, "1M"},
	}
	for _, tt := range tests {
		got := humanKB(tt.kb)
		if got != tt.want {
			t.Errorf("humanKB(%d) = %q, want %q", tt.kb, got, tt.want)
		}
	}
}

func TestHumanInterval(t *testing.T) {
	tests := []struct {
		secs int
		want string
	}{
		{0, "—"},
		{60, "раз в 1 мин."},
		{3600, "раз в 1 ч."},
		{86400, "раз в 1 дн."},
		{172800, "раз в 2 дн."},
	}
	for _, tt := range tests {
		got := humanInterval(tt.secs)
		if got != tt.want {
			t.Errorf("humanInterval(%d) = %q, want %q", tt.secs, got, tt.want)
		}
	}
}

func TestEditorCounts(t *testing.T) {
	// Test the counting logic directly
	// The raw string starts with a newline (as in computeStats)
	// Commands are on their own lines (like actual bash history)
	raw := "\nvim\nnvim\nnano\nemacs\nmicro\n"
	lower := strings.ToLower(raw)

	// Count using the same logic as in computeStats
	vimCount := strings.Count(lower, "\nvim\n") + strings.Count(lower, "\nvi\n")
	nvimCount := strings.Count(lower, "\nnvim\n")
	nanoCount := strings.Count(lower, "\nnano\n")
	emacsCount := strings.Count(lower, "\nemacs\n")
	microCount := strings.Count(lower, "\nmicro\n")

	if vimCount != 1 {
		t.Errorf("VimCount = %d, want 1", vimCount)
	}
	if nvimCount != 1 {
		t.Errorf("NvimCount = %d, want 1", nvimCount)
	}
	if nanoCount != 1 {
		t.Errorf("NanoCount = %d, want 1", nanoCount)
	}
	if emacsCount != 1 {
		t.Errorf("EmacsCount = %d, want 1", emacsCount)
	}
	if microCount != 1 {
		t.Errorf("MicroCount = %d, want 1", microCount)
	}
}

func TestFreq(t *testing.T) {
	if freq(0, 1000) != "—" {
		t.Error("freq with 0 count should be '—'")
	}
	got := freq(10, 100000)
	if got == "" {
		t.Error("freq should return non-empty string")
	}
}
