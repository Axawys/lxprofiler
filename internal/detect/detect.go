package detect

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"runtime"
	"sort"
	"strings"
	"time"

	"github.com/Axawys/lxprofiler/internal/data"
)

var Behavior string

func init() {
	home := os.Getenv("HOME")
	if home == "" {
		home = "/root"
	}
	files := []string{
		home + "/.bash_history",
		home + "/.zsh_history",
		home + "/.local/share/fish/fish_history",
		home + "/.bashrc",
		home + "/.zshrc",
		home + "/.bash_aliases",
		home + "/.config/fish/config.fish",
		home + "/.profile",
	}
	for _, f := range files {
		if d, err := os.ReadFile(f); err == nil {
			Behavior += "\n" + string(d)
		}
	}
}

func Used(pattern string) bool {
	if Behavior == "" {
		return true
	}
	// \b — только целые слова (как grep -w в bash-версии): иначе `wg` матчится
	// внутри wget, `vi` внутри vim и т.п., завышая очки.
	re, err := regexp.Compile(`(?i)\b(?:` + pattern + `)\b`)
	if err != nil {
		return strings.Contains(strings.ToLower(Behavior), strings.ToLower(pattern))
	}
	return re.MatchString(Behavior)
}

func HasUsed(cmd string, pattern ...string) bool {
	if !Has(cmd) {
		return false
	}
	p := cmd
	if len(pattern) > 0 {
		p = pattern[0]
	}
	return Used(p)
}

func behav(pat, class, reason string, pts, threshold int) {
	if Behavior == "" {
		return
	}
	// \b — считаем только целые слова (паритет с grep -owE ... -w в bash).
	re, err := regexp.Compile(`(?i)\b(?:` + pat + `)\b`)
	if err != nil {
		return
	}
	n := len(re.FindAllString(Behavior, -1))
	if n >= threshold {
		Add(class, pts, fmt.Sprintf("%s (×%d)", reason, n))
	}
}

func Detect() {
	detectDistro()
	detectDesktop()
	detectSession()
	detectHardware()
	detectDisk()
	detectInit()
	detectBootloader()
	detectLegacy()
	detectAnonymous()
	detectPentester()
	detectDevops()
	detectSysadmin()
	detectProgrammer()
	detectBrowser()
	detectPackageManagers()
	detectGamer()
	detectShell()
	detectTerminal()
	detectDotfiles()
	detectSystemState()
	detectMetaClasses()
	detectVirtualizer()
	detectTiling()
	detectNeovimAddict()
	detectShellCollector()
	detectSelfBuilder()
	detectWaylandWafer()
	detectConsoleLife()
	detectPackageFreak()
	detectMusician()
	detectPhotographer()
	detectVideoEditor()
	detectModeler3d()
	detectWriter()
	detectStreamer()
	detectEmbedded()
	analyzeBehavior()
}

func detectDistro() {
	content, err := os.ReadFile("/etc/os-release")
	if err != nil {
		return
	}
	distroAll := ""
	for _, line := range strings.Split(string(content), "\n") {
		if strings.HasPrefix(line, "PRETTY_NAME=") {
			distroAll += strings.Trim(strings.TrimPrefix(line, "PRETTY_NAME="), "\"") + " "
		}
		if strings.HasPrefix(line, "ID=") {
			distroAll += strings.Trim(strings.TrimPrefix(line, "ID="), "\"") + " "
		}
		if strings.HasPrefix(line, "ID_LIKE=") {
			distroAll += strings.Trim(strings.TrimPrefix(line, "ID_LIKE="), "\"")
		}
	}
	switch {
	case strings.Contains(distroAll, "Kali"):
		Add("pentester", 18, "Kali Linux"); Add("anonymous", 6, "арсенал аудита")
	case strings.Contains(distroAll, "Parrot"):
		Add("pentester", 16, "Parrot OS"); Add("anonymous", 10, "приватность + аудит")
	case strings.Contains(distroAll, "BlackArch"):
		Add("pentester", 18, "BlackArch"); Add("anonymous", 4, "арсенал аудита")
	case strings.Contains(distroAll, "ArchStrike"):
		Add("pentester", 16, "ArchStrike")
	case strings.Contains(distroAll, "Pentoo"):
		Add("pentester", 18, "Pentoo"); Add("old_hacker", 4, "Gentoo-основа")
	case strings.Contains(distroAll, "BackBox"):
		Add("pentester", 16, "BackBox")
	case strings.Contains(distroAll, "Athena"):
		Add("pentester", 16, "Athena OS")
	case strings.Contains(distroAll, "Wifislax"):
		Add("pentester", 14, "Wifislax (Wi-Fi аудит)")
	case strings.Contains(distroAll, "Samurai"):
		Add("pentester", 12, "Samurai WTF")
	case strings.Contains(distroAll, "CAINE"):
		Add("pentester", 12, "CAINE (форензика)")
	case strings.Contains(distroAll, "Network Security Toolkit"):
		Add("pentester", 12, "Network Security Toolkit")
	case strings.Contains(distroAll, "Security Lab") || strings.Contains(distroAll, "Security Spin"):
		Add("pentester", 12, "Fedora Security Lab"); Add("fresh_witness", 6, "Fedora")
	case strings.Contains(distroAll, "Tails"):
		Add("anonymous", 20, "Tails")
	case strings.Contains(distroAll, "Qubes"):
		Add("anonymous", 16, "Qubes OS"); Add("pentester", 6, "изоляция по доменам")
	case strings.Contains(distroAll, "Whonix"):
		Add("anonymous", 18, "Whonix")
	case strings.Contains(distroAll, "Astra"):
		Add("import_substituted", 28, "Astra Linux"); Add("sysadmin", 6, "корпоративная ОС"); Add("anonymous", 4, "мандатный доступ")
	case strings.Contains(distroAll, "RED OS") || strings.Contains(distroAll, "RedOS"):
		Add("import_substituted", 28, "RED OS"); Add("sysadmin", 6, "серверная ОС")
	case strings.Contains(distroAll, "ALT Atomic"):
		Add("atomic", 14, "ALT Atomic"); Add("import_substituted", 22, "ALT (отеч.)"); Add("fresh_witness", 4, "immutable")
	case strings.Contains(distroAll, "ALT"):
		Add("import_substituted", 24, "ALT Linux"); Add("old_hacker", 4, "Sisyphus")
	case strings.Contains(distroAll, "ROSA"):
		Add("import_substituted", 24, "ROSA Linux")
	case strings.Contains(distroAll, "Calculate"):
		Add("import_substituted", 24, "Calculate Linux"); Add("old_hacker", 6, "Gentoo-основа")
	case strings.Contains(distroAll, "Simply"):
		Add("import_substituted", 20, "Simply Linux")
	case strings.Contains(distroAll, "Garuda"):
		Add("gamer", 10, "Garuda"); Add("ricer", 8, "ricing из коробки")
	case strings.Contains(distroAll, "Artix"):
		Add("old_hacker", 12, "Artix (без systemd)"); Add("minimalist", 4, "выбор init")
	case strings.Contains(distroAll, "EndeavourOS"):
		Add("old_hacker", 8, "EndeavourOS"); Add("fresh_witness", 6, "близко к Arch (rolling)")
	case strings.Contains(distroAll, "Manjaro"):
		Add("fresh_witness", 5, "Manjaro (rolling)"); Add("minimalist", 4, "удобный Arch")
	case strings.Contains(distroAll, "Arch"):
		Add("old_hacker", 8, "Arch Linux"); Add("fresh_witness", 6, "rolling-release")
	case strings.Contains(distroAll, "Bazzite"):
		Add("atomic", 12, "Bazzite (атомарная)"); Add("gamer", 12, "Bazzite")
	case strings.Contains(distroAll, "Silverblue") || strings.Contains(distroAll, "Kinoite") || strings.Contains(distroAll, "Sericea"):
		Add("atomic", 14, "атомарная Fedora"); Add("fresh_witness", 8, "immutable")
	case strings.Contains(distroAll, "MicroOS") || strings.Contains(distroAll, "Aeon") || strings.Contains(distroAll, "Kalpa"):
		Add("atomic", 14, "openSUSE MicroOS/Aeon"); Add("devops", 5, "transactional-update")
	case strings.Contains(distroAll, "Vanilla OS") || strings.Contains(distroAll, "VanillaOS"):
		Add("atomic", 14, "Vanilla OS"); Add("fresh_witness", 6, "immutable")
	case strings.Contains(distroAll, "blendOS"):
		Add("atomic", 12, "blendOS"); Add("fresh_witness", 6, "мульти-дистро")
	case strings.Contains(distroAll, "GNOME OS"):
		Add("atomic", 12, "GNOME OS"); Add("fresh_witness", 6, "immutable")
	case strings.Contains(distroAll, "SteamOS"):
		Add("atomic", 8, "SteamOS"); Add("gamer", 12, "SteamOS (Steam Deck)")
	case strings.Contains(distroAll, "Endless"):
		Add("atomic", 8, "Endless OS"); Add("minimalist", 4, "из коробки")
	case strings.Contains(distroAll, "Nobara"):
		Add("gamer", 12, "Nobara"); Add("fresh_witness", 6, "Fedora для игр")
	case strings.Contains(distroAll, "Fedora"):
		Add("fresh_witness", 12, "Fedora")
	case strings.Contains(distroAll, "Pop!_OS"):
		Add("gamer", 10, "Pop!_OS")
	case strings.Contains(distroAll, "elementary"):
		Add("ricer", 12, "elementary OS")
	case strings.Contains(distroAll, "Zorin"):
		Add("ricer", 8, "Zorin OS")
	case strings.Contains(distroAll, "Mint"):
		Add("minimalist", 6, "Linux Mint"); Add("sysadmin", 4, "консерватизм")
	case strings.Contains(distroAll, "MX"):
		Add("old_hacker", 12, "MX Linux"); Add("sysadmin", 4, "antiX-корни")
	case strings.Contains(distroAll, "Ubuntu"):
		Add("devops", 4, "стандарт индустрии")
	case strings.Contains(distroAll, "Debian"):
		Add("old_hacker", 10, "Debian"); Add("sysadmin", 10, "стабильность серверов")
	case strings.Contains(distroAll, "openSUSE"):
		Add("sysadmin", 10, "openSUSE/YaST"); Add("devops", 6, "корпоративный баланс")
	case strings.Contains(distroAll, "NixOS"):
		Add("atomic", 14, "NixOS (декларативная)"); Add("fresh_witness", 8, "NixOS"); Add("old_hacker", 4, "тинкеринг")
	case strings.Contains(distroAll, "Gentoo"):
		Add("old_hacker", 14, "Gentoo")
	case strings.Contains(distroAll, "Slackware"):
		Add("old_hacker", 16, "Slackware"); Add("minimalist", 5, "классика")
	case strings.Contains(distroAll, "Void"):
		Add("old_hacker", 12, "Void Linux"); Add("minimalist", 8, "runit")
	case strings.Contains(distroAll, "Alpine"):
		Add("minimalist", 15, "Alpine"); Add("sysadmin", 6, "musl + busybox"); Add("anonymous", 4, "малая поверхность атаки")
	case strings.Contains(distroAll, "Red Hat Enterprise") || strings.Contains(distroAll, "RHEL") || strings.Contains(distroAll, "Rocky") || strings.Contains(distroAll, "AlmaLinux") || strings.Contains(distroAll, "CentOS"):
		Add("sysadmin", 8, "enterprise-дистрибутив"); Add("devops", 4, "корпоративный стек")
		if strings.Contains(distroAll, "Red Hat Enterprise") || strings.Contains(distroAll, "RHEL") {
			MetaCorporat = true; Reasons["corporat"] = "Red Hat Enterprise Linux"
		}
	default:
		Add("old_hacker", 4, "нестандартный дистрибутив")
	}
}

func detectDesktop() {
	desktop := os.Getenv("XDG_CURRENT_DESKTOP") + "|" + os.Getenv("DESKTOP_SESSION") + "|" + os.Getenv("XDG_SESSION_DESKTOP")
	switch {
	case strings.Contains(desktop, "KDE") || strings.Contains(desktop, "plasma"):
		Add("ricer", 10, "KDE Plasma"); Add("devops", 4, "гибкое окружение")
	case strings.Contains(desktop, "GNOME") || strings.Contains(desktop, "gnome"):
		Add("minimalist", 8, "GNOME"); Add("ricer", 4, "цельный дизайн")
	case strings.Contains(desktop, "Cinnamon"):
		Add("minimalist", 6, "Cinnamon")
	case strings.Contains(desktop, "MATE"):
		Add("old_hacker", 10, "MATE")
	case strings.Contains(desktop, "XFCE"):
		Add("minimalist", 8, "XFCE"); Add("old_hacker", 5, "лёгкость")
	case strings.Contains(desktop, "LXQt") || strings.Contains(desktop, "LXDE"):
		Add("minimalist", 10, "LXQt/LXDE")
	case strings.Contains(desktop, "Budgie"):
		Add("ricer", 8, "Budgie")
	case strings.Contains(desktop, "Pantheon"):
		Add("ricer", 10, "Pantheon")
	case strings.Contains(desktop, "Deepin"):
		Add("ricer", 10, "Deepin")
	case strings.Contains(desktop, "COSMIC"):
		Add("fresh_witness", 12, "COSMIC"); Add("programmer", 4, "окружение на Rust")
	}
	switch {
	case strings.Contains(desktop, "Hyprland"):
		Add("ricer", 12, "Hyprland"); Add("fresh_witness", 8, "ricing на Wayland")
	case strings.Contains(desktop, "niri"):
		Add("import_substituted", 5, "niri (отеч. разработка)"); Add("fresh_witness", 10, "скроллируемый WM на Wayland")
	case strings.Contains(desktop, "sway"):
		Add("minimalist", 12, "sway"); Add("fresh_witness", 5, "Wayland"); Add("old_hacker", 3, "конфиг как код")
	case strings.Contains(desktop, "river"):
		Add("fresh_witness", 12, "river (Wayland-WM)")
	case strings.Contains(desktop, "i3"):
		Add("minimalist", 12, "i3"); Add("old_hacker", 5, "тайлинг")
	case strings.Contains(desktop, "bspwm"):
		Add("minimalist", 10, "bspwm"); Add("programmer", 5, "скриптуемый WM")
	case strings.Contains(desktop, "dwm"):
		Add("old_hacker", 14, "dwm (suckless)"); Add("minimalist", 6, "патчи и пересборка")
	case strings.Contains(desktop, "awesome"):
		Add("programmer", 8, "awesome (конфиг на Lua)"); Add("old_hacker", 4, "тайлинг")
	case strings.Contains(desktop, "qtile"):
		Add("programmer", 10, "qtile (конфиг на Python)")
	case strings.Contains(desktop, "xmonad"):
		Add("programmer", 12, "xmonad (конфиг на Haskell)"); Add("old_hacker", 4, "тайлинг")
	}
	if !strings.Contains(desktop, "niri") && Has("niri") {
		Add("import_substituted", 5, "niri установлен (отеч. разработка)")
		Add("fresh_witness", 6, "скроллируемый WM на Wayland")
	}
	if desktop == "||" {
		Add("old_hacker", 6, "без графического окружения")
		Add("sysadmin", 6, "headless-режим")
	}
}

func detectSession() {
	session := os.Getenv("XDG_SESSION_TYPE")
	switch session {
	case "wayland":
		Add("fresh_witness", 10, "Wayland")
	case "x11":
		Add("old_hacker", 5, "X11")
	case "tty", "":
		Add("old_hacker", 5, "чистый TTY"); Add("sysadmin", 3, "без иксов")
	}
}

func detectHardware() {
	if d, err := os.ReadFile("/proc/meminfo"); err == nil {
		for _, line := range strings.Split(string(d), "\n") {
			if strings.HasPrefix(line, "MemTotal:") {
				fields := strings.Fields(line)
				if len(fields) >= 2 {
					ramKB := 0
					fmt.Sscanf(fields[1], "%d", &ramKB)
					ramGB := ramKB / 1024 / 1024
					if ramGB >= 64 {
						Add("devops", 5, "64+ GB RAM"); Add("gamer", 6, "флагманское железо")
					} else if ramGB >= 32 {
						Add("devops", 3, "32+ GB RAM"); Add("gamer", 5, "мощное железо")
					} else if ramGB <= 4 {
						Add("minimalist", 10, "≤4 GB RAM")
					} else if ramGB <= 8 {
						Add("minimalist", 5, "8 GB RAM")
					}
				}
				break
			}
		}
	}
	if runtime.NumCPU() >= 16 {
		Add("gamer", 4, "16+ ядер (мощная станция)")
	}
}

func detectDisk() {
	out, _ := exec.Command("df", "/").Output()
	lines := strings.Split(string(out), "\n")
	if len(lines) >= 2 {
		fields := strings.Fields(lines[1])
		if len(fields) >= 5 {
			pct := strings.TrimSuffix(fields[4], "%")
			used := 0
			fmt.Sscanf(pct, "%d", &used)
			if used > 85 {
				Add("old_hacker", 6, "диск почти полон")
			} else if used < 30 {
				Add("minimalist", 5, "чистота файловой системы")
			}
		}
	}
	if _, err := os.Stat("/etc/crypttab"); err == nil {
		Add("anonymous", 10, "шифрование диска (LUKS)")
	}
	if _, err := os.Stat("/proc/mdstat"); err == nil {
		d, _ := os.ReadFile("/proc/mdstat")
		if strings.Contains(string(d), "^md") && Has("mdadm") {
			Add("sysadmin", 6, "программный RAID (mdadm)")
		}
	}
	if Has("lvs") {
		if _, err := exec.Command("lvs").Output(); err == nil {
			Add("sysadmin", 5, "LVM")
		}
	}
}

func detectInit() {
	if _, err := os.Stat("/run/systemd/system"); err == nil {
		return
	}
	switch {
	case fileExists("/run/runit.stopit") || dirExists("/etc/sv"):
		Add("old_hacker", 12, "runit")
	case dirExists("/etc/s6") || Has("s6-svscan"):
		Add("old_hacker", 12, "s6")
	case fileExists("/etc/dinit") || Has("dinit"):
		Add("old_hacker", 10, "dinit")
	case Has("openrc") || fileExists("/etc/init.d/openrc"):
		Add("old_hacker", 10, "OpenRC")
	default:
		Add("old_hacker", 8, "альтернативный init")
	}
}

func detectBootloader() {
	bl := ""
	bootctlOut := ""
	if Has("bootctl") {
		if out, err := exec.Command("bootctl", "status").CombinedOutput(); err == nil {
			bootctlOut = string(out)
		}
	}
	switch {
	case fileExists("/etc/lilo.conf") || Has("lilo"):
		bl = "lilo"
	case dirExists("/boot/efi/EFI/refind") || dirExists("/boot/EFI/refind") || fileExists("/boot/refind_linux.conf") || Has("refind-install") || strings.Contains(bootctlOut, "rEFInd"):
		bl = "refind"
	case fileExists("/boot/limine.cfg") || fileExists("/boot/limine.conf") || fileExists("/boot/limine/limine.conf") || Has("limine"):
		bl = "limine"
	case strings.Contains(bootctlOut, "systemd-boot") || dirExists("/boot/efi/EFI/systemd") || dirExists("/boot/EFI/systemd"):
		bl = "systemd-boot"
	case dirExists("/boot/syslinux") || dirExists("/boot/extlinux") || Has("extlinux"):
		bl = "syslinux"
	case Has("efibootmgr"):
		if out, err := exec.Command("efibootmgr", "-v").CombinedOutput(); err == nil {
			s := strings.ToLower(string(out))
			if strings.Contains(s, "vmlinuz") || strings.Contains(s, "linux.efi") || strings.Contains(s, "efistub") {
				bl = "efistub"
			}
		}
	case dirExists("/boot/grub") || dirExists("/boot/grub2") || Has("grub-install") || Has("grub2-install") || Has("update-grub") || strings.Contains(bootctlOut, "GRUB"):
		bl = "grub"
	}
	switch bl {
	case "lilo":
		Add("old_hacker", 12, "LILO (олдскул-загрузчик)")
	case "refind":
		Add("ricer", 8, "rEFInd (красивый загрузчик)")
	case "limine":
		Add("fresh_witness", 8, "Limine"); Add("old_hacker", 3, "осознанный выбор загрузчика")
	case "systemd-boot":
		Add("minimalist", 6, "systemd-boot"); Add("fresh_witness", 3, "минимальный EFI-загрузчик")
	case "syslinux":
		Add("old_hacker", 6, "syslinux/extlinux"); Add("minimalist", 3, "лёгкий загрузчик")
	case "efistub":
		Add("minimalist", 8, "EFISTUB (без загрузчика)"); Add("old_hacker", 6, "прямая загрузка ядра")
	case "grub":
		if d, err := os.ReadFile("/etc/default/grub"); err == nil && strings.Contains(string(d), "GRUB_THEME=") {
			Add("ricer", 5, "кастомная тема GRUB")
		}
	}
}

func detectLegacy() {
	kver := ""
	if out, err := exec.Command("uname", "-r").Output(); err == nil {
		kver = strings.TrimSpace(string(out))
	}
	parts := strings.Split(kver, ".")
	if len(parts) > 0 {
		major := 0
		fmt.Sscanf(parts[0], "%d", &major)
		if major > 0 && major <= 4 {
			Add("legacy", 12, fmt.Sprintf("старое ядро %d.x", major))
		} else if major == 5 {
			Add("legacy", 4, "ядро 5.x (не самое свежее)")
		}
	}
	is32 := false
	arch, _ := exec.Command("uname", "-m").Output()
	switch strings.TrimSpace(string(arch)) {
	case "i386", "i486", "i586", "i686":
		is32 = true
	}
	if !is32 {
		// 32-битный CPU при нестандартном uname (напр. 32-битный userland).
		if d, err := os.ReadFile("/proc/cpuinfo"); err == nil {
			for _, m := range []string{"i386", "i486", "i586", "i686"} {
				if strings.Contains(string(d), m) { is32 = true; break }
			}
		}
	}
	if is32 {
		Add("legacy", 12, "32-битная система")
	}
	session := os.Getenv("XDG_SESSION_TYPE")
	if session == "x11" {
		Add("legacy", 4, "сессия X11/Xorg")
	}
	desktop := os.Getenv("XDG_CURRENT_DESKTOP") + "|" + os.Getenv("DESKTOP_SESSION") + "|" + os.Getenv("XDG_SESSION_DESKTOP")
	switch {
	case strings.Contains(desktop, "Trinity") || strings.Contains(desktop, "TDE"):
		Add("legacy", 12, "Trinity (KDE 3)")
	case strings.Contains(desktop, "MATE"):
		Add("legacy", 6, "MATE (наследник GNOME 2)")
	case strings.Contains(desktop, "LXDE"):
		Add("legacy", 6, "LXDE")
	case strings.Contains(desktop, "fvwm") || strings.Contains(desktop, "FVWM"):
		Add("legacy", 12, "FVWM")
	case strings.Contains(desktop, "windowmaker") || strings.Contains(desktop, "WindowMaker"):
		Add("legacy", 10, "Window Maker")
	case strings.Contains(desktop, "icewm") || strings.Contains(desktop, "IceWM"):
		Add("legacy", 10, "IceWM")
	case strings.Contains(desktop, "twm") || strings.Contains(desktop, "cwm"):
		Add("legacy", 10, "классический WM (twm/cwm)")
	case strings.Contains(desktop, "fluxbox") || strings.Contains(desktop, "blackbox"):
		Add("legacy", 8, "Fluxbox/Blackbox")
	case strings.Contains(desktop, "jwm") || strings.Contains(desktop, "JWM"):
		Add("legacy", 8, "JWM")
	case strings.Contains(desktop, "enlightenment") || strings.Contains(desktop, "Enlightenment"):
		Add("legacy", 6, "Enlightenment")
	case strings.Contains(desktop, "openbox") || strings.Contains(desktop, "Openbox"):
		Add("legacy", 5, "Openbox")
	}
	if Has("mplayer") && !Has("mpv") {
		Add("legacy", 5, "MPlayer (без mpv)")
	}
	if Has("pidgin") {
		Add("legacy", 5, "Pidgin")
	}
	if Has("xmms") || Has("audacious") {
		Add("legacy", 4, "XMMS/Audacious")
	}
}

func detectAnonymous() {
	if HasUsed("gpg") { Add("anonymous", 6, "GPG") }
	if HasUsed("pass") { Add("anonymous", 8, "pass") }
	if HasUsed("age") { Add("anonymous", 6, "age") }
	if Has("veracrypt") { Add("anonymous", 14, "VeraCrypt") }
	if HasUsed("cryptsetup") { Add("anonymous", 5, "cryptsetup") }
	if HasUsed("tomb") { Add("anonymous", 8, "Tomb") }
	if Has("keepassxc") { Add("anonymous", 5, "KeePassXC") }
	if HasUsed("gocryptfs") { Add("anonymous", 6, "gocryptfs") }
	if HasUsed("encfs") { Add("anonymous", 5, "EncFS") }
	if Has("tor") { Add("anonymous", 12, "Tor") }
	if Has("torbrowser-launcher") { Add("anonymous", 10, "Tor Browser") }
	if Has("i2prouter") { Add("anonymous", 10, "I2P") }
	if HasUsed("proxychains") { Add("anonymous", 8, "proxychains"); Add("pentester", 4, "цепочки прокси") }
	if Has("mullvad") { Add("anonymous", 10, "Mullvad VPN") }
	if Has("protonvpn") || Has("protonvpn-cli") { Add("anonymous", 8, "ProtonVPN") }
	if HasUsed("openvpn") { Add("anonymous", 6, "OpenVPN") }
	if HasUsed("wg") || HasUsed("wg-quick") { Add("anonymous", 6, "WireGuard") }
}

func detectPentester() {
	if HasUsed("nmap") { Add("pentester", 12, "nmap") }
	if HasUsed("masscan") { Add("pentester", 8, "masscan") }
	if Has("wireshark") { Add("pentester", 8, "Wireshark") }
	if HasUsed("tshark") { Add("pentester", 6, "tshark") }
	if HasUsed("tcpdump") { Add("pentester", 5, "tcpdump") }
	if HasUsed("msfconsole") { Add("pentester", 16, "Metasploit") }
	if HasUsed("aircrack-ng") { Add("pentester", 12, "aircrack-ng") }
	if HasUsed("hashcat") { Add("pentester", 10, "hashcat") }
	if HasUsed("john") { Add("pentester", 10, "John the Ripper") }
	if HasUsed("hydra") { Add("pentester", 10, "hydra") }
	if HasUsed("sqlmap") { Add("pentester", 10, "sqlmap") }
	if HasUsed("nikto") { Add("pentester", 8, "nikto") }
	if HasUsed("gobuster") { Add("pentester", 6, "gobuster") }
	if HasUsed("ffuf") { Add("pentester", 6, "ffuf") }
	if Has("burpsuite") { Add("pentester", 12, "Burp Suite") }
	if Has("zaproxy") { Add("pentester", 8, "OWASP ZAP") }
	if HasUsed("radare2") || HasUsed("r2") { Add("pentester", 10, "radare2") }
	if Has("ghidra") { Add("pentester", 12, "Ghidra") }
	if HasUsed("binwalk") { Add("pentester", 6, "binwalk") }
	if HasUsed("volatility") { Add("pentester", 8, "Volatility") }
	if HasUsed("wpscan") { Add("pentester", 6, "WPScan") }
	if Has("fail2ban") { Add("sysadmin", 6, "fail2ban") }
}

func detectDevops() {
	if HasUsed("docker") { Add("devops", 5, "Docker") }
	if HasUsed("podman") { Add("devops", 5, "Podman"); Add("anonymous", 3, "rootless-контейнеры") }
	if HasUsed("kubectl") { Add("devops", 14, "Kubernetes") }
	if HasUsed("k9s") { Add("devops", 8, "k9s") }
	if HasUsed("helm") { Add("devops", 8, "Helm") }
	if Has("minikube") || Has("kind") { Add("devops", 6, "локальный кластер") }
	if Has("argocd") || Has("flux") || Has("skaffold") { Add("devops", 8, "GitOps") }
	if HasUsed("terraform") { Add("devops", 12, "Terraform") }
	if HasUsed("opentofu") { Add("devops", 10, "OpenTofu") }
	if HasUsed("pulumi") { Add("devops", 8, "Pulumi") }
	if HasUsed("ansible") { Add("devops", 8, "Ansible"); Add("sysadmin", 4, "автоматизация") }
	if Has("puppet") || Has("chef") || Has("salt") { Add("devops", 8, "config management") }
	if HasUsed("vagrant") { Add("devops", 6, "Vagrant") }
	if Has("qemu-img") || Has("virt-manager") || Has("virsh") { Add("devops", 8, "QEMU/KVM") }
	if Has("lxc") || Has("lxd") || Has("incus") { Add("devops", 6, "LXC/Incus") }
	if Has("vault") || Has("consul") || Has("nomad") { Add("devops", 8, "HashiCorp-стек") }
	if HasUsed("aws") || HasUsed("gcloud") || HasUsed("az") || HasUsed("doctl") { Add("devops", 8, "облачный CLI") }
	if Has("gitlab-runner") || Has("act") { Add("devops", 5, "CI-раннеры") }
}

func detectSysadmin() {
	if Has("nginx") {
		if Used("nginx") {
			Add("sysadmin", 8, "nginx"); Add("import_substituted", 4, "nginx (Игорь Сысоев)")
		}
	}
	if Has("apache2") || Has("httpd") { Add("sysadmin", 6, "Apache"); Add("old_hacker", 4, "httpd") }
	if Has("psql") || Has("postgres") || dirExists("/var/lib/pgsql") || dirExists("/var/lib/postgresql") {
		Add("sysadmin", 8, "PostgreSQL"); Add("import_substituted", 4, "PostgreSQL (Postgres Pro)")
	}
	if Has("mysql") || Has("mariadb") { Add("sysadmin", 6, "MySQL/MariaDB") }
	if Has("redis-cli") { Add("sysadmin", 5, "Redis") }
	if Has("mongod") { Add("devops", 5, "MongoDB") }
	if Has("sshd") || fileExists("/etc/ssh/sshd_config") { Add("sysadmin", 5, "SSH-сервер") }
	if Has("htop") || Has("btop") || Has("glances") { Add("sysadmin", 3, "мониторинг процессов") }
	if Has("prometheus") { Add("sysadmin", 6, "Prometheus"); Add("devops", 3, "метрики") }
	if Has("grafana") { Add("sysadmin", 6, "Grafana"); Add("devops", 3, "дашборды") }
	if Has("borg") || Has("restic") || Has("rsnapshot") || Has("duplicity") { Add("sysadmin", 8, "бэкапы (borg/restic)") }
	if Has("smartctl") { Add("sysadmin", 3, "S.M.A.R.T.-мониторинг") }
}

func detectProgrammer() {
	devCount := 0
	for _, lang := range []string{"node", "deno", "bun", "rustc", "go", "zig", "elixir", "julia", "haskell", "scala", "kotlin", "crystal", "nim", "ocaml"} {
		if Has(lang) { Add("programmer", 3, lang); devCount++ }
	}
	if Has("ruby") { Add("programmer", 1, "ruby"); devCount++ }
	home := os.Getenv("HOME")
	if home == "" { home = "/root" }
	if Has("python3") || Has("python") {
		if Has("pyenv") || Has("poetry") || Has("pipx") || Has("virtualenv") || dirExists(home+"/.virtualenvs") || hasPythonSitePackages(home) {
			Add("programmer", 3, "Python + инструменты разработки"); devCount++
		}
	}
	if Has("javac") || Has("mvn") || Has("gradle") || Has("jenv") || dirExists(home+"/.m2") || dirExists(home+"/.gradle") || dirExists(home+"/.sdkman") {
		Add("programmer", 3, "Java (JDK/сборка)"); devCount++
	}
	if dirExists(home + "/.cargo") { Add("programmer", 5, "Cargo"); Add("fresh_witness", 2, "Rust toolchain") }
	if dirExists(home + "/.rustup") { Add("programmer", 2, "rustup") }
	if dirExists(home + "/go") { Add("programmer", 3, "Go workspace") }
	if dirExists(home + "/.npm") { Add("programmer", 3, "npm-проекты") }
	if Has("asdf") || Has("nvm") || Has("pyenv") || Has("rbenv") || Has("sdk") { Add("programmer", 5, "менеджер версий языков") }
	if _, err := os.Stat(home + "/.gitconfig"); err == nil {
		d, _ := os.ReadFile(home + "/.gitconfig")
		if strings.Contains(strings.ToLower(string(d)), "[user]") { Add("programmer", 6, "настроенный git (user.*)") }
	}
	if devCount >= 5 { Add("programmer", 10, "полиглот (5+ языков)")
	} else if devCount >= 3 { Add("programmer", 4, "несколько языков") }
	if Has("gcc") || Has("clang") { Add("old_hacker", 3, "gcc/clang") }
	if Has("make") && Has("cmake") { Add("programmer", 3, "сборочные системы") }
	if Has("rustc") { Add("fresh_witness", 6, "Rust") }
	if Has("docker-compose") { Add("devops", 4, "Compose") }
	if Has("vim") || Has("nvim") { Add("old_hacker", 6, "Vim/Neovim"); Add("programmer", 4, "модальное редактирование") }
	if Has("emacs") { Add("old_hacker", 12, "Emacs") }
	if Has("helix") { Add("fresh_witness", 6, "Helix") }
	if Has("code") { Add("programmer", 6, "VS Code") }
	if Has("zed") || Has("zeditor") { Add("programmer", 6, "Zed"); Add("fresh_witness", 3, "редактор на Rust") }
	if Has("codium") { Add("programmer", 6, "VSCodium") }
	if Has("cursor") { Add("programmer", 6, "Cursor") }
	// AI coding assistants
	if Has("aider") { Add("programmer", 5, "aider (AI pair programming)") }
	if Has("codex") { Add("programmer", 5, "Codex (OpenAI)") }
	if Has("claude") || Has("claude-code") { Add("programmer", 5, "Claude Code") }
	if Has("continue") { Add("programmer", 5, "Continue (AI assistant)") }
	if Has("copilot") { Add("programmer", 5, "GitHub Copilot CLI") }
	if Has("cody") { Add("programmer", 5, "Cody (Sourcegraph)") }
	if Has("phind") { Add("programmer", 4, "Phind") }
	if Has("tabnine") { Add("programmer", 4, "Tabnine") }
	if Has("supermaven") { Add("programmer", 4, "Supermaven") }
	if Has("codeium") { Add("programmer", 5, "Codeium") }
	if Has("subl") { Add("programmer", 5, "Sublime Text") }
	if Has("lapce") { Add("programmer", 5, "Lapce"); Add("fresh_witness", 3, "редактор на Rust") }
	if Has("micro") { Add("programmer", 3, "micro") }
	if Has("nvim") && dirExists(home+"/.config/nvim") { Add("ricer", 5, "кастомный Neovim") }
	if dirExists(home+"/.config/JetBrains") || Has("idea") { Add("programmer", 6, "JetBrains IDE") }
	if Has("tmux") || Has("screen") { Add("old_hacker", 6, "мультиплексор") }
	if Has("gdb") || Has("lldb") { Add("programmer", 4, "отладчики") }

	// AI plugins in IDEs
	// VS Code extensions
	vibeScore := 0
	if dirExists(home + "/.vscode/extensions") {
		entries, _ := os.ReadDir(home + "/.vscode/extensions")
		for _, e := range entries {
			name := strings.ToLower(e.Name())
			if strings.Contains(name, "copilot") || strings.Contains(name, "codeium") ||
				strings.Contains(name, "cody") || strings.Contains(name, "tabnine") ||
				strings.Contains(name, "supermaven") || strings.Contains(name, "continue") ||
				strings.Contains(name, "phind") || strings.Contains(name, "aider") {
				vibeScore += 3
				break
			}
		}
	}
	// JetBrains plugins
	if dirExists(home + "/.config/JetBrains") {
		entries, _ := os.ReadDir(home + "/.config/JetBrains")
		for _, e := range entries {
			pluginDir := home + "/.config/JetBrains/" + e.Name() + "/plugins"
			if dirExists(pluginDir) {
				plugins, _ := os.ReadDir(pluginDir)
				for _, p := range plugins {
					name := strings.ToLower(p.Name())
					if strings.Contains(name, "copilot") || strings.Contains(name, "codeium") ||
						strings.Contains(name, "cody") || strings.Contains(name, "tabnine") {
						vibeScore += 3
						break
					}
				}
			}
		}
	}
	// Neovim plugins (check lazy.nvim or packer)
	if dirExists(home + "/.local/share/nvim/lazy") {
		entries, _ := os.ReadDir(home + "/.local/share/nvim/lazy")
		for _, e := range entries {
			name := strings.ToLower(e.Name())
			if strings.Contains(name, "copilot") || strings.Contains(name, "codeium") ||
				strings.Contains(name, "cody") || strings.Contains(name, "tabnine") ||
				strings.Contains(name, "copilot-cmp") || strings.Contains(name, "avante") ||
				strings.Contains(name, "codecompanion") || strings.Contains(name, "gen.nvim") {
				vibeScore += 3
				break
			}
		}
	}
	// AI CLI tools
	if Has("aider") { vibeScore += 5 }
	if Has("codex") { vibeScore += 5 }
	if Has("claude") || Has("claude-code") { vibeScore += 5 }
	if Has("continue") { vibeScore += 4 }
	if Has("copilot") { vibeScore += 4 }
	if Has("cody") { vibeScore += 4 }
	if Has("phind") { vibeScore += 3 }
	if Has("tabnine") { vibeScore += 3 }
	if Has("supermaven") { vibeScore += 3 }
	if Has("codeium") { vibeScore += 4 }

	if vibeScore >= 10 {
		Add("vibe_coder", vibeScore, fmt.Sprintf("AI-инструменты (score: %d)", vibeScore))
	} else if vibeScore >= 5 {
		Add("vibe_coder", vibeScore, fmt.Sprintf("AI-инструменты (score: %d)", vibeScore))
	}
}

func detectBrowser() {
	if Has("yandex-browser") || Has("yandex_browser") || Has("yandex-browser-stable") {
		Add("import_substituted", 0, "Яндекс.Браузер")
	}
	locale := os.Getenv("LANG") + "|" + os.Getenv("LC_ALL") + "|" + os.Getenv("LC_CTYPE") + "|" + os.Getenv("LC_MESSAGES")
	if strings.Contains(locale, "ru_RU") || strings.Contains(locale, "ru_") {
		Add("import_substituted", 2, "русская локаль")
	}
	if Has("librewolf") { Add("anonymous", 6, "LibreWolf"); Add("ricer", 3, "приватный форк") }
	if Has("brave") { Add("anonymous", 4, "Brave") }
	if Has("mullvad-browser") { Add("anonymous", 8, "Mullvad Browser") }
}

func detectPackageManagers() {
	if Has("flatpak") { Add("atomic", 4, "Flatpak") }
	if Has("snap") { Add("minimalist", 2, "Snap") }
	if Has("nix") { Add("atomic", 8, "Nix"); Add("fresh_witness", 4, "декларативные пакеты") }
	if Has("brew") { Add("programmer", 4, "Homebrew") }
	if Has("distrobox") || Has("toolbox") || Has("toolbx") { Add("atomic", 8, "distrobox/toolbox") }
	if Has("rpm-ostree") || Has("ostree") || Has("bootc") { Add("atomic", 10, "ostree-система") }
	if Has("yay") || Has("paru") { Add("old_hacker", 4, "AUR-хелпер") }
	if Has("guix") { Add("atomic", 8, "GNU Guix"); Add("old_hacker", 6, "функциональный пакетинг") }
}

func detectGamer() {
	if Has("steam") { Add("gamer", 12, "Steam") }
	if Has("lutris") { Add("gamer", 8, "Lutris") }
	if Has("heroic") { Add("gamer", 6, "Heroic") }
	if Has("bottles") { Add("gamer", 6, "Bottles") }
	if Has("wine") || Has("wine64") { Add("gamer", 6, "Wine") }
	if Has("gamemoderun") { Add("gamer", 6, "GameMode") }
	if Has("mangohud") { Add("gamer", 5, "MangoHud") }
	if Has("protontricks") || Has("protontricks-launch") { Add("gamer", 5, "Proton") }
	if Has("retroarch") { Add("gamer", 6, "RetroArch (эмуляция)") }
	gameRe := regexp.MustCompile(`(?i)^Categories=.*Game`)
	gameCount := 0
	for _, dir := range []string{"/usr/share/applications", os.Getenv("HOME") + "/.local/share/applications", "/var/lib/flatpak/exports/share/applications", os.Getenv("HOME") + "/.local/share/flatpak/exports/share/applications"} {
		if entries, err := os.ReadDir(dir); err == nil {
			for _, e := range entries {
				if strings.HasSuffix(e.Name(), ".desktop") {
					if d, err := os.ReadFile(filepath.Join(dir, e.Name())); err == nil {
						for _, line := range strings.Split(string(d), "\n") {
							if gameRe.MatchString(line) {
								gameCount++
								break
							}
						}
					}
				}
			}
		}
	}
	if gameCount >= 15 { Add("gamer", 14, fmt.Sprintf("%d игр в меню", gameCount))
	} else if gameCount >= 5 { Add("gamer", 8, fmt.Sprintf("%d игр в меню", gameCount))
	} else if gameCount >= 1 { Add("gamer", 3, fmt.Sprintf("%d игр в меню", gameCount)) }
	if d, err := os.ReadFile("/proc/modules"); err == nil && strings.Contains(string(d), "nvidia") {
		Add("gamer", 6, "NVIDIA GPU")
	}
	home := os.Getenv("HOME")
	if home == "" { home = "/root" }
	for _, gd := range []string{home + "/Games", home + "/games", home + "/Игры", home + "/игры"} {
		if info, err := os.Stat(gd); err == nil && info.IsDir() {
			entries, _ := os.ReadDir(gd)
			count := len(entries)
			if count >= 5 { Add("gamer", 12, fmt.Sprintf("каталог %s (%d шт.)", filepath.Base(gd), count))
			} else if count >= 1 { Add("gamer", 8, fmt.Sprintf("каталог %s (%d шт.)", filepath.Base(gd), count)) }
			break
		}
	}
}

func detectShell() {
	shell := os.Getenv("SHELL")
	if shell == "" { shell = "bash" }
	shell = filepath.Base(shell)
	switch shell {
	case "zsh": Add("ricer", 8, "zsh")
	case "fish": Add("fresh_witness", 6, "fish")
	case "bash": Add("old_hacker", 3, "bash")
	case "dash": Add("minimalist", 8, "dash")
	case "nu": Add("fresh_witness", 10, "nushell"); Add("programmer", 3, "структурный shell")
	case "xonsh": Add("fresh_witness", 5, "xonsh"); Add("programmer", 4, "Python-shell")
	case "elvish": Add("fresh_witness", 6, "elvish")
	}
}

func detectTerminal() {
	if Has("kitty") { Add("ricer", 6, "kitty"); Add("fresh_witness", 3, "GPU-терминал") }
	if Has("alacritty") { Add("ricer", 5, "Alacritty"); Add("fresh_witness", 3, "GPU-терминал") }
	if Has("wezterm") { Add("ricer", 6, "WezTerm"); Add("programmer", 3, "конфиг на Lua") }
	if Has("foot") { Add("minimalist", 5, "foot"); Add("fresh_witness", 3, "Wayland-терминал") }
	if Has("st") { Add("old_hacker", 6, "st (suckless)"); Add("minimalist", 4, "патчи под себя") }
	if Has("xterm") { Add("old_hacker", 3, "xterm") }
}

func detectDotfiles() {
	home := os.Getenv("HOME")
	if home == "" { home = "/root" }
	riceConfigs := 0
	for _, cfg := range []string{"hypr", "waybar", "polybar", "picom", "compton", "rofi", "wofi", "dunst", "mako", "eww", "sway", "i3", "bspwm", "awesome", "qtile", "river", "niri", "kitty", "alacritty", "wezterm", "foot", "fastfetch", "neofetch", "wal", "swaylock", "wlogout", "starship"} {
		if _, err := os.Stat(home + "/.config/" + cfg); err == nil { riceConfigs++ }
	}
	if riceConfigs >= 6 { Add("ricer", 12, fmt.Sprintf("кастомных конфигов: %d", riceConfigs))
	} else if riceConfigs >= 3 { Add("ricer", 6, fmt.Sprintf("кастомизация окружения (%d)", riceConfigs))
	} else if riceConfigs == 0 { Add("minimalist", 5, "без кастомизации окружения") }
	if dirExists(home+"/dotfiles/.git") || dirExists(home+"/.dotfiles/.git") {
		Add("ricer", 8, "dotfiles в Git"); Add("programmer", 4, "управление конфигами")
	}
	if Has("stow") || Has("chezmoi") { Add("ricer", 5, "менеджер dotfiles"); Add("devops", 3, "воспроизводимые конфиги") }
}

func detectSystemState() {
	if d, err := os.ReadFile("/proc/uptime"); err == nil {
		fields := strings.Fields(string(d))
		if len(fields) > 0 {
			sec := 0
			fmt.Sscanf(fields[0], "%d", &sec)
			days := sec / 86400
			if days >= 30 { Add("sysadmin", 10, fmt.Sprintf("аптайм %d дн.", days))
			} else if days >= 7 { Add("sysadmin", 5, fmt.Sprintf("аптайм %d дн.", days)) }
		}
	}
	installEpoch := 0
	if out, err := exec.Command("stat", "-c", "%W", "/").Output(); err == nil {
		fmt.Sscanf(strings.TrimSpace(string(out)), "%d", &installEpoch)
	}
	if installEpoch == 0 {
		if out, err := exec.Command("stat", "-c", "%Y", "/etc/machine-id").Output(); err == nil {
			fmt.Sscanf(strings.TrimSpace(string(out)), "%d", &installEpoch)
		}
	}
	if installEpoch > 0 {
		ageDays := (int(time.Now().Unix()) - installEpoch) / 86400
		if ageDays >= 1095 {
			Add("sysadmin", 8, fmt.Sprintf("система живёт %d дн. без переустановки", ageDays))
			Add("old_hacker", 4, "не распыляется на переустановки")
		} else if ageDays >= 365 {
			Add("sysadmin", 4, "установлена больше года назад")
		} else if ageDays >= 0 && ageDays <= 14 {
			Add("fresh_witness", 4, fmt.Sprintf("свежая установка (%d дн.)", ageDays))
		}
	}
	pkgCount := 0
	switch {
	case Has("pacman"):
		if out, err := exec.Command("pacman", "-Qq").Output(); err == nil { pkgCount = strings.Count(string(out), "\n") }
	case Has("dpkg-query"):
		if out, err := exec.Command("dpkg-query", "-f", ".\n", "-W").Output(); err == nil { pkgCount = strings.Count(string(out), "\n") }
	case Has("rpm"):
		if out, err := exec.Command("rpm", "-qa").Output(); err == nil { pkgCount = strings.Count(string(out), "\n") }
	case Has("apk"):
		if out, err := exec.Command("apk", "info").Output(); err == nil { pkgCount = strings.Count(string(out), "\n") }
	case Has("xbps-query"):
		if out, err := exec.Command("xbps-query", "-l").Output(); err == nil { pkgCount = strings.Count(string(out), "\n") }
	}
	if pkgCount > 0 {
		if pkgCount <= 300 { Add("minimalist", 10, fmt.Sprintf("очень мало пакетов (%d)", pkgCount))
		} else if pkgCount <= 600 { Add("minimalist", 5, fmt.Sprintf("немного пакетов (%d)", pkgCount)) }
	}
	if Has("pacman") {
		if out, err := exec.Command("pacman", "-Qqm").Output(); err == nil {
			aurCount := strings.Count(string(out), "\n")
			if aurCount >= 20 { Add("old_hacker", 6, fmt.Sprintf("%d пакетов из AUR", aurCount))
			} else if aurCount >= 5 { Add("old_hacker", 3, "сборки из AUR") }
		}
	}
	kernel, _ := exec.Command("uname", "-r").Output()
	kernelStr := strings.TrimSpace(string(kernel))
	switch {
	case strings.Contains(kernelStr, "zen"):
		Add("gamer", 5, "ядро Zen"); Add("fresh_witness", 3, "тюнинг отзывчивости")
	case strings.Contains(kernelStr, "xanmod"):
		Add("gamer", 5, "ядро XanMod")
	case strings.Contains(kernelStr, "lqx") || strings.Contains(kernelStr, "liquorix"):
		Add("gamer", 5, "ядро Liquorix")
	case strings.Contains(kernelStr, "tkg"):
		Add("gamer", 5, "ядро TkG"); Add("old_hacker", 3, "сборка ядра")
	case strings.Contains(kernelStr, "hardened"):
		Add("anonymous", 6, "hardened-ядро"); Add("pentester", 3, "защищённое ядро")
	case strings.Contains(kernelStr, "-rt"):
		Add("sysadmin", 4, "realtime-ядро")
	case strings.Contains(kernelStr, "lts"):
		Add("sysadmin", 4, "LTS-ядро (стабильность)")
	}
	if Has("snapper") || dirExists("/.snapshots") { Add("sysadmin", 5, "снапшоты (snapper)"); Add("atomic", 4, "откаты ФС") }
	if Has("timeshift") { Add("sysadmin", 5, "Timeshift") }
	home := os.Getenv("HOME")
	if home == "" { home = "/root" }
	if dirExists("/etc/nixos") || dirExists(home+"/.config/home-manager") { Add("atomic", 6, "Nix/home-manager") }
	if Has("zfs") { Add("sysadmin", 6, "ZFS"); Add("old_hacker", 3, "ZFS-энтузиаст") }
	if dirExists(home + "/.ssh") {
		entries, _ := os.ReadDir(home + "/.ssh")
		sshKeys := 0
		for _, e := range entries {
			if strings.HasPrefix(e.Name(), "id_") && !strings.HasSuffix(e.Name(), ".pub") { sshKeys++ }
		}
		if sshKeys >= 1 { Add("devops", 4, "SSH-ключи"); Add("sysadmin", 3, "доступ к хостам") }
	}
	if dirExists(home) {
		gitRepos := 0
		if out, err := exec.Command("find", home, "-maxdepth", "4", "-name", ".git", "-type", "d").Output(); err == nil {
			gitRepos = strings.Count(string(out), "\n")
		}
		if gitRepos >= 10 { Add("programmer", 8, fmt.Sprintf("%d+ git-репозиториев", gitRepos))
		} else if gitRepos >= 3 { Add("programmer", 4, fmt.Sprintf("%d git-репозитория", gitRepos)) }
	}
	if Has("flatpak") {
		if out, err := exec.Command("flatpak", "list", "--app", "--columns=application").Output(); err == nil {
			fpApps := string(out)
			fpCount := strings.Count(fpApps, "\n")
			if fpCount >= 15 { Add("atomic", 8, fmt.Sprintf("%d flatpak-приложений", fpCount))
			} else if fpCount >= 5 { Add("atomic", 5, fmt.Sprintf("%d flatpak-приложений", fpCount)) }
			if strings.Contains(fpApps, "torproject") || strings.Contains(fpApps, "mullvad") || strings.Contains(fpApps, "signalapp") { Add("anonymous", 5, "приватные flatpak") }
			if strings.Contains(fpApps, "valvesoftware.Steam") || strings.Contains(fpApps, "heroicgameslauncher") || strings.Contains(fpApps, "net.lutris") { Add("gamer", 6, "игровые flatpak") }
			if strings.Contains(fpApps, "visualstudio") || strings.Contains(fpApps, "jetbrains") || strings.Contains(fpApps, "gnome.Builder") { Add("programmer", 5, "dev-flatpak") }
			if strings.Contains(fpApps, "blender") || strings.Contains(fpApps, "gimp") || strings.Contains(fpApps, "inkscape") || strings.Contains(fpApps, "kdenlive") || strings.Contains(fpApps, "obsproject") || strings.Contains(fpApps, "darktable") || strings.Contains(fpApps, "krita") || strings.Contains(fpApps, "Audacity") {
				Add("ricer", 4, "креатив/медиа flatpak"); Add("creative", 5, "креатив/медиа flatpak")
			}
		}
	}
	if Has("1cv8") || Has("1cv8c") || Has("1cestart") || Has("1c") || dirExists("/opt/1cv8") {
		Add("import_substituted", 8, "1С:Предприятие")
	}
	if Has("lxprofile") || os.Getenv("LXPROFILE_ROOT") != "" || fileExists(os.Args[0]) {
		Add("import_substituted", 1, "запущен lxprofiler ;)")
	}
	for _, app := range []string{"gimp", "krita", "darktable", "rawtherapee", "inkscape", "blender", "kdenlive", "shotcut", "openshot", "flowblade", "ardour", "lmms", "hydrogen", "mixxx", "audacity", "bitwig-studio", "renoise"} {
		if Has(app) { Add("creative", 4, app) }
	}
	if Has("resolve") { Add("creative", 6, "DaVinci Resolve") }
	desktop := os.Getenv("XDG_CURRENT_DESKTOP") + "|" + os.Getenv("DESKTOP_SESSION") + "|" + os.Getenv("XDG_SESSION_DESKTOP")
	if Has("gnome-shell") || strings.Contains(desktop, "GNOME") || strings.Contains(desktop, "gnome") {
		gnomeExt := 0
		if Has("gsettings") {
			if out, err := exec.Command("gsettings", "get", "org.gnome.shell", "enabled-extensions").Output(); err == nil {
				gnomeExt = strings.Count(string(out), "'") / 2
			}
		}
		if gnomeExt == 0 && dirExists(home+"/.local/share/gnome-shell/extensions") {
			entries, _ := os.ReadDir(home + "/.local/share/gnome-shell/extensions")
			gnomeExt = len(entries)
		}
		if gnomeExt >= 1 {
			if gnomeExt > 6 { gnomeExt = 6 }
			Add("ricer", gnomeExt, fmt.Sprintf("%d расширений GNOME", gnomeExt))
		}
	}
}

func detectMetaClasses() {
	home := os.Getenv("HOME")
	if home == "" { home = "/root" }
	if os.Getenv("WSL_DISTRO_NAME") != "" || os.Getenv("WSL_INTEROP") != "" {
		MetaWSL = true; Reasons["wsl"] = "окружение WSL"
	} else {
		if d, err := os.ReadFile("/proc/sys/kernel/osrelease"); err == nil {
			s := strings.ToLower(string(d))
			if strings.Contains(s, "microsoft") || strings.Contains(s, "wsl") { MetaWSL = true; Reasons["wsl"] = "окружение WSL" }
		}
		if !MetaWSL {
			if d, err := os.ReadFile("/proc/version"); err == nil {
				s := strings.ToLower(string(d))
				if strings.Contains(s, "microsoft") || strings.Contains(s, "wsl") { MetaWSL = true; Reasons["wsl"] = "окружение WSL" }
			}
		}
	}
	if os.Getenv("TERMUX_VERSION") != "" || strings.Contains(os.Getenv("PREFIX"), "com.termux") || dirExists("/data/data/com.termux") {
		MetaAndroid = true; Reasons["android"] = "Termux на Android"
	}
	if !MetaWSL && !MetaAndroid {
		if Has("systemd-detect-virt") {
			if out, err := exec.Command("systemd-detect-virt").Output(); err == nil {
				virt := strings.TrimSpace(string(out))
				if virt != "" && virt != "none" { MetaVM = true; Reasons["vm"] = virt }
			}
		}
		if _, err := os.Stat("/.dockerenv"); err == nil { MetaVM = true; if Reasons["vm"] == "" { Reasons["vm"] = "docker" } }
		if d, err := os.ReadFile("/proc/1/cgroup"); err == nil {
			s := string(d)
			if strings.Contains(s, "docker") || strings.Contains(s, "lxc") || strings.Contains(s, "kubepods") || strings.Contains(s, "containerd") || strings.Contains(s, "podman") {
				MetaVM = true; if Reasons["vm"] == "" { Reasons["vm"] = "container" }
			}
		}
		if container := os.Getenv("container"); container != "" { MetaVM = true; Reasons["vm"] = container }
		for _, f := range []string{"/sys/class/dmi/id/product_name", "/sys/class/dmi/id/sys_vendor"} {
			if d, err := os.ReadFile(f); err == nil {
				s := strings.ToLower(string(d))
				if strings.Contains(s, "qemu") || strings.Contains(s, "kvm") || strings.Contains(s, "virtualbox") || strings.Contains(s, "vmware") || strings.Contains(s, "innotek") || strings.Contains(s, "bochs") || strings.Contains(s, "xen") || strings.Contains(s, "hyper-v") || strings.Contains(s, "virtual machine") {
					MetaVM = true; if Reasons["vm"] == "" { Reasons["vm"] = "hypervisor" }
				}
			}
		}
	}
	if dirExists("/boot/efi/EFI/Microsoft") || dirExists("/boot/EFI/Microsoft") || dirExists("/efi/EFI/Microsoft") {
		MetaDualboot = true; Reasons["dualboot"] = "рядом обнаружен Windows"
	}
	if Has("lsblk") {
		if out, err := exec.Command("lsblk", "-no", "FSTYPE").Output(); err == nil {
			s := strings.ToLower(string(out))
			if strings.Contains(s, "ntfs") || strings.Contains(s, "bitlocker") { MetaDualboot = true; Reasons["dualboot"] = "рядом обнаружен Windows" }
		}
	}
	for _, f := range []string{"/boot/grub/grub.cfg", "/boot/grub2/grub.cfg"} {
		if d, err := os.ReadFile(f); err == nil {
			s := strings.ToLower(string(d))
			if strings.Contains(s, "windows") || strings.Contains(s, "microsoft") { MetaDualboot = true; Reasons["dualboot"] = "рядом обнаружен Windows" }
		}
	}
}

func analyzeBehavior() {
	if Behavior == "" { return }
	behav(`docker|docker-compose`, "devops", "docker в работе", 7, 5)
	behav(`kubectl|helm|k9s|kustomize`, "devops", "kubernetes в работе", 9, 3)
	behav(`terraform|tofu|ansible|pulumi`, "devops", "IaC в истории", 8, 3)
	behav(`systemctl|journalctl`, "sysadmin", "управление сервисами", 7, 6)
	behav(`ssh|scp|sftp`, "sysadmin", "удалённые хосты", 5, 8)
	behav(`nginx|certbot|iptables|nft|ufw`, "sysadmin", "серверная эксплуатация", 6, 4)
	behav(`psql|mysql|mariadb|redis-cli`, "sysadmin", "работа с БД", 5, 4)
	behav(`make|cargo|npm|pnpm|yarn|pip|pip3|gradle|mvn`, "programmer", "сборка/разработка", 7, 6)
	behav(`git`, "programmer", "git в повседневной работе", 6, 12)
	behav(`vim|nvim|emacs`, "programmer", "редактор кода в истории", 4, 10)
	behav(`nmap|nikto|sqlmap|msfconsole|hydra|hashcat|aircrack-ng|gobuster`, "pentester", "пентест в истории", 11, 2)
	behav(`tor|proxychains|openvpn|wg|gpg|veracrypt`, "anonymous", "приватность в истории", 6, 3)
	behav(`pacman|yay|paru|makepkg|emerge`, "old_hacker", "ручное управление пакетами", 4, 8)
	behav(`flatpak|distrobox|toolbox|rpm-ostree|nixos-rebuild|nix-shell|nix-env`, "atomic", "иммутабельный workflow", 7, 3)
	behav(`nix`, "fresh_witness", "Nix в истории", 4, 3)
	behav(`steam|lutris|wine|proton|protontricks`, "gamer", "игры в истории", 6, 2)
}

type ArchetypeResult struct {
	Key       string
	Label     string
	Score     int
	NormScore int
	Reason    string
}

func Normalize() []ArchetypeResult {
	maxScore := 1
	for key, s := range Score {
		if data.Hidden[key] { continue }
		if s > maxScore { maxScore = s }
	}
	normScore := map[string]int{}
	for key := range Score {
		if data.Mystery[key] || key == "normis" { continue }
		normScore[key] = Score[key] * 100 / maxScore
		// Скрытые классы нормализуются по максимуму видимых и могут его
		// превысить — без ограничения makeBar уйдёт в отрицательный Repeat.
		if normScore[key] > 100 { normScore[key] = 100 }
	}
	if maxScore <= 12 { normScore["normis"] = 100
	} else if maxScore >= 40 { normScore["normis"] = 0
	} else { normScore["normis"] = (40 - maxScore) * 100 / 28 }
	Reasons["normis"] = "мало ярких сигналов других классов"
	normScore["vm"] = boolToInt(MetaVM)
	normScore["wsl"] = boolToInt(MetaWSL)
	normScore["android"] = boolToInt(MetaAndroid)
	normScore["dualboot"] = boolToInt(MetaDualboot)
	normScore["corporat"] = boolToInt(MetaCorporat)
	var results []ArchetypeResult
	for key, pct := range normScore {
		if data.Hidden[key] && pct <= 20 { continue }
		results = append(results, ArchetypeResult{Key: key, Label: data.Labels[key], Score: Score[key], NormScore: pct, Reason: Reasons[key]})
	}
	sort.Slice(results, func(i, j int) bool { return results[i].NormScore > results[j].NormScore })
	var normal, mystery []ArchetypeResult
	for _, r := range results {
		if data.Mystery[r.Key] { mystery = append(mystery, r) } else { normal = append(normal, r) }
	}
	return append(normal, mystery...)
}

func boolToInt(b bool) int { if b { return 100 }; return 0 }

var MetaVM, MetaWSL, MetaAndroid, MetaDualboot, MetaCorporat bool

func fileExists(path string) bool { _, err := os.Stat(path); return err == nil }
func dirExists(path string) bool { info, err := os.Stat(path); return err == nil && info.IsDir() }

func hasPythonSitePackages(home string) bool {
	pattern := home + "/.local/lib/python*/site-packages"
	matches, _ := filepath.Glob(pattern)
	return len(matches) > 0
}

func detectVirtualizer() {
	home := os.Getenv("HOME")
	if home == "" { home = "/root" }
	score := 0
	// Docker
	if Has("docker") { score += 3 }
	// Compose v2 — не бинарник в PATH, а cli-плагин docker.
	if Has("docker-compose") ||
		fileExists("/usr/lib/docker/cli-plugins/docker-compose") ||
		fileExists("/usr/libexec/docker/cli-plugins/docker-compose") ||
		fileExists(home+"/.docker/cli-plugins/docker-compose") { score += 2 }
	// Podman
	if Has("podman") { score += 3 }
	if Has("podman-compose") { score += 2 }
	// LXC/LXD
	if Has("lxc") || Has("lxd") || Has("incus") { score += 4 }
	// KVM/QEMU
	if Has("qemu-system-x86_64") || Has("qemu-system-aarch64") { score += 4 }
	if Has("virsh") || Has("virt-manager") { score += 3 }
	// Proxmox
	if Has("pvesh") || Has("qm") { score += 5 }
	// VirtualBox
	if Has("VBoxManage") { score += 3 }
	// VMware
	if Has("vmware") { score += 3 }
	// Docker containers running
	if Has("docker") {
		out, _ := exec.Command("docker", "ps", "-q").Output()
		if len(strings.TrimSpace(string(out))) > 0 { score += 2 }
	}
	if score >= 8 {
		Add("virtualizer", score, fmt.Sprintf("виртуализация (score: %d)", score))
	}
}

func detectTiling() {
	home := os.Getenv("HOME")
	if home == "" { home = "/root" }
	score := 0
	desktop := os.Getenv("XDG_CURRENT_DESKTOP") + "|" + os.Getenv("DESKTOP_SESSION") + "|" + os.Getenv("XDG_SESSION_DESKTOP")
	// WM detection
	if strings.Contains(desktop, "i3") || Has("i3") { score += 4 }
	if strings.Contains(desktop, "sway") || Has("sway") { score += 4 }
	// Бинарник и XDG_CURRENT_DESKTOP у Hyprland — с большой буквы.
	if strings.Contains(desktop, "Hyprland") || strings.Contains(desktop, "hyprland") || Has("Hyprland") { score += 4 }
	if strings.Contains(desktop, "bspwm") || Has("bspwm") { score += 3 }
	if strings.Contains(desktop, "dwm") || Has("dwm") { score += 3 }
	if strings.Contains(desktop, "xmonad") || Has("xmonad") { score += 3 }
	if strings.Contains(desktop, "qtile") || Has("qtile") { score += 3 }
	if strings.Contains(desktop, "awesome") || Has("awesome") { score += 3 }
	if strings.Contains(desktop, "river") || Has("river") { score += 3 }
	if strings.Contains(desktop, "niri") || Has("niri") { score += 3 }
	// Tiling configs
	if dirExists(home + "/.config/i3") { score += 2 }
	if dirExists(home + "/.config/sway") { score += 2 }
	if dirExists(home + "/.config/hypr") { score += 2 }
	if dirExists(home + "/.config/bspwm") { score += 2 }
	if score >= 6 {
		Add("tiling", score, fmt.Sprintf("тайлинг (score: %d)", score))
	}
}

func detectNeovimAddict() {
	home := os.Getenv("HOME")
	if home == "" { home = "/root" }
	score := 0
	if Has("nvim") { score += 3 }
	if Has("vim") { score += 1 }
	// Neovim config
	if dirExists(home + "/.config/nvim") {
		score += 3
		// Count config files
		entries, _ := os.ReadDir(home + "/.config/nvim")
		fileCount := 0
		for _, e := range entries {
			if strings.HasSuffix(e.Name(), ".lua") || strings.HasSuffix(e.Name(), ".vim") {
				fileCount++
			}
		}
		if fileCount >= 5 { score += 2 }
		if fileCount >= 10 { score += 2 }
	}
	// Lazy.nvim or packer
	if dirExists(home + "/.local/share/nvim/lazy") { score += 2 }
	if dirExists(home + "/.local/share/nvim/site/pack") { score += 2 }
	if score >= 6 {
		Add("neovim_addict", score, fmt.Sprintf("Neovim (score: %d)", score))
	}
}

func detectShellCollector() {
	home := os.Getenv("HOME")
	if home == "" { home = "/root" }
	score := 0
	// Zsh
	if Has("zsh") { score += 2 }
	if dirExists(home + "/.oh-my-zsh") { score += 3 }
	if dirExists(home + "/.zsh") { score += 2 }
	// Starship
	if Has("starship") { score += 2 }
	if dirExists(home + "/.config/starship") { score += 1 }
	// P10k
	if fileExists(home + "/.p10k.zsh") { score += 2 }
	// Fish
	if Has("fish") { score += 2 }
	// Nushell
	if Has("nu") { score += 2 }
	// Elvish
	if Has("elvish") { score += 2 }
	// Many aliases/functions in config
	if fileExists(home + "/.zshrc") {
		d, _ := os.ReadFile(home + "/.zshrc")
		content := string(d)
		aliasCount := strings.Count(content, "alias ")
		functionCount := strings.Count(content, "function ") + strings.Count(content, "()")
		if aliasCount >= 5 { score += 2 }
		if functionCount >= 3 { score += 2 }
	}
	if score >= 6 {
		Add("shell_collector", score, fmt.Sprintf("shell (score: %d)", score))
	}
}

func detectSelfBuilder() {
	score := 0
	// Gentoo
	if Has("emerge") || Has("portage") { score += 5 }
	if fileExists("/etc/portage/make.conf") { score += 3 }
	// NixOS
	if Has("nixos-rebuild") || Has("nix") { score += 4 }
	// LFS
	if fileExists("/etc/lfs-release") { score += 5 }
	// Source builds
	if Has("make") && Has("gcc") { score += 2 }
	if Has("cmake") { score += 1 }
	// Custom kernels
	if Has("make") && fileExists("/usr/src/linux") { score += 3 }
	if score >= 6 {
		Add("self_builder", score, fmt.Sprintf("самосборка (score: %d)", score))
	}
}

func detectWaylandWafer() {
	score := 0
	session := os.Getenv("XDG_SESSION_TYPE")
	if session == "wayland" { score += 3 }
	// Wayland compositors
	if Has("sway") { score += 3 }
	if Has("Hyprland") { score += 3 }
	if Has("river") { score += 3 }
	if Has("niri") { score += 3 }
	// Wayland tools
	if Has("waybar") { score += 2 }
	if Has("wofi") { score += 2 }
	if Has("rofi") { score += 1 }
	if Has("wlogout") { score += 2 }
	if Has("swaylock") { score += 2 }
	if Has("mako") { score += 2 }
	if Has("dunst") { score += 1 }
	if Has("foot") { score += 2 }
	if Has("kitty") { score += 1 }
	if score >= 8 {
		Add("wayland_wafer", score, fmt.Sprintf("Wayland (score: %d)", score))
	}
}

func detectConsoleLife() {
	score := 0
	// RSS readers
	if Has("newsboat") { score += 3 }
	if Has("newsblur") { score += 3 }
	if Has("snownews") { score += 3 }
	// Email
	if Has("mutt") || Has("neomutt") { score += 3 }
	if Has("offlineimap") { score += 2 }
	if Has("isync") { score += 2 }
	// Music
	if Has("ncmpcpp") { score += 2 }
	if Has("mpd") { score += 2 }
	if Has("cmus") { score += 2 }
	// IRC/Chat
	if Has("irssi") { score += 2 }
	if Has("weechat") { score += 2 }
	// File manager
	if Has("ranger") { score += 2 }
	if Has("lf") { score += 2 }
	if Has("nnn") { score += 2 }
	if Has("mc") { score += 1 }
	// Other console tools
	if Has("htop") || Has("btop") { score += 1 }
	if Has("tmux") || Has("screen") { score += 1 }
	if Has("fzf") { score += 1 }
	if Has("rg") { score += 1 }
	if Has("fd") { score += 1 }
	if Has("bat") { score += 1 }
	if Has("exa") || Has("eza") { score += 1 }
	if score >= 8 {
		Add("console_life", score, fmt.Sprintf("консоль (score: %d)", score))
	}
}

func detectPackageFreak() {
	home := os.Getenv("HOME")
	if home == "" { home = "/root" }
	score := 0
	// Flatpak
	if Has("flatpak") {
		score += 2
		out, _ := exec.Command("flatpak", "list", "--app", "--columns=application").Output()
		fpCount := strings.Count(string(out), "\n")
		if fpCount >= 10 { score += 2 }
		if fpCount >= 20 { score += 2 }
	}
	// Snap
	if Has("snap") {
		score += 2
		out, _ := exec.Command("snap", "list").Output()
		snapCount := strings.Count(string(out), "\n") - 1
		if snapCount >= 10 { score += 2 }
		if snapCount >= 20 { score += 2 }
	}
	// AppImage
	if dirExists(home + "/Applications") {
		entries, _ := os.ReadDir(home + "/Applications")
		appimageCount := 0
		for _, e := range entries {
			if strings.HasSuffix(e.Name(), ".AppImage") { appimageCount++ }
		}
		if appimageCount >= 5 { score += 3 }
		if appimageCount >= 10 { score += 3 }
	}
	// Multiple package managers
	pkgManagers := 0
	if Has("apt") { pkgManagers++ }
	if Has("dnf") { pkgManagers++ }
	if Has("pacman") { pkgManagers++ }
	if Has("zypper") { pkgManagers++ }
	if Has("xbps-install") { pkgManagers++ }
	if Has("apk") { pkgManagers++ }
	if pkgManagers >= 3 { score += 2 }
	if score >= 8 {
		Add("package_freak", score, fmt.Sprintf("пакеты (score: %d)", score))
	}
}

func detectMusician() {
	score := 0
	// DAW
	if Has("lmms") { score += 4 }
	if Has("ardour") { score += 4 }
	if Has("bitwig-studio") { score += 4 }
	if Has("reaper") { score += 4 }
	if Has("hydrogen") { score += 3 }
	if Has("mixxx") { score += 3 }
	// MIDI
	if Has("timidity") { score += 2 }
	if Has("fluidsynth") { score += 2 }
	// Audio
	if Has("mpd") { score += 2 }
	if Has("ncmpcpp") { score += 2 }
	if Has("cmus") { score += 2 }
	if Has("audacity") { score += 2 }
	// Synths
	if Has("synthesia") { score += 2 }
	if Has("musescore") { score += 3 }
	// JACK
	if Has("jackd") { score += 2 }
	if Has("pipewire") { score += 1 }
	if score >= 6 {
		Add("musician", score, fmt.Sprintf("музыка (score: %d)", score))
	}
}

func detectPhotographer() {
	score := 0
	if Has("darktable") { score += 4 }
	if Has("rawtherapee") { score += 4 }
	if Has("digikam") { score += 3 }
	if Has("gimp") { score += 2 }
	if Has("krita") { score += 2 }
	if Has("shotwell") { score += 2 }
	// RAW processing
	if Has("dcraw") { score += 2 }
	if Has("exiftool") { score += 1 }
	if Has("exiv2") { score += 1 }
	if score >= 5 {
		Add("photographer", score, fmt.Sprintf("фото (score: %d)", score))
	}
}

func detectVideoEditor() {
	score := 0
	if Has("kdenlive") { score += 4 }
	if Has("shotcut") { score += 4 }
	if Has("openshot") { score += 3 }
	if Has("pitivi") { score += 3 }
	if Has("flowblade") { score += 3 }
	if Has("olive") { score += 3 }
	if Has("ffmpeg") { score += 2 }
	if Has("obs") { score += 3 }
	if Has("mkvtoolnix") { score += 1 }
	if Has("handbrake") { score += 2 }
	if score >= 5 {
		Add("video_editor", score, fmt.Sprintf("видео (score: %d)", score))
	}
}

func detectModeler3d() {
	score := 0
	if Has("blender") { score += 5 }
	if Has("freecad") { score += 4 }
	if Has("openscad") { score += 3 }
	if Has("solvespace") { score += 3 }
	if Has("bambu-studio") { score += 3 }
	if Has("prusa-slicer") { score += 3 }
	if Has("cura") { score += 3 }
	if Has("slic3r") { score += 3 }
	if Has("orca-slicer") { score += 3 }
	if Has("meshlab") { score += 2 }
	if Has("cloudcompare") { score += 2 }
	if score >= 5 {
		Add("modeler3d", score, fmt.Sprintf("3D (score: %d)", score))
	}
}

func detectWriter() {
	score := 0
	if Has("pdflatex") || Has("xelatex") || Has("lualatex") { score += 4 }
	if Has("pandoc") { score += 3 }
	if Has("lyx") { score += 3 }
	if Has("zathura") { score += 2 }
	if Has("evince") { score += 1 }
	if Has("okular") { score += 1 }
	if Has("calibre") { score += 2 }
	if Has("hledger") { score += 2 }
	if Has("ledger") { score += 2 }
	if Has("typst") { score += 3 }
	if Has("quarto") { score += 2 }
	if score >= 5 {
		Add("writer", score, fmt.Sprintf("тексты (score: %d)", score))
	}
}

func detectStreamer() {
	score := 0
	if Has("obs") { score += 4 }
	if Has("streamlabs") { score += 4 }
	if Has("ffmpeg") { score += 2 }
	if Has("parecord") || Has("pactl") { score += 1 }
	if Has("pulseaudio") { score += 1 }
	if Has("pipewire") { score += 1 }
	if Has("wireplumber") { score += 1 }
	// v4l2loopback — модуль ядра, ищем его среди загруженных, а не в PATH.
	if fileExists("/sys/module/v4l2loopback") { score += 2 }
	if Has("xdotool") { score += 1 }
	if Has("xprop") { score += 1 }
	if score >= 5 {
		Add("streamer", score, fmt.Sprintf("стримы (score: %d)", score))
	}
}

func detectEmbedded() {
	score := 0
	home := os.Getenv("HOME")
	if home == "" { home = "/root" }
	// STM32
	if Has("st-flash") || Has("st-info") || Has("st-util") { score += 4 }
	if Has("stm32flash") { score += 3 }
	// ESP
	if Has("esptool") || Has("esptool.py") { score += 4 }
	if Has("idf.py") { score += 4 }
	if Has("platformio") || Has("pio") { score += 4 }
	// Arduino
	if Has("arduino-cli") { score += 3 }
	if Has("arduino") { score += 3 }
	// ARM
	if Has("arm-none-eabi-gcc") { score += 4 }
	if Has("arm-none-eabi-objcopy") { score += 2 }
	// RISC-V
	if Has("riscv64-unknown-elf-gcc") { score += 4 }
	// Debug
	if Has("openocd") { score += 3 }
	if Has("JLinkExe") || Has("JLinkGDBServer") { score += 3 }
	if Has("stlink-server") { score += 3 }
	// Other embedded tools
	if Has("dfu-util") { score += 2 }
	if Has("picocom") { score += 2 }
	if Has("minicom") { score += 2 }
	if Has("screen") { score += 1 }
	// VS Code extensions for embedded
	if dirExists(home + "/.vscode/extensions") {
		entries, _ := os.ReadDir(home + "/.vscode/extensions")
		for _, e := range entries {
			name := strings.ToLower(e.Name())
			if strings.Contains(name, "espressif") || strings.Contains(name, "esp-idf") ||
				strings.Contains(name, "stm32") || strings.Contains(name, "platformio") ||
				strings.Contains(name, "arduino") || strings.Contains(name, "cmsis") ||
				strings.Contains(name, "cortex") || strings.Contains(name, "riscv") {
				score += 4
				break
			}
		}
	}
	if score >= 6 {
		Add("embedded", score, fmt.Sprintf("embedded (score: %d)", score))
	}
}
