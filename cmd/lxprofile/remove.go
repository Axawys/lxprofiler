package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func dirExists(path string) bool {
	fi, err := os.Stat(path)
	return err == nil && fi.IsDir()
}

// doRemove удаляет бинарник, наши короткие симлинки и каталог состояния.
func doRemove() error {
	exe := ourBinary()
	if exe == "" {
		return fmt.Errorf("не удалось определить путь бинарника")
	}
	bin := filepath.Dir(exe)

	var links []string
	for _, name := range shortCmds {
		cand := filepath.Join(bin, name)
		if isOurSymlink(cand, exe) {
			links = append(links, cand)
		}
		if found, err := exec.LookPath(name); err == nil && isOurSymlink(found, exe) {
			links = append(links, found)
		}
	}
	links = dedupe(links)
	state := stateDir()

	fmt.Println("Будет удалено:")
	fmt.Printf("  • бинарник  %s\n", exe)
	for _, l := range links {
		fmt.Printf("  • симлинк   %s (%s)\n", l, filepath.Base(l))
	}
	if dirExists(state) {
		fmt.Printf("  • состояние %s\n", state)
	}

	fmt.Print("Продолжить удаление? [y/N] ")
	reader := bufio.NewReader(os.Stdin)
	line, _ := reader.ReadString('\n')
	switch strings.ToLower(strings.TrimSpace(line)) {
	case "y", "yes", "да", "д":
	default:
		fmt.Println("Отменено.")
		return fmt.Errorf("отменено")
	}

	for _, l := range links {
		_ = os.Remove(l)
	}
	_ = os.RemoveAll(state)
	_ = os.Remove(exe)
	fmt.Println("lxprofile удалён. Спасибо, что пользовались! 🐧")
	return nil
}
