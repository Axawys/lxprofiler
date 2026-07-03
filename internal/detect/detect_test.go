package detect

import (
	"os"
	"testing"

	"github.com/Axawys/lxprofiler/internal/data"
)

func TestHas(t *testing.T) {
	if !Has("ls") {
		t.Error("ls should exist")
	}
	if Has("nonexistent_command_xyz123") {
		t.Error("nonexistent command should not exist")
	}
}

func TestAdd(t *testing.T) {
	Score["test"] = 0
	Reasons["test"] = ""
	Add("test", 10, "reason1")
	if Score["test"] != 10 {
		t.Errorf("Score = %d, want 10", Score["test"])
	}
	if Reasons["test"] != "reason1" {
		t.Errorf("Reasons = %q, want %q", Reasons["test"], "reason1")
	}
	Add("test", 5, "reason2")
	if Score["test"] != 15 {
		t.Errorf("Score = %d, want 15", Score["test"])
	}
	if Reasons["test"] != "reason1, reason2" {
		t.Errorf("Reasons = %q, want %q", Reasons["test"], "reason1, reason2")
	}
}

func TestUsed(t *testing.T) {
	origBehavior := Behavior
	defer func() { Behavior = origBehavior }()

	Behavior = "docker run hello\nkubectl get pods\n"
	if !Used("docker") {
		t.Error("should find docker")
	}
	if !Used("kubectl") {
		t.Error("should find kubectl")
	}
	if Used("nonexistent") {
		t.Error("should not find nonexistent")
	}
}

func TestHasUsed(t *testing.T) {
	origBehavior := Behavior
	defer func() { Behavior = origBehavior }()

	Behavior = "ls -la\ncat file.txt\n"
	if !HasUsed("ls") {
		t.Error("ls should be used")
	}
	if HasUsed("nonexistent_xyz") {
		t.Error("nonexistent should not be used")
	}
}

func TestDetect(t *testing.T) {
	for key := range data.Labels {
		Score[key] = 0
		Reasons[key] = ""
	}
	Detect()

	if Score["programmer"] == 0 && Score["devops"] == 0 {
		t.Error("at least one class should have non-zero score")
	}
}

func TestNormalize(t *testing.T) {
	for key := range data.Labels {
		Score[key] = 0
		Reasons[key] = ""
	}
	Score["programmer"] = 100
	Score["devops"] = 50
	Reasons["programmer"] = "test"
	Reasons["devops"] = "test"

	results := Normalize()
	if len(results) == 0 {
		t.Error("Normalize returned empty results")
	}
	if results[0].NormScore != 100 {
		t.Errorf("top score = %d, want 100", results[0].NormScore)
	}
}

func TestFileExists(t *testing.T) {
	if !fileExists("/etc/passwd") {
		t.Error("/etc/passwd should exist")
	}
	if fileExists("/nonexistent/file") {
		t.Error("nonexistent file should not exist")
	}
}

func TestDirExists(t *testing.T) {
	if !dirExists("/tmp") {
		t.Error("/tmp should exist")
	}
	if dirExists("/nonexistent/dir") {
		t.Error("nonexistent dir should not exist")
	}
}

func TestHasPythonSitePackages(t *testing.T) {
	home := os.Getenv("HOME")
	if home == "" {
		home = "/root"
	}
	// This may or may not exist, just test it doesn't panic
	hasPythonSitePackages(home)
}
