package tui

import "testing"

func TestExtractCommand(t *testing.T) {
	tests := []struct {
		line, want string
	}{
		{"vim foo.txt", "vim foo.txt"},                 // bash / zsh plain
		{": 1609459200:0;vim foo.txt", "vim foo.txt"}, // zsh extended history
		{": 1609459200:12;git commit -m x", "git commit -m x"},
		{"- cmd: nvim config.lua", "nvim config.lua"}, // fish
		{"  when: 1609459200", ""},                    // fish metadata
		{"#1609459200", ""},                           // bash HISTTIMEFORMAT timestamp
		{"   ", ""},                                    // blank
	}
	for _, tt := range tests {
		if got := extractCommand(tt.line); got != tt.want {
			t.Errorf("extractCommand(%q) = %q, want %q", tt.line, got, tt.want)
		}
	}
}

func TestTallyCommandsEditors(t *testing.T) {
	// vim вызывался 5 раз с аргументами и под sudo — все должны считаться;
	// helix (hx) и code — тоже редакторы, их надо учесть по имени.
	commands := []string{
		"vim a.txt", "vim b.txt", "vim c.go", "sudo vim /etc/hosts", "vi notes",
		"nvim init.lua", "nvim plugin.lua",
		"nano quick.txt",
		"hx main.rs", "code .",
		"ls -la", "ls", "git status",
	}
	var r StatsResult
	tallyCommands(commands, nil, &r)

	if r.Editors["vim"] != 5 {
		t.Errorf("Editors[vim] = %d, want 5 (vim×3 + sudo vim + vi)", r.Editors["vim"])
	}
	if r.Editors["nvim"] != 2 {
		t.Errorf("Editors[nvim] = %d, want 2", r.Editors["nvim"])
	}
	if r.Editors["nano"] != 1 {
		t.Errorf("Editors[nano] = %d, want 1", r.Editors["nano"])
	}
	if r.Editors["helix"] != 1 {
		t.Errorf("Editors[helix] = %d, want 1 (hx)", r.Editors["helix"])
	}
	if r.Editors["vscode"] != 1 {
		t.Errorf("Editors[vscode] = %d, want 1 (code)", r.Editors["vscode"])
	}
	if r.SudoCount != 1 {
		t.Errorf("SudoCount = %d, want 1", r.SudoCount)
	}
	if r.TotalCmds != len(commands) {
		t.Errorf("TotalCmds = %d, want %d", r.TotalCmds, len(commands))
	}
	// топ-1 команда — литеральный "vim" (4×: vim×3 + sudo vim); "vi" считается
	// отдельной командой, но в войне редакторов суммируется с vim (VimCount=5).
	if len(r.TopCmds) == 0 || r.TopCmds[0].Cmd != "vim" || r.TopCmds[0].Count != 4 {
		t.Errorf("TopCmds[0] = %+v, want {vim 4}", r.TopCmds)
	}
}

func TestTallyCommandsTypos(t *testing.T) {
	commands := []string{"sl", "gti status", "claer", "cd..", "ls", "grep foo", "sl"}
	var r StatsResult
	tallyCommands(commands, nil, &r)
	// sl×2, gti, claer, cd.. = 5 опечаток
	if r.TypoCount != 5 {
		t.Errorf("TypoCount = %d, want 5", r.TypoCount)
	}
}

func TestFetchAliasNames(t *testing.T) {
	config := `
alias ff='fastfetch'
alias neo=neofetch
alias ll='ls -la'
abbr -a nf fastfetch
alias sysinfo="command fastfetch --logo arch"
`
	got := map[string]bool{}
	for _, n := range fetchAliasNames(config) {
		got[n] = true
	}
	for _, want := range []string{"ff", "neo", "nf", "sysinfo"} {
		if !got[want] {
			t.Errorf("fetchAliasNames: alias %q should resolve to a fetch tool", want)
		}
	}
	if got["ll"] {
		t.Errorf("fetchAliasNames: 'll' (ls -la) must NOT be a fetch alias")
	}
}

func TestTallyCommandsFetchAliases(t *testing.T) {
	commands := []string{"ff", "fastfetch", "ff", "neofetch", "ls"}
	fetchAlias := map[string]bool{"ff": true}
	var r StatsResult
	tallyCommands(commands, fetchAlias, &r)
	// ff×2 (алиас) + fastfetch + neofetch = 4
	if r.FFCount != 4 {
		t.Errorf("FFCount = %d, want 4 (ff×2 + fastfetch + neofetch)", r.FFCount)
	}
}
