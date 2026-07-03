package data

import "testing"

func TestDescribe(t *testing.T) {
	tests := []struct {
		key    string
		pct    int
		expect string
	}{
		{"devops", 100, DescHigh["devops"]},
		{"devops", 60, DescMid["devops"]},
		{"devops", 30, DescLow["devops"]},
		{"devops", 10, DescNone["devops"]},
		{"anonymous", 50, DescLow["anonymous"] + " И как ты вообще не побоялся ставить этот скрипт?"},
		{"anonymous", 30, DescLow["anonymous"]},
	}
	for _, tt := range tests {
		got := Describe(tt.key, tt.pct)
		if got != tt.expect {
			t.Errorf("Describe(%q, %d) = %q, want %q", tt.key, tt.pct, got, tt.expect)
		}
	}
}

func TestLabels(t *testing.T) {
	if len(Labels) != 37 {
		t.Errorf("Labels has %d entries, want 37", len(Labels))
	}
	if Labels["devops"] != "DevOps" {
		t.Errorf("Labels[devops] = %q, want %q", Labels["devops"], "DevOps")
	}
	if Labels["vibe_coder"] != "Вайбкодер" {
		t.Errorf("Labels[vibe_coder] = %q, want %q", Labels["vibe_coder"], "Вайбкодер")
	}
	if Labels["musician"] != "Музыкант" {
		t.Errorf("Labels[musician] = %q, want %q", Labels["musician"], "Музыкант")
	}
	if Labels["embedded"] != "Embedded-разработчик" {
		t.Errorf("Labels[embedded] = %q, want %q", Labels["embedded"], "Embedded-разработчик")
	}
}

func TestHidden(t *testing.T) {
	if !Hidden["normis"] {
		t.Error("normis should be hidden")
	}
	if Hidden["devops"] {
		t.Error("devops should not be hidden")
	}
}

func TestMystery(t *testing.T) {
	if !Mystery["vm"] {
		t.Error("vm should be mystery")
	}
	if Mystery["programmer"] {
		t.Error("programmer should not be mystery")
	}
}
