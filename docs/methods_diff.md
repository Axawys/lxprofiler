# Разница в методиках подсчёта: Bash vs Go

## Текущий статус

Go-версия содержит **все** проверки Bash **плюс** дополнительные Go-добавления.

## Сравнение статического вывода

### Одинаковые описания (порядок идентичен):
- Компиляторы, языки, редакторы. Машина для тебя — мастерская, а не витрина.
- Cattle, не pets. Кластеры, пайплайны и декларативная инфраструктура важнее отдельной машины.
- Ты помнишь, как всё должно работать, и был здесь до systemd. Философия Unix — твоя привычка.
- Сервисы и демоны тебе привычны, ты держишь систему в порядке.
- Местами проскакивает российский софт, но не как принцип.
- Иногда запускаешь игру-другую, но геймингом это не назвать.
- Пара Flatpak'ов есть, но до атомарной философии далеко.
- Дефолтная тема — и так сойдёт. Скриншот рабочего стола ты не запостишь.
- nmap и Metasploit — не твои инструменты. Чужие порты тебя не зовут.
- Хайповые новинки тебя не трогают — пусть сначала отлежатся пару лет.
- Творческого софта на машине не видно — она не про искусство.
- Минимализм — не про тебя: пусть стоит всё, что может пригодиться.

### Различия в детекции:

| Признак | Bash | Go | Причина |
|---------|------|-----|---------|
| `git в повседневной работе` | ×59 | ×85 | Go читает больше файлов (.bashrc, .zshrc, .bash_aliases, .profile) |
| `сборка/разработка` | нет | ×6 | Go-добавление: `behav('make\|cargo\|npm...', programmer, 7)` |
| `LVM` | нет | да | Go-добавление: `Has("lvs")` |
| `игры в истории` | нет | ×5 | Go-добавление: `behav('steam\|lutris\|wine...', gamer, 6)` |
| `приватность в истории` | нет | ×34 | Go-добавление: `behav('tor\|proxychains\|openvpn...', anonymous, 6)` |
| `иммутабельный workflow` | ×11 | ×13 | Разные источники BEHAVIOR |

### Go-добавления (нет в Bash):

1. **сборка/разработка (×6)** — make/cargo/npm/pip/gradle/mvn
2. **LVM** — проверка `Has("lvs")`
3. **игры в истории (×5)** — steam/lutris/wine/proton
4. **приватность в истории (×34)** — tor/proxychains/gpg/veracrypt
5. **44 игр в меню** — подсчёт .desktop файлов
6. **1С:Предприятие** — определение по PATH
7. **запущен lxprofiler** — self-detection
8. **Вайбкодер** — AI-инструменты (cursor, copilot, aider, codex, claude, continue, cody, phind, tabnine, supermaven, codeium + плагины в VS Code/JetBrains/Neovim)
9. **Виртуализатор** — Docker, Podman, LXC, KVM, QEMU, Proxmox, VirtualBox, VMware
10. **Тайлинг** — i3, sway, hyprland, bspwm, dwm, xmonad, qtile, awesome, river, niri + конфиги
11. **Neovim-аддикт** — nvim + конфиг ~/.config/nvim + lazy.nvim/packer
12. **Shell-коллекционер** — Zsh + oh-my-zsh + starship + p10k + aliases/functions
13. **Самосборщик** — emerge/portage, nixos-rebuild, make+gcc, кастомные ядра
14. **Wayland-вафлер** — Wayland + compositor + waybar + wofi + rofi + wlogout + swaylock
15. **Консольный жизни** — newsboat, mutt, ncmpcpp, mpd, irssi, ranger, lf, tmux, fzf, ripgrep (rg)
16. **Пакетоман** — flatpak/snap/appimage с количеством + несколько пакетных менеджеров
17. **Музыкант** — LMMS, Ardour, Bitwig, Reaper, Hydrogen, Mixxx, MPD, ncmpcpp, cmus, Audacity, MuseScore, JACK
18. **Фотограф** — darktable, RawTherapee, digiKam, GIMP, Krita, Shotwell, dcraw, exiftool
19. **Видеомонтажёр** — kdenlive, Shotcut, OpenShot, Pitivi, Flowblade, Olive, FFmpeg, OBS, HandBrake
20. **3D-моделлер** — Blender, FreeCAD, OpenSCAD, Bambu Studio, PrusaSlicer, Cura, Slic3r, OrcaSlicer, MeshLab
21. **Писатель** — pdflatex/xelatex/lualatex, Pandoc, LyX, Zathura, Calibre, Typst, Quarto
22. **Стример** — OBS, Streamlabs, FFmpeg, PipeWire, v4l2loopback
23. **Ветеран** — 32-bit, X11, ancient WM (Trinity/FVWM/IceWM), MPlayer, Pidgin, старое ядро
24. **Тинкерер** — make+gcc, CMake, Meson, custom kernels, Gentoo/NixOS, AUR, /opt
25. **Embedded** — STM32 (st-flash), ESP32 (esptool/idf.py), PlatformIO, Arduino CLI, ARM/RISC-V toolchains, JLinkExe/JLinkGDBServer, OpenOCD

### Причины расхождения в git (×59 vs ×85):

Bash читает только:
- `~/.bash_history`
- `~/.zsh_history`

Go читает:
- `~/.bash_history`
- `~/.zsh_history`
- `~/.local/share/fish/fish_history`
- `~/.bashrc`
- `~/.zshrc`
- `~/.bash_aliases`
- `~/.config/fish/config.fish`
- `~/.profile`

## Интерактивный режим (проценты):

| Архетип | Bash | Go | Разница |
|---------|------|-----|---------|
| Программист | 100% | 100% | 0% |
| DevOps | 95% | 87% | -8% |
| Последователь Столлмана | 83% | 77% | -6% |
| Сис-админ | 58% | 59% | +1% |
| Импортозаместитель | 45% | 40% | -5% |
| Геймер | 32% | 36% | +4% |
| Атомарник | 25% | 23% | -2% |
| Райсер | 12% | 11% | -1% |
| Хацкер | 9% | 9% | 0% |
| Анонимус | 0% | 6% | +6% |
| Свидетель свежего ПО | 7% | 6% | -1% |
| Творческая снежинка | 4% | 4% | 0% |
| Минималист | 3% | 3% | 0% |

### Причины расхождений в процентах:

1. **DevOps (-8%)**: Go не детектирует `Terraform` через `has_used`, Bash — через `has_used terraform`
2. **Последователь Столлмана (-6%)**: Go не детектирует `OpenRC` через `has openrc`, Bash — через `has openrc || [[ -f /etc/init.d/openrc ]]`
3. **Импортозаместитель (-5%)**: Go не детектирует `1С:Предприятие` (нет в PATH), Bash — через `compgen -G "/opt/1[Cc]*"`
4. **Анонимус (+6%)**: Go добавил `behav('tor|proxychains...', anonymous, 6, "приватность в истории", 3)` — этих проверок не было в оригинале
5. **Геймер (+4%)**: Go добавил `behav('steam|lutris|wine...', gamer, 6, "игры в истории", 2)` — этих проверок не было в оригинале

## Вывод:

Go-версия содержит:
- Все 16 поведенческих проверок из Bash
- Все проверки наличия инструментов из Bash
- 7 дополнительных Go-добавлений
- Больше источников для поведенческого анализа (8 файлов вместо 2)
