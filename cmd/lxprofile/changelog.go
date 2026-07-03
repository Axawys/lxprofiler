package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strings"

	root "github.com/Axawys/lxprofiler"
)

var verHeaderRe = regexp.MustCompile(`^##\s+v?([0-9]+\.[0-9]+\.[0-9]+)`)

// doChanges печатает changelog из встроенного CHANGELOG.md:
//
//	want пуст     → последние 5 версий;
//	want = версия → запись этой версии (или список доступных при ошибке).
func doChanges(want string) {
	text := root.ChangelogMD
	if strings.TrimSpace(text) == "" {
		fmt.Fprintln(os.Stderr, "CHANGELOG пуст или не встроен в сборку.")
		return
	}
	want = strings.TrimPrefix(want, "v")
	if want != "" {
		block := changelogBlock(text, want)
		if strings.TrimSpace(block) == "" {
			fmt.Fprintf(os.Stderr, "Версия %s не найдена. Доступные версии:\n", want)
			for _, v := range changelogVersions(text) {
				fmt.Fprintf(os.Stderr, "  %s\n", v)
			}
			return
		}
		fmt.Println(strings.TrimRight(block, "\n"))
		return
	}
	fmt.Println(strings.TrimRight(changelogLastN(text, 5), "\n"))
}

// changelogBlock возвращает блок `## ver` … до следующего заголовка версии.
func changelogBlock(text, ver string) string {
	var out []string
	capturing := false
	scanner := bufio.NewScanner(strings.NewReader(text))
	for scanner.Scan() {
		line := scanner.Text()
		if m := verHeaderRe.FindStringSubmatch(line); m != nil {
			if capturing {
				break
			}
			if m[1] == ver {
				capturing = true
				out = append(out, line)
			}
			continue
		}
		if capturing {
			out = append(out, line)
		}
	}
	return strings.Join(out, "\n")
}

// changelogLastN возвращает первые n версий сверху файла.
func changelogLastN(text string, n int) string {
	var out []string
	count := 0
	scanner := bufio.NewScanner(strings.NewReader(text))
	for scanner.Scan() {
		line := scanner.Text()
		if verHeaderRe.MatchString(line) {
			count++
			if count > n {
				break
			}
		}
		if count >= 1 {
			out = append(out, line)
		}
	}
	return strings.Join(out, "\n")
}

// changelogVersions возвращает все версии из файла (в порядке следования).
func changelogVersions(text string) []string {
	var vers []string
	scanner := bufio.NewScanner(strings.NewReader(text))
	for scanner.Scan() {
		if m := verHeaderRe.FindStringSubmatch(scanner.Text()); m != nil {
			vers = append(vers, m[1])
		}
	}
	return vers
}

// changelogSince форматирует пункты всех версий новее current (для уведомления).
func changelogSince(text, current string) string {
	var b strings.Builder
	inBlock := false
	any := false
	scanner := bufio.NewScanner(strings.NewReader(text))
	for scanner.Scan() {
		line := scanner.Text()
		if m := verHeaderRe.FindStringSubmatch(line); m != nil {
			if versionGT(m[1], current) {
				inBlock = true
				if !any {
					b.WriteString("Что нового:\n")
					any = true
				}
				fmt.Fprintf(&b, "  %s:\n", m[1])
			} else {
				inBlock = false
			}
			continue
		}
		if inBlock && strings.TrimSpace(line) != "" {
			fmt.Fprintf(&b, "    %s\n", line)
		}
	}
	return b.String()
}
