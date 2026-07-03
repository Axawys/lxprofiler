# lxprofiler — Расчёт категорий

## Скрытые категории

Скрытые категории показываются только при заполненности > 20%.

### Нормис
**Порог:** > 20% (обратная зависимость от макс. баллов других классов)
**Формула:** `(40 - MAX_SCORE) * 100 / 28` (если MAX_SCORE < 40)
**Критерий:** Мало ярких сигналов других классов

### ЭТО ВИРТУАЛКА!
**Порог:** 100% при обнаружении
**Критерий:** Запуск в виртуальной машине или контейнере
**Детекция:** systemd-detect-virt, /.dockerenv, /proc/1/cgroup, DMI

### WSL
**Порог:** 100% при обнаружении
**Критерий:** Запуск внутри WSL
**Детекция:** WSL_DISTRO_NAME, WSL_INTEROP, /proc/sys/kernel/osrelease

### Дуалбут
**Порог:** 100% при обнаружении
**Критерий:** Рядом обнаружен Windows
**Детекция:** /boot/efi/EFI/Microsoft, NTFS, grub.cfg

### АндроЕд
**Порог:** 100% при обнаружении
**Критерий:** Запуск в Termux на Android
**Детекция:** TERMUX_VERSION, PREFIX, /data/data/com.termux

### Корпорат
**Порог:** 100% при обнаружении
**Критерий:** Дистрибутив RHEL
**Детекция:** Red Hat Enterprise Linux в /etc/os-release

---

## Основные категории

### DevOps
**Кластеры, пайплайны и декларативная инфраструктура**
- Docker (+5, поведение ×5: +7)
- Podman (+5, поведение)
- Kubernetes (+14, поведение ×3: +9)
- k9s (+8, поведение)
- Helm (+8, поведение)
- Локальный кластер (minikube/kind) (+6)
- GitOps (argocd/flux/skaffold) (+8)
- Terraform (+12, поведение)
- OpenTofu (+10, поведение)
- Pulumi (+8, поведение)
- Ansible (+8, поведение ×3: +8)
- Config management (puppet/chef/salt) (+8)
- Vagrant (+6, поведение)
- QEMU/KVM (+8)
- LXC/Incus (+6)
- HashiCorp-стек (vault/consul/nomad) (+8)
- Облачный CLI (aws/gcloud/az/doctl) (+8, поведение)
- CI-раннеры (gitlab-runner/act) (+5)
- Docker Compose (+4)
- MongoDB (+5)
- Метрики (Prometheus) (+3)
- Дашборды (Grafana) (+3)

### Программист
**Компиляторы, языки, редакторы**
- Языки (node/deno/bun/rustc/go/zig/elixir/julia/haskell/scala/kotlin/crystal/nim/ocaml) (+3 за каждый)
- Ruby (+1)
- Python + инструменты разработки (+3)
- Java JDK/сборка (+3)
- Cargo (+5)
- rustup (+2)
- Go workspace (+3)
- npm-проекты (+3)
- Менеджер версий (asdf/nvm/pyenv/rbenv/sdk) (+5)
- Настроенный git (user.*) (+6)
- Полиглот (5+ языков) (+10)
- Несколько языков (3+) (+4)
- gcc/clang (+3)
- Сборочные системы (make+cmake) (+3)
- Rust (+6)
- Vim/Neovim (+6, поведение ×10: +4)
- Emacs (+12)
- Helix (+6)
- VS Code (+6)
- Zed (+6)
- VSCodium (+6)
- Cursor (+6)
- Sublime Text (+5)
- Lapce (+5)
- micro (+3)
- Кастомный Neovim (+5)
- JetBrains IDE (+6)
- Мультиплексор (tmux/screen) (+6)
- Отладчики (gdb/lldb) (+4)
- Сборка/разработка (поведение ×6: +7)
- git в повседневной работе (поведение ×12: +6)

### Сис-админ
**Живые серверы и демоны**
- nginx (+8, поведение ×4: +6)
- Apache (+6)
- PostgreSQL (+8)
- MySQL/MariaDB (+6)
- Redis (+5)
- SSH-сервер (+5)
- Мониторинг (htop/btop/glances) (+3)
- Prometheus (+6)
- Grafana (+6)
- Бэкапы (borg/restic/rsnapshot/duplicity) (+8)
- S.M.A.R.T. (+3)
- fail2ban (+6)
- APT-обновления (поведение)
- AUR (поведение ×8: +4)

### Минималист
**Мало процессов, лёгкие WM**
- GNOME (+8)
- XFCE (+8)
- LXQt/LXDE (+10)
- Cinnamon (+6)
- sway (+12)
- i3 (+12)
- bspwm (+10)
- foot (+5)
- ≤4 GB RAM (+10)
- 8 GB RAM (+5)
- Chistaya FS (<30%) (+5)
- Мало пакетов (≤300: +10, ≤600: +5)
- dash (+8)
- systemd-boot (+6)
- EFISTUB (+8)
- syslinux (+3)
- Snap (+2)
- Без кастомизации окружения (+5)

### Последователь Столлмана
**Был здесь до systemd**
- Artix (+12)
- EndeavourOS (+8)
- Arch (+8)
- MX Linux (+12)
- Gentoo (+14)
- Slackware (+16)
- Void Linux (+12)
- Debian (+10)
- MATE (+10)
- dwm (+14)
- Emacs (+12)
- X11 (+5)
- Старое ядро (≤4.x: +12, 5.x: +4)
- 32-бит (+12)
- LILO (+12)
- runit/s6 (+12)
- dinit/OpenRC (+10)
- gcc/clang (+3)
- Мультиплексор (+6)
- AUR-хелпер (+4)
- ручное управление пакетами (поведение ×8: +4)

### Райсер
**Вылизанный рабочий стол**
- KDE Plasma (+10)
- Hyprland (+12)
- COSMIC (+12)
- elementary OS (+12)
- Zorin OS (+8)
- Deepin (+10)
- Pantheon (+10)
- Budgie (+8)
- niri (+5)
- rEFInd (+8)
- Кастомная тема GRUB (+5)
- zsh (+8)
- kitty (+6)
- Alacritty (+5)
- WezTerm (+6)
- 6+ кастомных конфигов (+12)
- 3+ кастомных конфигов (+6)
- Dotfiles в Git (+8)
- Менеджер dotfiles (stow/chezmoi) (+5)
- fetch-алиас (+4)
- Кастомный Neovim (+5)
- Расширения GNOME (+1-6)
- Креатив/медиа flatpak (+4)

### Геймер
**Steam, Proton, Wine**
- Steam (+12)
- Lutris (+8)
- Heroic (+6)
- Bottles (+6)
- Wine (+6)
- GameMode (+6)
- MangoHud (+5)
- Proton (+5)
- RetroArch (+6)
- NVIDIA GPU (+6)
- Каталог Games/Игры (5+: +12, 1+: +8)
- 15+ игр в меню (+14)
- 5+ игр в меню (+8)
- 1+ игр в меню (+3)
- Ядро Zen/XanMod/Liquorix/TkG (+5)
- SteamOS (+12)
- Nobara (+12)
- Pop!_OS (+10)
- Garuda (+10)
- Bazzite (+12)
- Флагманское железо (64+ GB RAM) (+6)
- Мощное железо (32+ GB RAM) (+5)
- 16+ ядер (+4)
- Игры в истории (поведение ×2: +6)

### Анонимус
**Tor, VPN, шифрование**
- Tor (+12)
- Tor Browser (+10)
- I2P (+10)
- VeraCrypt (+14)
- GPG (+6, поведение)
- pass (+8, поведение)
- age (+6, поведение)
- cryptsetup (+5, поведение)
- Tomb (+8, поведение)
- KeePassXC (+5)
- gocryptfs (+6, поведение)
- EncFS (+5, поведение)
- proxychains (+8, поведение ×4: +4)
- Mullvad VPN (+10)
- ProtonVPN (+8)
- OpenVPN (+6, поведение)
- WireGuard (+6, поведение)
- Mullvad Browser (+8)
- LibreWolf (+6)
- Brave (+4)
- Шифрование диска (LUKS) (+10)
- Tails (+20)
- Qubes OS (+16)
- Whonix (+18)
- Приватность в истории (поведение ×3: +6)

### Хацкер
**nmap, Metasploit, Wireshark**
- nmap (+12, поведение)
- masscan (+8, поведение)
- Wireshark (+8)
- tshark (+6, поведение)
- tcpdump (+5, поведение)
- Metasploit (+16, поведение)
- aircrack-ng (+12, поведение)
- hashcat (+10, поведение)
- John the Ripper (+10, поведение)
- hydra (+10, поведение)
- sqlmap (+10, поведение)
- nikto (+8, поведение)
- gobuster (+6, поведение)
- ffuf (+6, поведение)
- Burp Suite (+12)
- OWASP ZAP (+8)
- radare2 (+10, поведение)
- Ghidra (+12)
- binwalk (+6, поведение)
- Volatility (+8, поведение)
- WPScan (+6, поведение)
- Пентест в истории (поведение ×2: +11)
- Kali Linux (+18)
- Parrot OS (+16)
- BlackArch (+18)
- Pentoo (+18)

### Импортозаместитель
**Отечественный софт**
- Astra Linux (+28)
- RED OS (+28)
- ALT Linux (+24)
- ALT Atomic (+22)
- ROSA Linux (+24)
- Calculate Linux (+24)
- Simply Linux (+20)
- PostgreSQL (+4)
- nginx (+4)
- 1С:Предприятие (+8)
- ALT Tuner (+5)
- Яндекс.Браузер (+0, только описание)
- Русская локаль (+2)
- Запущен lxprofiler (+1)

### Свидетель свежего ПО
**Wayland, Rust, rolling-release**
- Fedora (+12)
- Manjaro (+5)
- EndeavourOS (+6)
- NixOS (+8)
- COSMIC (+12)
- Wayland (+10)
- fish (+6)
- nushell (+10)
- xonsh (+5)
- elvish (+6)
- Helix (+6)
- Rust (+6)
- Rust toolchain (+2)
- Редактор на Rust (Zed/Lapce) (+3)
- GPU-терминал (kitty/Alacritty/WezTerm) (+3)
- Wayland-терминал (foot) (+3)
- Декларативные пакеты (nix) (+4)
-.Immutable (Silverblue/Vanilla OS/blendOS) (+6-8)
- Limine (+8)
- systemd-boot (+3)
- Nix в истории (поведение ×3: +4)
- Свежая установка (≤14 дн.) (+4)

### Атомарник
**Иммутабельность и откаты**
- NixOS (+14)
- Silverblue/Kinoite/Sericea (+14)
- MicroOS/Aeon/Kalpa (+14)
- Vanilla OS (+14)
- ALT Atomic (+14)
- blendOS (+12)
- GNOME OS (+12)
- Bazzite (+12)
- SteamOS (+8)
- Endless OS (+8)
- Flatpak (+4)
- Nix (+8)
- distrobox/toolbox (+8)
- ostree-система (rpm-ostree/ostree/bootc) (+10)
- GNU Guix (+8)
- Nix/home-manager (+6)
- Откаты ФС (snapper) (+4)
- Flatpak-приложения (15+: +8, 5+: +5)
- Иммутабельный workflow (поведение ×3: +7)

### Творческая снежинка
**Фото/видео-редакторы и биты**
- GIMP (+4)
- Krita (+4)
- Darktable (+4)
- RawTherapee (+4)
- Inkscape (+4)
- Blender (+4)
- Kdenlive (+4)
- Shotcut (+4)
- OpenShot (+4)
- Flowblade (+4)
- Ardour (+4)
- LMMS (+4)
- Hydrogen (+4)
- Mixxx (+4)
- Audacity (+4)
- Bitwig Studio (+4)
- Renoise (+4)
- DaVinci Resolve (+6)
- Креатив/медиа flatpak (+5)

---

## Go-добавления (нет в оригинале)

### Вайбкодер
**AI-ассистенты для кода**
- aider (+5)
- codex (+5)
- claude/claude-code (+5)
- continue (+4)
- copilot (+4)
- cody (+4)
- phind (+3)
- tabnine (+3)
- supermaven (+3)
- codeium (+4)
- AI плагин в VS Code (+3)
- AI плагин в JetBrains (+3)
- AI плагин in Neovim (+3)
- **Порог:** 5+ (показывается всегда)

### Виртуализатор
**Виртуализация как стиль жизни**
- Docker (+3)
- docker-compose (+2)
- Podman (+3)
- podman-compose (+2)
- LXC/LXD/Incus (+4)
- QEMU (+4)
- virsh/virt-manager (+3)
- Proxmox (+5)
- VirtualBox (+3)
- VMware (+3)
- Запущенные контейнеры Docker (+2)
- **Порог:** 8+

### Тайлинг
**Тайлинг-менеджеры окон**
- i3 (+4)
- sway (+4)
- hyprland (+4)
- bspwm (+3)
- dwm (+3)
- xmonad (+3)
- qtile (+3)
- awesome (+3)
- river (+3)
- niri (+3)
- Конфиги tiling (+2)
- **Порог:** 6+

### Neovim-аддикт
**Конфиг Neovim длиннее кода**
- nvim (+3)
- vim (+1)
- Конфиг ~/.config/nvim (+3)
- 5+ конфиг-файлов (+2)
- 10+ конфиг-файлов (+2)
- lazy.nvim/packer (+2)
- **Порог:** 6+

### Shell-коллекционер
**50 плагинов в zsh**
- zsh (+2)
- oh-my-zsh (+3)
- .zsh (+2)
- starship (+2)
- p10k (+2)
- fish (+2)
- nushell (+2)
- elvish (+2)
- 5+ aliases (+2)
- 3+ functions (+2)
- **Порог:** 6+

### Самосборщик
**Собирает систему из исходников**
- emerge/portage (+5)
- make.conf (+3)
- nixos-rebuild/nix (+4)
- LFS (+5)
- make+gcc (+2)
- cmake (+1)
- Кастомные ядра (+3)
- **Порог:** 6+

### Wayland-вафлер
**Полный Wayland-стек**
- Wayland сессия (+3)
- sway/hyprland/river/niri (+3)
- waybar (+2)
- wofi (+2)
- rofi (+1)
- wlogout (+2)
- swaylock (+2)
- wlroots (+2)
- mako (+2)
- foot (+2)
- **Порог:** 8+

### Консольный жизни
**Всё через терминал**
- newsboat/newsblur/snownews (+3)
- mutt/neomutt (+3)
- offlineimap/isync (+2)
- ncmpcpp/mpd/cmus (+2)
- irssi/weechat (+2)
- ranger/lf/nnn (+2)
- htop/btop (+1)
- tmux/screen (+1)
- fzf/ripgrep (rg)/fd/bat/exa (+1)
- **Порог:** 8+

### Пакетоман
**Ставит всё что найдёт**
- flatpak (+2)
- flatpak 10+ apps (+2)
- flatpak 20+ apps (+2)
- snap (+2)
- snap 10+ packages (+2)
- snap 20+ packages (+2)
- AppImage 5+ (+3)
- AppImage 10+ (+3)
- 3+ пакетных менеджера (+2)
- **Порог:** 8+

### Музыкант
**DAW, сэмплеры, MIDI**
- LMMS (+4)
- Ardour (+4)
- Bitwig Studio (+4)
- Reaper (+4)
- Hydrogen (+3)
- Mixxx (+3)
- Timidity (+2)
- FluidSynth (+2)
- MPD (+2)
- ncmpcpp (+2)
- cmus (+2)
- Audacity (+2)
- MuseScore (+3)
- JACK (+2)
- **Порог:** 6+

### Фотограф
**darktable, rawtherapee**
- darktable (+4)
- RawTherapee (+4)
- digiKam (+3)
- GIMP (+2)
- Krita (+2)
- Shotwell (+2)
- dcraw (+2)
- exiftool (+1)
- **Порог:** 5+

### Видеомонтажёр
**kdenlive, shotcut**
- kdenlive (+4)
- Shotcut (+4)
- OpenShot (+3)
- Pitivi (+3)
- Flowblade (+3)
- Olive (+3)
- FFmpeg (+2)
- OBS (+3)
- HandBrake (+2)
- **Порог:** 5+

### 3D-моделлер
**Blender, FreeCAD**
- Blender (+5)
- FreeCAD (+4)
- OpenSCAD (+3)
- SolveSpace (+3)
- Bambu Studio (+3)
- PrusaSlicer (+3)
- Cura (+3)
- Slic3r (+3)
- OrcaSlicer (+3)
- MeshLab (+2)
- **Порог:** 5+

### Писатель
**LaTeX, pandoc**
- pdflatex/xelatex/lualatex (+4)
- Pandoc (+3)
- LyX (+3)
- Zathura (+2)
- Calibre (+2)
- Typst (+3)
- Quarto (+2)
- hledger (+2)
- **Порог:** 5+

### Стример
**OBS, streamlabs**
- OBS (+4)
- Streamlabs (+4)
- FFmpeg (+2)
- PipeWire (+1)
- v4l2loopback (+2)
- **Порог:** 5+

### Ветеран
**Старое железо и софт**
- 32-bit архитектура (+5)
- Старое ядро (≤4.x: +5, 5.x: +2)
- X11 сессия (+2)
- Trinity/KDE 3 (+5)
- FVWM/IceWM/twm (+4)
- MPlayer без mpv (+2)
- Pidgin (+2)
- **Порог:** 6+

### Тинкерер
**Собирает из исходников**
- make+gcc (+2)
- CMake/Meson (+1)
- Кастомные ядра (+3)
- Gentoo/emerge (+4)
- NixOS/nixos-rebuild (+4)
- AUR (yay/paru) (+1)
- /opt с 5+ проектами (+2)
- **Порог:** 6+

### Embedded-разработчик
**STM32, ESP32, Arduino**
- st-flash/st-info/st-util (+4)
- stm32flash (+3)
- esptool/idf.py (+4)
- PlatformIO/pio (+4)
- arduino-cli (+3)
- arm-none-eabi-gcc (+4)
- riscv64-unknown-elf-gcc (+4)
- OpenOCD (+3)
- JLinkExe/JLinkGDBServer (+3)
- dfu-util (+2)
- picocom/minicom (+2)
- **Порог:** 6+
