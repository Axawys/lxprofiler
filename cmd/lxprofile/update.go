package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

const repoSlug = "Axawys/lxprofiler"

// parseVer разбирает "vX.Y.Z" в [3]int, игнорируя нецифровые хвосты.
func parseVer(s string) [3]int {
	s = strings.TrimPrefix(strings.TrimSpace(s), "v")
	parts := strings.SplitN(s, ".", 3)
	var v [3]int
	for i := 0; i < 3 && i < len(parts); i++ {
		n := 0
		for _, c := range parts[i] {
			if c < '0' || c > '9' {
				break
			}
			n = n*10 + int(c-'0')
		}
		v[i] = n
	}
	return v
}

// versionGT → true, если версия a строго новее b (semver).
func versionGT(a, b string) bool {
	va, vb := parseVer(a), parseVer(b)
	for i := 0; i < 3; i++ {
		if va[i] != vb[i] {
			return va[i] > vb[i]
		}
	}
	return false
}

// remoteLatest узнаёт версию последнего релиза на GitHub (без префикса v).
func remoteLatest() (string, error) {
	url := "https://api.github.com/repos/" + repoSlug + "/releases/latest"
	req, _ := http.NewRequest(http.MethodGet, url, nil)
	req.Header.Set("Accept", "application/vnd.github+json")
	req.Header.Set("User-Agent", "lxprofile")
	client := &http.Client{Timeout: 6 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("GitHub API: %s", resp.Status)
	}
	var payload struct {
		TagName string `json:"tag_name"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&payload); err != nil {
		return "", err
	}
	return strings.TrimPrefix(payload.TagName, "v"), nil
}

// assetName — имя релизного бинарника под текущую ОС/архитектуру.
func assetName() string {
	return fmt.Sprintf("lxprofile-%s-%s", runtime.GOOS, runtime.GOARCH)
}

// doUpdate обновляет (или откатывает) бинарник до релиза с GitHub.
// version пуст → последний релиз; иначе конкретная версия (откат/фиксация).
func doUpdate(version string) error {
	exe := ourBinary()
	if exe == "" {
		return fmt.Errorf("не удалось определить путь бинарника")
	}

	var ver string
	if version == "" {
		latest, err := remoteLatest()
		if err != nil {
			return fmt.Errorf("не удалось узнать последнюю версию: %v", err)
		}
		ver = latest
		if !versionGT(ver, Version) {
			fmt.Printf("У вас уже последняя версия: %s\n", Version)
			return nil
		}
	} else {
		ver = strings.TrimPrefix(version, "v")
	}

	tag := "v" + ver
	url := fmt.Sprintf("https://github.com/%s/releases/download/%s/%s", repoSlug, tag, assetName())
	fmt.Printf("Загружаю %s (%s)…\n", ver, assetName())
	if err := selfReplace(url, exe); err != nil {
		return fmt.Errorf("обновление не удалось: %v", err)
	}

	_ = os.Remove(declinedFile()) // ручное обновление снимает «отказ»
	ensureShortCommands(true)
	fmt.Printf("Готово. Текущая версия: %s\n", ver)
	return nil
}

// selfReplace скачивает бинарник и атомарно заменяет текущий exe.
// На Linux rename поверх запущенного файла работает (inode остаётся открытым).
func selfReplace(url, exe string) error {
	req, _ := http.NewRequest(http.MethodGet, url, nil)
	req.Header.Set("User-Agent", "lxprofile")
	client := &http.Client{Timeout: 90 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("скачивание %s: %s", url, resp.Status)
	}

	dir := filepath.Dir(exe)
	tmp, err := os.CreateTemp(dir, ".lxprofile-new-*")
	if err != nil {
		return err
	}
	tmpName := tmp.Name()
	defer os.Remove(tmpName)

	if _, err := io.Copy(tmp, resp.Body); err != nil {
		tmp.Close()
		return err
	}
	if err := tmp.Close(); err != nil {
		return err
	}
	if err := os.Chmod(tmpName, 0o755); err != nil {
		return err
	}
	return os.Rename(tmpName, exe)
}

// fetchRemoteChangelog тянет CHANGELOG.md из репозитория на указанном теге.
func fetchRemoteChangelog(tag string) string {
	url := fmt.Sprintf("https://raw.githubusercontent.com/%s/%s/CHANGELOG.md", repoSlug, tag)
	req, _ := http.NewRequest(http.MethodGet, url, nil)
	req.Header.Set("User-Agent", "lxprofile")
	client := &http.Client{Timeout: 6 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return ""
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return ""
	}
	b, _ := io.ReadAll(resp.Body)
	return string(b)
}

// startBackgroundCheck запускает проверку обновлений в фоне (если не отказывались).
// Возвращает канал, в который попадёт новая версия (или закроется пустым).
func startBackgroundCheck() chan string {
	ch := make(chan string, 1)
	if markerExists(declinedFile()) {
		close(ch)
		return ch
	}
	go func() {
		defer close(ch)
		ver, err := remoteLatest()
		if err == nil && ver != "" && versionGT(ver, Version) {
			ch <- ver
		}
	}()
	return ch
}

// finishBackgroundCheck ждёт результат проверки и предлагает обновиться.
func finishBackgroundCheck(ch chan string) {
	if ch == nil {
		return
	}
	select {
	case ver, ok := <-ch:
		if ok && ver != "" {
			offerUpdate(ver)
		}
	case <-time.After(3 * time.Second):
	}
}

// offerUpdate печатает предложение обновиться и обрабатывает ответ.
func offerUpdate(remoteVer string) {
	fmt.Printf("Доступна новая версия: \033[31m%s\033[0m -> \033[32m%s\033[0m\n", Version, remoteVer)
	if cl := fetchRemoteChangelog("v" + remoteVer); cl != "" {
		fmt.Print(changelogSince(cl, Version))
	}
	fmt.Print("Обновить сейчас? [y/N] ")

	reader := bufio.NewReader(os.Stdin)
	line, _ := reader.ReadString('\n')
	switch strings.ToLower(strings.TrimSpace(line)) {
	case "y", "yes", "да", "д":
		if err := doUpdate(""); err != nil {
			fmt.Fprintln(os.Stderr, err)
		}
	default:
		writeMarker(declinedFile(), remoteVer+"\n")
		fmt.Println("Хорошо, больше не буду предлагать — пока не обновитесь вручную (lxprofile --update).")
	}
}
