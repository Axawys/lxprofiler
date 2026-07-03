package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
)

// shortCmds — короткая команда lx и её «слитные» формы с флагом.
// lxu == lx -u, lxs == lx -s, lxv == lx -v, lxh == lx -h, lxc == lx -c,
// lxrm == lx --rm (разбор имени вызова — в main).
var shortCmds = []string{"lx", "lxu", "lxs", "lxv", "lxh", "lxc", "lxrm"}

// stateDir — каталог для маркеров (отказ от обновления, разовые операции).
func stateDir() string {
	base := os.Getenv("XDG_STATE_HOME")
	if base == "" {
		home, _ := os.UserHomeDir()
		base = filepath.Join(home, ".local", "state")
	}
	return filepath.Join(base, "lxprofiler")
}

func declinedFile() string { return filepath.Join(stateDir(), "update_declined") }
func lxMarker() string     { return filepath.Join(stateDir(), "lx_setup_done") }

func markerExists(path string) bool { _, err := os.Stat(path); return err == nil }

func writeMarker(path, content string) {
	_ = os.MkdirAll(filepath.Dir(path), 0o755)
	_ = os.WriteFile(path, []byte(content), 0o644)
}

// ourBinary возвращает реальный путь запущенного бинарника (симлинки развёрнуты).
func ourBinary() string {
	exe, err := os.Executable()
	if err != nil {
		return ""
	}
	if resolved, err := filepath.EvalSymlinks(exe); err == nil {
		return resolved
	}
	return exe
}

func isOurSymlink(link, exe string) bool {
	fi, err := os.Lstat(link)
	if err != nil || fi.Mode()&os.ModeSymlink == 0 {
		return false
	}
	tgt, err := filepath.EvalSymlinks(link)
	return err == nil && tgt == exe
}

// shortTaken — занято ли имя чужой командой в PATH или алиасом/функцией shell.
func shortTaken(name, link, exe string) bool {
	if found, err := exec.LookPath(name); err == nil && found != link {
		if resolved, e := filepath.EvalSymlinks(found); e != nil || resolved != exe {
			return true
		}
	}
	return aliasDefined(name)
}

// aliasDefined ищет alias/function name в типичных shell-конфигах.
func aliasDefined(name string) bool {
	home, _ := os.UserHomeDir()
	re := regexp.MustCompile(`(?i)^\s*alias\s+` + name + `[\s=]|(^|\s)function\s+` + name + `(\s|\(|$)|(^|\s)` + name + `\s*\(\)`)
	files := []string{
		".bashrc", ".bash_aliases", ".zshrc", ".zsh_aliases",
		".aliases", ".profile", ".config/fish/config.fish",
	}
	for _, f := range files {
		path := filepath.Join(home, f)
		fh, err := os.Open(path)
		if err != nil {
			continue
		}
		sc := bufio.NewScanner(fh)
		for sc.Scan() {
			if re.MatchString(sc.Text()) {
				fh.Close()
				return true
			}
		}
		fh.Close()
	}
	return false
}

// ensureShortCommands создаёт симлинки коротких команд для свободных имён.
// Дорогая часть выполняется один раз (маркер lxMarker); force — перепроверить.
func ensureShortCommands(force bool) {
	marker := lxMarker()
	if !force && markerExists(marker) {
		return
	}
	verbose := !markerExists(marker)

	exe := ourBinary()
	if exe == "" {
		return
	}
	bin := filepath.Dir(exe)

	for _, name := range shortCmds {
		link := filepath.Join(bin, name)
		if isOurSymlink(link, exe) {
			continue
		}
		if _, err := os.Lstat(link); err == nil {
			if name == "lx" && verbose {
				fmt.Fprintln(os.Stderr, "Короткая команда lx уже занята — используйте lxprofile.")
			}
			continue
		}
		if shortTaken(name, link, exe) {
			if name == "lx" && verbose {
				fmt.Fprintln(os.Stderr, "Короткая команда lx уже занята — используйте lxprofile.")
			}
			continue
		}
		if err := os.Symlink(exe, link); err == nil && name == "lx" && verbose {
			fmt.Fprintln(os.Stderr, "Создана короткая команда: lx (= lxprofile) и слитные формы lxu/lxs/lxc…")
		}
	}
	writeMarker(marker, "")
}

// currentBinDir — каталог, где лежит наш бинарник (для установки/удаления симлинков).
func currentBinDir() string { return filepath.Dir(ourBinary()) }

// dedupe removes duplicate strings preserving order.
func dedupe(in []string) []string {
	seen := map[string]bool{}
	var out []string
	for _, s := range in {
		if s != "" && !seen[s] {
			seen[s] = true
			out = append(out, s)
		}
	}
	return out
}
