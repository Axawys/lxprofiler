#!/usr/bin/env bash
# lib/detect.sh — анализ системы, нормализация и сортировка
# (требует lib/data.sh и lib/helpers.sh)

# ──────────────────────────────────────────────
# Дистрибутив
# ──────────────────────────────────────────────

DISTRO="Unknown"
DISTRO_ID=""
DISTRO_LIKE=""
if [[ -f /etc/os-release ]]; then
  DISTRO=$(grep '^PRETTY_NAME=' /etc/os-release | cut -d= -f2- | tr -d '"')
  DISTRO_ID=$(grep '^ID=' /etc/os-release | cut -d= -f2- | tr -d '"')
  DISTRO_LIKE=$(grep '^ID_LIKE=' /etc/os-release | cut -d= -f2- | tr -d '"')
fi
DISTRO_ALL="${DISTRO} ${DISTRO_ID} ${DISTRO_LIKE}"

case "$DISTRO_ALL" in
  # ── Отечественные ──
  *Astra*|*astra*)             add import_substituted 18 "Astra Linux"; add sysadmin 6 "корпоративная ОС"; add anonymous 4 "мандатный доступ" ;;
  *RED\ OS*|*RedOS*|*redos*)   add import_substituted 18 "RED OS"; add sysadmin 6 "серверная ОС" ;;
  *ALT*Atomic*|*"ALT Atomic"*) add atomic 14 "ALT Atomic"; add import_substituted 14 "ALT (отеч.)"; add fresh_witness 4 "immutable" ;;
  *ALT\ *|*altlinux*|*"ALT "*) add import_substituted 16 "ALT Linux"; add old_hacker 4 "Sisyphus" ;;
  *ROSA*|*rosa*)               add import_substituted 16 "ROSA Linux" ;;
  *Calculate*|*calculate*)     add import_substituted 14 "Calculate Linux"; add old_hacker 6 "Gentoo-основа" ;;
  *Simply*)                    add import_substituted 14 "Simply Linux" ;;
  # ── Arch-семейство ──
  *Garuda*)      add gamer 10 "Garuda"; add ricer 8 "ricing из коробки" ;;
  *Artix*)       add old_hacker 12 "Artix (без systemd)"; add minimalist 4 "выбор init" ;;
  *EndeavourOS*) add old_hacker 8 "EndeavourOS"; add fresh_witness 6 "близко к Arch (rolling)" ;;
  *Manjaro*)     add fresh_witness 5 "Manjaro (rolling)"; add minimalist 4 "удобный Arch" ;;
  *Arch*)        add old_hacker 8 "Arch Linux"; add fresh_witness 6 "rolling-release"; add programmer 4 "контроль над окружением" ;;
  # ── Атомарные / иммутабельные ──
  *Bazzite*)                        add atomic 12 "Bazzite (атомарная)"; add gamer 12 "Bazzite" ;;
  *Silverblue*|*Kinoite*|*Sericea*) add atomic 14 "атомарная Fedora"; add fresh_witness 8 "immutable" ;;
  *MicroOS*|*Aeon*|*Kalpa*)         add atomic 14 "openSUSE MicroOS/Aeon"; add devops 5 "transactional-update" ;;
  *Vanilla*OS*|*VanillaOS*)         add atomic 14 "Vanilla OS"; add fresh_witness 6 "immutable" ;;
  *blendOS*)                        add atomic 12 "blendOS"; add fresh_witness 6 "мульти-дистро" ;;
  *"GNOME OS"*)                     add atomic 12 "GNOME OS"; add fresh_witness 6 "immutable" ;;
  *SteamOS*)                        add atomic 8 "SteamOS"; add gamer 12 "SteamOS (Steam Deck)" ;;
  *Endless*)                        add atomic 8 "Endless OS"; add minimalist 4 "из коробки" ;;
  # ── Fedora / игровые ──
  *Nobara*)      add gamer 12 "Nobara"; add fresh_witness 6 "Fedora для игр" ;;
  *Fedora*)      add fresh_witness 12 "Fedora"; add programmer 8 "свежий инструментарий" ;;
  # ── Debian / Ubuntu семейство ──
  *Pop\!_OS*)    add gamer 10 "Pop!_OS"; add programmer 5 "удобство из коробки" ;;
  *elementary*)  add ricer 12 "elementary OS" ;;
  *Zorin*)       add ricer 8 "Zorin OS" ;;
  *Mint*)        add minimalist 6 "Linux Mint"; add sysadmin 4 "консерватизм" ;;
  *MX\ *|*MX_*)  add old_hacker 12 "MX Linux"; add sysadmin 4 "antiX-корни" ;;
  *Ubuntu*)      add programmer 8 "Ubuntu"; add devops 4 "стандарт индустрии" ;;
  *Debian*)      add old_hacker 10 "Debian"; add sysadmin 10 "стабильность серверов" ;;
  # ── Security-ориентированные ──
  *Kali*)        add pentester 18 "Kali Linux"; add anonymous 6 "арсенал аудита" ;;
  *Parrot*)      add pentester 16 "Parrot OS"; add anonymous 10 "приватность + аудит" ;;
  *Tails*)       add anonymous 20 "Tails" ;;
  *Qubes*)       add anonymous 16 "Qubes OS"; add pentester 6 "изоляция по доменам" ;;
  *Whonix*)      add anonymous 18 "Whonix" ;;
  # ── Инженерные / нишевые ──
  *openSUSE*)    add sysadmin 10 "openSUSE/YaST"; add devops 6 "корпоративный баланс" ;;
  *NixOS*)       add atomic 14 "NixOS (декларативная)"; add fresh_witness 8 "NixOS"; add old_hacker 4 "тинкеринг" ;;
  *Gentoo*)      add old_hacker 14 "Gentoo"; add programmer 5 "сборка из исходников" ;;
  *Slackware*)   add old_hacker 16 "Slackware"; add minimalist 5 "классика" ;;
  *Void*)        add old_hacker 12 "Void Linux"; add minimalist 8 "runit" ;;
  *Alpine*)      add minimalist 15 "Alpine"; add sysadmin 6 "musl + busybox"; add anonymous 4 "малая поверхность атаки" ;;
  *)             add old_hacker 4 "нестандартный дистрибутив" ;;
esac

# ──────────────────────────────────────────────
# DE / WM
# ──────────────────────────────────────────────

DESKTOP="${XDG_CURRENT_DESKTOP:-}|${DESKTOP_SESSION:-}|${XDG_SESSION_DESKTOP:-}"

case "$DESKTOP" in
  *KDE*|*plasma*)    add ricer 10 "KDE Plasma"; add devops 4 "гибкое окружение" ;;
  *GNOME*|*gnome*)   add minimalist 8 "GNOME"; add ricer 4 "цельный дизайн" ;;
  *Cinnamon*)        add minimalist 6 "Cinnamon" ;;
  *MATE*|*mate*)     add old_hacker 10 "MATE" ;;
  *XFCE*|*xfce*)     add minimalist 8 "XFCE"; add old_hacker 5 "лёгкость" ;;
  *LXQt*|*LXDE*)     add minimalist 10 "LXQt/LXDE" ;;
  *Budgie*)          add ricer 8 "Budgie" ;;
  *Pantheon*)        add ricer 10 "Pantheon" ;;
  *Deepin*|*deepin*) add ricer 10 "Deepin" ;;
  *COSMIC*|*cosmic*) add fresh_witness 12 "COSMIC"; add programmer 4 "окружение на Rust" ;;
esac

# Оконные менеджеры (могут сосуществовать с пустым XDG_CURRENT_DESKTOP)
case "$DESKTOP" in
  *Hyprland*|*hyprland*) add ricer 12 "Hyprland"; add fresh_witness 8 "ricing на Wayland" ;;
  *niri*)               add import_substituted 10 "niri (отеч. разработка)"; add fresh_witness 10 "скроллируемый WM на Wayland" ;;
  *sway*)               add minimalist 12 "sway"; add fresh_witness 5 "Wayland"; add old_hacker 3 "конфиг как код" ;;
  *river*)              add fresh_witness 12 "river (Wayland-WM)" ;;
  *i3*)                 add minimalist 12 "i3"; add old_hacker 5 "тайлинг" ;;
  *bspwm*)              add minimalist 10 "bspwm"; add programmer 5 "скриптуемый WM" ;;
  *dwm*)                add old_hacker 14 "dwm (suckless)"; add minimalist 6 "патчи и пересборка" ;;
  *awesome*)           add programmer 8 "awesome (конфиг на Lua)"; add old_hacker 4 "тайлинг" ;;
  *qtile*)             add programmer 10 "qtile (конфиг на Python)" ;;
  *xmonad*)            add programmer 12 "xmonad (конфиг на Haskell)"; add old_hacker 4 "тайлинг" ;;
  *)                   : ;;
esac

# Если ни DE, ни WM не определились — вероятно, голая консоль/сервер
if [[ "$DESKTOP" == "||" ]]; then
  add old_hacker 6 "без графического окружения"
  add sysadmin 6 "headless-режим"
fi

# ──────────────────────────────────────────────
# Wayland / X11
# ──────────────────────────────────────────────

SESSION_TYPE="${XDG_SESSION_TYPE:-tty}"
case "$SESSION_TYPE" in
  wayland) add fresh_witness 10 "Wayland" ;;
  x11)     add old_hacker 5 "X11" ;;
  tty)     add old_hacker 5 "чистый TTY"; add sysadmin 3 "без иксов" ;;
esac

# ──────────────────────────────────────────────
# Железо
# ──────────────────────────────────────────────

RAM_KB=$(awk '/MemTotal/ {print $2}' /proc/meminfo 2>/dev/null || echo 0)
RAM=$(( RAM_KB / 1024 / 1024 ))
CPU=$(nproc 2>/dev/null || echo 1)

if safe_ge "$RAM" 64; then add devops 5 "64+ GB RAM"; add gamer 6 "флагманское железо"
elif safe_ge "$RAM" 32; then add devops 3 "32+ GB RAM"; add gamer 5 "мощное железо"
elif safe_le "$RAM" 4;  then add minimalist 10 "≤4 GB RAM"
elif safe_le "$RAM" 8;  then add minimalist 5 "8 GB RAM"
fi

# Число ядер — слабый сигнал (не делает человека devops); только мощная станция
if safe_ge "$CPU" 16; then add gamer 4 "16+ ядер (мощная станция)"; fi

# ──────────────────────────────────────────────
# Диск
# ──────────────────────────────────────────────

USED=$(df / 2>/dev/null | awk 'NR==2{gsub(/%/,"",$5); print $5}' || echo 50)

if safe_gt "$USED" 85; then add old_hacker 6 "диск почти полон"
elif safe_lt "$USED" 30; then add minimalist 5 "чистота файловой системы"
fi

# Шифрование диска (LUKS)
if [[ -f /etc/crypttab && -s /etc/crypttab ]] || lsblk -o TYPE 2>/dev/null | grep -q crypt; then
  add anonymous 10 "шифрование диска (LUKS)"
fi

# RAID / LVM → классический системный администратор
if has mdadm && [[ -f /proc/mdstat ]] && grep -q '^md' /proc/mdstat 2>/dev/null; then
  add sysadmin 6 "программный RAID (mdadm)"
fi
if has lvs && lvs >/dev/null 2>&1; then add sysadmin 5 "LVM"; fi

# ──────────────────────────────────────────────
# Init
# ──────────────────────────────────────────────

if pidof systemd >/dev/null 2>&1; then
  :   # systemd — init по умолчанию в большинстве дистрибутивов, не считаем сигналом
elif [[ -f /run/runit.stopit ]] || [[ -d /etc/sv ]]; then
  add old_hacker 12 "runit"
elif [[ -d /etc/s6 ]] || has s6-svscan; then
  add old_hacker 12 "s6"
elif [[ -f /etc/dinit ]] || has dinit; then
  add old_hacker 10 "dinit"
elif has openrc || [[ -f /etc/init.d/openrc ]]; then
  add old_hacker 10 "OpenRC"
else
  add old_hacker 8 "альтернативный init"
fi

# ──────────────────────────────────────────────
# Шифрование и приватность → Анонимус
# ──────────────────────────────────────────────

has gpg         && add anonymous 6  "GPG"
has pass        && add anonymous 8  "pass"
has age         && add anonymous 6  "age"
has veracrypt   && add anonymous 14 "VeraCrypt"
has cryptsetup  && add anonymous 5  "cryptsetup"
has tomb        && add anonymous 8  "Tomb"
has keepassxc   && add anonymous 5  "KeePassXC"
has gocryptfs   && add anonymous 6  "gocryptfs"
has encfs       && add anonymous 5  "EncFS"

# Анонимность сети
has tor                 && add anonymous 12 "Tor"
has torbrowser-launcher && add anonymous 10 "Tor Browser"
has i2prouter           && add anonymous 10 "I2P"
has proxychains         && { add anonymous 8 "proxychains"; add pentester 4 "цепочки прокси"; }
has mullvad             && add anonymous 10 "Mullvad VPN"
has protonvpn           && add anonymous 8  "ProtonVPN"
has protonvpn-cli       && add anonymous 8  "ProtonVPN"
has openvpn             && add anonymous 6  "OpenVPN"
if has wg || has wg-quick; then add anonymous 6 "WireGuard"; fi

# ──────────────────────────────────────────────
# Аудит безопасности и взлом → Хакер
# ──────────────────────────────────────────────

has nmap        && add pentester 12 "nmap"
has masscan     && add pentester 8  "masscan"
has wireshark   && add pentester 8  "Wireshark"
has tshark      && add pentester 6  "tshark"
has tcpdump     && add pentester 5  "tcpdump"
has msfconsole  && add pentester 16 "Metasploit"
has aircrack-ng && add pentester 12 "aircrack-ng"
has hashcat     && add pentester 10 "hashcat"
has john        && add pentester 10 "John the Ripper"
has hydra       && add pentester 10 "hydra"
has sqlmap      && add pentester 10 "sqlmap"
has nikto       && add pentester 8  "nikto"
has gobuster    && add pentester 6  "gobuster"
has ffuf        && add pentester 6  "ffuf"
has burpsuite   && add pentester 12 "Burp Suite"
has zaproxy     && add pentester 8  "OWASP ZAP"
has radare2     && add pentester 10 "radare2"
has r2          && add pentester 10 "radare2"
has ghidra      && add pentester 12 "Ghidra"
has binwalk     && add pentester 6  "binwalk"
has volatility  && add pentester 8  "Volatility"
has wpscan      && add pentester 6  "WPScan"

# Защитные механизмы. SELinux/AppArmor/файрвол включены ПО УМОЛЧАНИЮ во многих
# дистрибутивах (Fedora, Ubuntu, RHEL), поэтому сами по себе не считаются
# сигналом. Учитываем только то, что ставят и настраивают осознанно.
has fail2ban && add sysadmin 6 "fail2ban"

# ──────────────────────────────────────────────
# Контейнеры, оркестрация, IaC → DevOps (эфемерное, декларативное)
# ──────────────────────────────────────────────

has docker     && add devops 5 "Docker (установлен)"
has podman     && { add devops 5 "Podman (установлен)"; add anonymous 3 "rootless-контейнеры"; }
has kubectl    && add devops 14 "Kubernetes"
has k9s        && add devops 8  "k9s"
has helm       && add devops 8  "Helm"
if has minikube || has kind; then add devops 6 "локальный кластер"; fi
if has argocd || has flux || has skaffold; then add devops 8 "GitOps"; fi
has terraform  && add devops 12 "Terraform"
has opentofu   && add devops 10 "OpenTofu"
has pulumi     && add devops 8  "Pulumi"
has ansible    && { add devops 8 "Ansible"; add sysadmin 4 "автоматизация"; }
if has puppet || has chef || has salt; then add devops 8 "config management"; fi
has vagrant    && add devops 6 "Vagrant"
if has qemu-img || has virt-manager || has virsh; then add devops 8 "QEMU/KVM"; fi
if has lxc || has lxd || has incus; then add devops 6 "LXC/Incus"; fi
if has vault || has consul || has nomad; then add devops 8 "HashiCorp-стек"; fi
if has aws || has gcloud || has az || has doctl; then add devops 8 "облачный CLI"; fi
if has gitlab-runner || has act; then add devops 5 "CI-раннеры"; fi

# ──────────────────────────────────────────────
# Серверы, демоны, бэкапы → Системный администратор (живые pets)
# ──────────────────────────────────────────────

if has nginx || [[ -d /etc/nginx ]]; then
  add sysadmin 8 "nginx"; add import_substituted 8 "nginx (Игорь Сысоев)"
fi
if has apache2 || has httpd; then add sysadmin 6 "Apache"; add old_hacker 4 "httpd"; fi
if has psql || has postgres || [[ -d /var/lib/pgsql ]] || [[ -d /var/lib/postgresql ]]; then
  add sysadmin 8 "PostgreSQL"; add import_substituted 8 "PostgreSQL (Postgres Pro)"
fi
if has mysql || has mariadb; then add sysadmin 6 "MySQL/MariaDB"; fi
has redis-cli && add sysadmin 5 "Redis"
has mongod    && add devops 5 "MongoDB"
if has sshd || [[ -f /etc/ssh/sshd_config ]]; then add sysadmin 5 "SSH-сервер"; fi
if has htop || has btop || has glances; then add sysadmin 3 "мониторинг процессов"; fi
has prometheus && { add sysadmin 6 "Prometheus"; add devops 3 "метрики"; }
has grafana    && { add sysadmin 6 "Grafana"; add devops 3 "дашборды"; }
# Бэкапы — визитная карточка классического админа (ставятся осознанно)
if has borg || has restic || has rsnapshot || has duplicity; then add sysadmin 8 "бэкапы (borg/restic)"; fi
has smartctl  && add sysadmin 3 "S.M.A.R.T.-мониторинг"
# rsync, cron, logrotate предустановлены почти везде — не сигнал;
# реальное использование ловим в поведенческом анализе ниже.

# ──────────────────────────────────────────────
# Разработка → Программист
# ──────────────────────────────────────────────

# Языки и рантаймы. python/perl/java часто предустановлены или приходят
# зависимостью, поэтому одно лишь их наличие не делает человека программистом.
# Засчитываем осознанно установленные рантаймы и реальные следы разработки.
dev_count=0
for lang in node deno bun rustc go zig elixir julia haskell scala kotlin crystal nim ocaml; do
  if has "$lang"; then add programmer 3 "$lang"; dev_count=$(( dev_count + 1 )); fi
done
# Интерпретатор, который нередко предустановлен — слабый сигнал
has ruby && { add programmer 1 "ruby"; dev_count=$(( dev_count + 1 )); }

# Python есть почти везде по умолчанию: засчитываем только при следах разработки
if has python3 || has python; then
  if has pyenv || has poetry || has pipx || has virtualenv \
     || [[ -d "${HOME:-}/.virtualenvs" ]] \
     || compgen -G "${HOME:-}/.local/lib/python*/site-packages" >/dev/null 2>&1; then
    add programmer 3 "Python + инструменты разработки"; dev_count=$(( dev_count + 1 ))
  fi
fi

# Пользовательские тулчейны и менеджеры версий — явный признак разработчика
[[ -d "${HOME:-}/.cargo"  ]] && { add programmer 5 "Cargo"; add fresh_witness 2 "Rust toolchain"; }
[[ -d "${HOME:-}/.rustup" ]] && add programmer 2 "rustup"
[[ -d "${HOME:-}/go"      ]] && add programmer 3 "Go workspace"
[[ -d "${HOME:-}/.npm"    ]] && add programmer 3 "npm-проекты"
if has asdf || has nvm || has pyenv || has rbenv || has sdk; then add programmer 5 "менеджер версий языков"; fi

# Настроенный git (есть секция [user]) — пользователь реально коммитит
if [[ -f "${HOME:-}/.gitconfig" ]] && grep -qi '\[user\]' "${HOME:-}/.gitconfig" 2>/dev/null; then
  add programmer 6 "настроенный git (user.*)"
fi

if safe_ge "$dev_count" 5; then add programmer 10 "полиглот (5+ языков)"
elif safe_ge "$dev_count" 3; then add programmer 4 "несколько языков"
fi

if has gcc || has clang; then add old_hacker 3 "gcc/clang"; fi
if has make && has cmake; then add programmer 3 "сборочные системы"; fi
has rustc && add fresh_witness 6 "Rust"
has docker-compose && add devops 4 "Compose"

# Редакторы
if has vim || has nvim; then add old_hacker 6 "Vim/Neovim"; add programmer 4 "модальное редактирование"; fi
has emacs && add old_hacker 12 "Emacs"
has helix && add fresh_witness 6 "Helix"
has code  && add programmer 6 "VS Code"
if has nvim && [[ -d "${HOME:-}/.config/nvim" ]]; then add ricer 5 "кастомный Neovim"; fi
if [[ -d "${HOME:-}/.config/JetBrains" ]] || has idea; then add programmer 6 "JetBrains IDE"; fi

if has tmux || has screen; then add old_hacker 6 "мультиплексор"; fi
if has gdb || has lldb; then add programmer 4 "отладчики"; fi

# ──────────────────────────────────────────────
# Браузеры
# ──────────────────────────────────────────────

if has yandex-browser || has yandex_browser || has yandex-browser-stable; then
  add import_substituted 12 "Яндекс.Браузер"
fi

# Русская локаль — слабый намёк на импортозамещение
LOCALE_ALL="${LANG:-}|${LC_ALL:-}|${LC_CTYPE:-}|${LC_MESSAGES:-}"
if [[ "$LOCALE_ALL" == *ru_RU* || "$LOCALE_ALL" == *ru_* ]]; then
  add import_substituted 3 "русская локаль"
fi
has librewolf       && { add anonymous 6 "LibreWolf"; add ricer 3 "приватный форк"; }
has brave           && add anonymous 4 "Brave"
has mullvad-browser && add anonymous 8 "Mullvad Browser"

# ──────────────────────────────────────────────
# Пакетные менеджеры
# ──────────────────────────────────────────────

has flatpak && add atomic 4 "Flatpak"
has snap    && add minimalist 2 "Snap"
has nix     && { add atomic 8 "Nix"; add fresh_witness 4 "декларативные пакеты"; }
has brew    && add programmer 4 "Homebrew"
# Контейнеры для разработки поверх иммутабельной ОС — характерны для атомарников
if has distrobox || has toolbox || has toolbx; then add atomic 8 "distrobox/toolbox"; fi
if has rpm-ostree || has ostree || has bootc; then add atomic 10 "ostree-система"; fi
if has yay || has paru; then add old_hacker 4 "AUR-хелпер"; fi
if has guix; then add atomic 8 "GNU Guix"; add old_hacker 6 "функциональный пакетинг"; fi

# ──────────────────────────────────────────────
# Игры → Геймер
# ──────────────────────────────────────────────

has steam       && add gamer 12 "Steam"
has lutris      && add gamer 8  "Lutris"
has heroic      && add gamer 6  "Heroic"
has bottles     && add gamer 6  "Bottles"
if has wine || has wine64; then add gamer 6 "Wine"; fi
has gamemoderun && add gamer 6 "GameMode"
has mangohud    && add gamer 5 "MangoHud"
if has protontricks || has protontricks-launch; then add gamer 5 "Proton"; fi
has retroarch   && add gamer 6 "RetroArch (эмуляция)"

# Подсчёт установленных игр по .desktop-файлам категории Game
GAME_DESKTOPS=$(grep -rliE '^Categories=.*Game' \
  /usr/share/applications \
  "${HOME:-}/.local/share/applications" \
  /var/lib/flatpak/exports/share/applications \
  "${HOME:-}/.local/share/flatpak/exports/share/applications" \
  2>/dev/null | wc -l)
if safe_ge "$GAME_DESKTOPS" 15; then add gamer 14 "${GAME_DESKTOPS} игр в меню"
elif safe_ge "$GAME_DESKTOPS" 5; then add gamer 8 "${GAME_DESKTOPS} игр в меню"
elif safe_ge "$GAME_DESKTOPS" 1; then add gamer 3 "${GAME_DESKTOPS} игр в меню"
fi

if [[ -f /proc/modules ]] && grep -q nvidia /proc/modules 2>/dev/null; then
  add gamer 6 "NVIDIA GPU"
fi

# ──────────────────────────────────────────────
# Shell
# ──────────────────────────────────────────────

SHELL_NAME=$(basename "${SHELL:-bash}")
case "$SHELL_NAME" in
  zsh)    add ricer 8 "zsh" ;;
  fish)   add fresh_witness 6 "fish" ;;
  bash)   add old_hacker 3 "bash" ;;
  dash)   add minimalist 8 "dash" ;;
  nu)     add fresh_witness 10 "nushell"; add programmer 3 "структурный shell" ;;
  xonsh)  add fresh_witness 5 "xonsh"; add programmer 4 "Python-shell" ;;
  elvish) add fresh_witness 6 "elvish" ;;
esac

# ──────────────────────────────────────────────
# Терминальные эмуляторы → Райсер / Свидетель свежего ПО
# ──────────────────────────────────────────────

has kitty     && { add ricer 6 "kitty"; add fresh_witness 3 "GPU-терминал"; }
has alacritty && { add ricer 5 "Alacritty"; add fresh_witness 3 "GPU-терминал"; }
has wezterm   && { add ricer 6 "WezTerm"; add programmer 3 "конфиг на Lua"; }
has foot      && { add minimalist 5 "foot"; add fresh_witness 3 "Wayland-терминал"; }
has st        && { add old_hacker 6 "st (suckless)"; add minimalist 4 "патчи под себя"; }
has xterm     && add old_hacker 3 "xterm"

# ──────────────────────────────────────────────
# Dotfiles
# ──────────────────────────────────────────────

# Райсинг определяем по конфигам инструментов кастомизации, а НЕ по числу
# скрытых папок в $HOME — их и так десятки у любого DE (.config, .cache,
# .mozilla, .pki, .gnupg…), из-за чего раньше любой пользователь GNOME
# попадал в «райсеры».
RICE_CONFIGS=0
for cfg in hypr waybar polybar picom compton rofi wofi dunst mako eww \
           sway i3 bspwm awesome qtile river niri \
           kitty alacritty wezterm foot \
           fastfetch neofetch wal swaylock wlogout starship; do
  [[ -e "${HOME:-}/.config/$cfg" ]] && RICE_CONFIGS=$(( RICE_CONFIGS + 1 ))
done
if safe_ge "$RICE_CONFIGS" 6;   then add ricer 12 "кастомных конфигов: ${RICE_CONFIGS}"
elif safe_ge "$RICE_CONFIGS" 3; then add ricer 6  "кастомизация окружения (${RICE_CONFIGS})"
elif [[ "$RICE_CONFIGS" -eq 0 ]]; then add minimalist 5 "без кастомизации окружения"
fi

if [[ -d "${HOME:-}/dotfiles/.git" ]] || [[ -d "${HOME:-}/.dotfiles/.git" ]]; then
  add ricer 8 "dotfiles в Git"
  add programmer 4 "управление конфигами"
fi
if has stow || has chezmoi; then add ricer 5 "менеджер dotfiles"; add devops 3 "воспроизводимые конфиги"; fi

# ──────────────────────────────────────────────
# Состояние системы: аптайм, возраст, пакеты, ядро, flatpak
# ──────────────────────────────────────────────

# Аптайм — длинный характерен для серверов и админов
UPTIME_SEC=$(cut -d. -f1 /proc/uptime 2>/dev/null || echo 0)
UPTIME_DAYS=$(( UPTIME_SEC / 86400 ))
if   safe_ge "$UPTIME_DAYS" 30; then add sysadmin 10 "аптайм ${UPTIME_DAYS} дн."
elif safe_ge "$UPTIME_DAYS" 7;  then add sysadmin 5  "аптайм ${UPTIME_DAYS} дн."
fi

# Возраст установки: birth-время корня, иначе mtime machine-id как запасной вариант
INSTALL_EPOCH=$(stat -c %W / 2>/dev/null || echo 0)
if [[ -z "$INSTALL_EPOCH" || "$INSTALL_EPOCH" == 0 ]]; then
  INSTALL_EPOCH=$(stat -c %Y /etc/machine-id 2>/dev/null || echo 0)
fi
if safe_gt "$INSTALL_EPOCH" 0; then
  AGE_DAYS=$(( ( $(date +%s) - INSTALL_EPOCH ) / 86400 ))
  if   safe_ge "$AGE_DAYS" 1095; then add sysadmin 8 "система живёт ${AGE_DAYS} дн. без переустановки"; add old_hacker 4 "не распыляется на переустановки"
  elif safe_ge "$AGE_DAYS" 365;  then add sysadmin 4 "установлена больше года назад"
  elif safe_ge "$AGE_DAYS" 0 && safe_le "$AGE_DAYS" 14; then add fresh_witness 4 "свежая установка (${AGE_DAYS} дн.)"
  fi
fi

# Число установленных пакетов — мало пакетов = минимализм
PKG_COUNT=0
if   has pacman;     then PKG_COUNT=$(pacman -Qq 2>/dev/null | wc -l)
elif has dpkg-query; then PKG_COUNT=$(dpkg-query -f '.\n' -W 2>/dev/null | wc -l)
elif has rpm;        then PKG_COUNT=$(rpm -qa 2>/dev/null | wc -l)
elif has apk;        then PKG_COUNT=$(apk info 2>/dev/null | wc -l)
elif has xbps-query; then PKG_COUNT=$(xbps-query -l 2>/dev/null | wc -l)
fi
if safe_gt "$PKG_COUNT" 0; then
  if   safe_le "$PKG_COUNT" 300; then add minimalist 10 "очень мало пакетов (${PKG_COUNT})"
  elif safe_le "$PKG_COUNT" 600; then add minimalist 5  "немного пакетов (${PKG_COUNT})"
  fi
fi

# Чужие пакеты (AUR) — практика тинкеров Arch
if has pacman; then
  AUR_COUNT=$(pacman -Qqm 2>/dev/null | wc -l)
  if   safe_ge "$AUR_COUNT" 20; then add old_hacker 6 "${AUR_COUNT} пакетов из AUR"
  elif safe_ge "$AUR_COUNT" 5;  then add old_hacker 3 "сборки из AUR"
  fi
fi

# Кастомное / специальное ядро
case "$(uname -r 2>/dev/null)" in
  *zen*)            add gamer 5 "ядро Zen"; add fresh_witness 3 "тюнинг отзывчивости" ;;
  *xanmod*)         add gamer 5 "ядро XanMod" ;;
  *lqx*|*liquorix*) add gamer 5 "ядро Liquorix" ;;
  *tkg*)            add gamer 5 "ядро TkG"; add old_hacker 3 "сборка ядра" ;;
  *hardened*)       add anonymous 6 "hardened-ядро"; add pentester 3 "защищённое ядро" ;;
  *-rt*|*rt[0-9]*)  add sysadmin 4 "realtime-ядро" ;;
  *lts*)            add sysadmin 4 "LTS-ядро (стабильность)" ;;
esac

# Снапшоты ФС и декларативная конфигурация
if has snapper || [[ -d /.snapshots ]]; then add sysadmin 5 "снапшоты (snapper)"; add atomic 4 "откаты ФС"; fi
has timeshift && add sysadmin 5 "Timeshift"
if [[ -d /etc/nixos ]] || [[ -d "${HOME:-}/.config/home-manager" ]]; then add atomic 6 "Nix/home-manager"; fi
if has zfs || [[ -d /sys/module/zfs ]]; then add sysadmin 6 "ZFS"; add old_hacker 3 "ZFS-энтузиаст"; fi

# SSH-ключи — рабочий инструмент админа/devops/разработчика
if [[ -d "${HOME:-}/.ssh" ]]; then
  SSH_KEYS=$(find "${HOME}/.ssh" -maxdepth 1 -name 'id_*' ! -name '*.pub' 2>/dev/null | wc -l)
  if safe_ge "$SSH_KEYS" 1; then add devops 4 "SSH-ключи"; add sysadmin 3 "доступ к хостам"; fi
fi

# Git-репозитории в домашней директории — реальная разработка
if [[ -d "${HOME:-}" ]]; then
  GIT_REPOS=$(find "$HOME" -maxdepth 4 -name .git -type d 2>/dev/null | head -200 | wc -l)
  if   safe_ge "$GIT_REPOS" 10; then add programmer 8 "${GIT_REPOS}+ git-репозиториев"
  elif safe_ge "$GIT_REPOS" 3;  then add programmer 4 "${GIT_REPOS} git-репозитория"
  fi
fi

# Установленные flatpak-приложения: количество и категории
if has flatpak; then
  FLATPAK_APPS=$(flatpak list --app --columns=application 2>/dev/null)
  FP_COUNT=$(printf '%s\n' "$FLATPAK_APPS" | grep -c .)
  if   safe_ge "$FP_COUNT" 15; then add atomic 8 "${FP_COUNT} flatpak-приложений"
  elif safe_ge "$FP_COUNT" 5;  then add atomic 5 "${FP_COUNT} flatpak-приложений"
  fi
  grep -qiE 'torproject|mullvad|signalapp|briar|protonvpn|monero' <<< "$FLATPAK_APPS" && add anonymous 5 "приватные flatpak"
  grep -qiE 'valvesoftware\.Steam|heroicgameslauncher|net\.lutris|Bottles|RetroArch|prismlauncher' <<< "$FLATPAK_APPS" && add gamer 6 "игровые flatpak"
  grep -qiE 'visualstudio|jetbrains|gnome\.Builder|vscodium|dev\.zed|neovim|GitKraken' <<< "$FLATPAK_APPS" && add programmer 5 "dev-flatpak"
  grep -qiE 'blender|gimp|inkscape|kdenlive|obsproject|darktable|krita|Audacity' <<< "$FLATPAK_APPS" && add ricer 4 "креатив/медиа flatpak"
fi

# ──────────────────────────────────────────────
# Поведенческий анализ: shell-конфиги и история команд
#   Отличаем тех, кто РЕАЛЬНО пользуется инструментом, от тех, у кого он
#   просто установлен (зависимостью или по умолчанию в дистрибутиве).
# ──────────────────────────────────────────────

BEHAVIOR=""
for f in "${HOME:-}/.bash_history" "${HOME:-}/.zsh_history" \
         "${HOME:-}/.local/share/fish/fish_history" \
         "${HOME:-}/.bashrc" "${HOME:-}/.zshrc" "${HOME:-}/.bash_aliases" \
         "${HOME:-}/.config/fish/config.fish" "${HOME:-}/.profile"; do
  [[ -r $f ]] && BEHAVIOR+=$'\n'"$(cat "$f" 2>/dev/null)"
done

# behav PATTERN CLASS PTS REASON THRESHOLD
# Начисляет очки, если команда встречается в истории/конфигах не реже THRESHOLD раз.
behav() {
  local pat=$1 class=$2 pts=$3 reason=$4 thr=$5 n
  [[ -z $BEHAVIOR ]] && return
  n=$(grep -aowE "$pat" <<< "$BEHAVIOR" 2>/dev/null | wc -l)
  n=${n:-0}
  if (( n >= thr )); then add "$class" "$pts" "${reason} (×${n})"; fi
}

if [[ -n $BEHAVIOR ]]; then
  behav 'docker|docker-compose'                            devops        7  "docker в работе"             5
  behav 'kubectl|helm|k9s|kustomize'                       devops        9  "kubernetes в работе"         3
  behav 'terraform|tofu|ansible|pulumi'                    devops        8  "IaC в истории"               3
  behav 'systemctl|journalctl'                             sysadmin      7  "управление сервисами"        6
  behav 'ssh|scp|sftp'                                     sysadmin      5  "удалённые хосты"             8
  behav 'nginx|certbot|iptables|nft|ufw'                   sysadmin      6  "серверная эксплуатация"      4
  behav 'psql|mysql|mariadb|redis-cli'                     sysadmin      5  "работа с БД"                 4
  behav 'make|cargo|npm|pnpm|yarn|pip|pip3|gradle|mvn'     programmer    7  "сборка/разработка"           6
  behav 'git'                                              programmer    6  "git в повседневной работе"  12
  behav 'vim|nvim|emacs'                                   programmer    4  "редактор кода в истории"    10
  behav 'nmap|nikto|sqlmap|msfconsole|hydra|hashcat|aircrack-ng|gobuster' pentester 11 "пентест в истории" 2
  behav 'tor|proxychains|openvpn|wg|gpg|veracrypt'         anonymous     6  "приватность в истории"       3
  behav 'pacman|yay|paru|makepkg|emerge'                   old_hacker    4  "ручное управление пакетами"  8
  behav 'flatpak|distrobox|toolbox|rpm-ostree|nixos-rebuild|nix-shell|nix-env' atomic 7 "иммутабельный workflow" 3
  behav 'nix'                                              fresh_witness 4  "Nix в истории"               3
  behav 'steam|lutris|wine|proton|protontricks'            gamer         6  "игры в истории"              2
fi

# ──────────────────────────────────────────────
# Алиасы для fetch-утилит (fastfetch и аналоги)
#   Многие заводят `alias ff='fastfetch --config ...'`. Считаем такой алиас
#   равнозначным fastfetch (флаги/конфиги игнорируем — это просто fastfetch),
#   а само наличие красивого fetch-алиаса — лёгкий сигнал райсинга.
# ──────────────────────────────────────────────

FETCH_RE='fastfetch|neofetch|screenfetch|pfetch|hyfetch|nerdfetch|macchina|cpufetch'
FF_MATCH_RE="$FETCH_RE"
FF_ALIAS_LABEL=""
declare -a FF_ALIASES=()

_rc_for_alias=$(
  for f in "${HOME:-}/.bashrc" "${HOME:-}/.bash_aliases" "${HOME:-}/.zshrc" \
           "${HOME:-}/.zsh_aliases" "${HOME:-}/.aliases" \
           "${HOME:-}/.config/fish/config.fish"; do
    [[ -r $f ]] && cat "$f" 2>/dev/null
  done
)
# alias name=...fetch  /  alias name 'fetch' (fish)
while IFS= read -r _an; do
  [[ -n $_an ]] && FF_ALIASES+=("$_an")
done < <(
  grep -hoiE "alias[[:space:]]+[A-Za-z0-9_.-]+[[:space:]]*=?[[:space:]]*['\"]?(${FETCH_RE})" <<< "$_rc_for_alias" \
    | sed -E "s/^[Aa][Ll][Ii][Aa][Ss][[:space:]]+([A-Za-z0-9_.-]+).*/\1/"
)
# однострочные функции: name() { ... fastfetch
while IFS= read -r _an; do
  [[ -n $_an ]] && FF_ALIASES+=("$_an")
done < <(
  grep -hoiE "[A-Za-z0-9_.-]+[[:space:]]*\(\)[[:space:]]*\{[^}]*(${FETCH_RE})" <<< "$_rc_for_alias" \
    | sed -E "s/^([A-Za-z0-9_.-]+).*/\1/"
)

if (( ${#FF_ALIASES[@]} > 0 )); then
  add ricer 4 "fetch-алиас (${FF_ALIASES[0]})"
  FF_ALIAS_LABEL=" +${FF_ALIASES[0]}"
  for _a in "${FF_ALIASES[@]}"; do
    FF_MATCH_RE+="|$(printf '%s' "$_a" | sed 's/\./\\./g')"
  done
fi

# ──────────────────────────────────────────────
# Нормализация и сортировка
# ──────────────────────────────────────────────

MAX_SCORE=1
for key in "${!score[@]}"; do
  if safe_gt "${score[$key]}" "$MAX_SCORE"; then
    MAX_SCORE=${score[$key]}
  fi
done

declare -A norm_score=()
for key in "${!score[@]}"; do
  norm_score[$key]=$(( score[$key] * 100 / MAX_SCORE ))
done

sorted_keys=()
while IFS= read -r line; do
  sorted_keys+=("${line#* }")
done < <(
  for key in "${!norm_score[@]}"; do
    printf '%d %s\n' "${norm_score[$key]}" "$key"
  done | sort -rn
)
