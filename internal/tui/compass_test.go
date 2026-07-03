package tui

import (
	"testing"

	"github.com/Axawys/lxprofiler/internal/detect"
)

func TestComputeCompass(t *testing.T) {
	// Save original scores
	origScore := make(map[string]int)
	for k, v := range detect.Score {
		origScore[k] = v
	}
	defer func() {
		for k, v := range origScore {
			detect.Score[k] = v
		}
	}()

	detect.Score["devops"] = 100
	detect.Score["programmer"] = 50
	compass := ComputeCompass()
	if compass.CX == 0 && compass.CY == 0 {
		t.Error("compass should not be at origin with non-zero scores")
	}
	if compass.Quadrant == "" {
		t.Error("quadrant should not be empty")
	}
}

func TestGetQuadrant(t *testing.T) {
	tests := []struct {
		cx, cy int
		want   string
	}{
		{-20, 20, "Лаборатория"},
		{20, 20, "Цитадель Unix"},
		{-20, -20, "Гладкое будущее"},
		{20, -20, "Тёплая гавань"},
		{0, 0, "Центрист"},
	}
	for _, tt := range tests {
		got := getQuadrant(tt.cx, tt.cy)
		if len(got) == 0 {
			t.Errorf("getQuadrant(%d, %d) = empty", tt.cx, tt.cy)
		}
	}
}
