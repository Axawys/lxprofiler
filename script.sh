#!/usr/bin/env bash
# linux_profile.sh — Linux Psychological Profiler v4.0
# Requires: bash 4.0+, standard coreutils, tput (ncurses)

# Без set -e: grep, pidof, (( )) штатно возвращают ненулевой код
set -uo pipefail

declare -A LABEL=(
  [devops]="DevOps-инженер"
  [programmer]="Программист"
  [sysadmin]="Системный администратор"
  [minimalist]="Минималист"
  [old_hacker]="Бородатый юниксоид"
  [ricer]="Райсер"
  [gamer]="Геймер"
  [anonymous]="Анонимус"
  [pentester]="Хакер"
  [import_substituted]="Импортозамещённый"
  [fresh_witness]="Свидетель свежего ПО"
)

declare -A DESCRIPTION=(
  [devops]="Cattle, не pets. Кластеры, пайплайны и декларативная инфраструктура важнее отдельной машины."
  [programmer]="Компиляторы, языки, редакторы. Машина для тебя — мастерская, а не витрина."
  [sysadmin]="Pets, не cattle. Живые серверы и демоны, аптайм в годах — ты знаешь каждую машину по имени."
  [minimalist]="Ты ценишь тишину. Меньше процессов — больше смысла."
  [old_hacker]="Ты помнишь, как всё должно работать, и был здесь до systemd. Философия Unix — твоя привычка."
  [ricer]="Рабочий стол — холст. Конфиги вылизаны, шрифты идеальны, скриншот готов для r/unixporn."
  [gamer]="Linux — это не только работа. Proton запущен, пингвин тащит твои игры."
  [anonymous]="Тебя здесь не было. Tor, VPN и шифрование заметают следы за тобой."
  [pentester]="Ты знаешь, где сломается чужая система. nmap греется не просто так."
  [import_substituted]="Сделано в России. Отечественный софт — твой осознанный выбор."
  [fresh_witness]="Ты крестишься на свежие релизы. Если версия не последняя — это уже легаси."
)

declare -A score=()
declare -A reasons=()
for key in "${!LABEL[@]}"; do
  score[$key]=0
  reasons[$key]=""
done

# ──────────────────────────────────────────────
# Вспомогательные функции
# ──────────────────────────────────────────────

add() {
  local key=$1 pts=$2 reason=$3
  score[$key]=$(( score[$key] + pts ))
  reasons[$key]+="${reasons[$key]:+, }${reason}"
}

has() {
  command -v "$1" >/dev/null 2>&1
}

# (( expr )) возвращает 1 при нуле — оборачиваем, чтобы не падать под pipefail
safe_gt() { [[ $(( $1 > $2 )) -eq 1 ]]; }
safe_lt() { [[ $(( $1 < $2 )) -eq 1 ]]; }
safe_ge() { [[ $(( $1 >= $2 )) -eq 1 ]]; }
safe_le() { [[ $(( $1 <= $2 )) -eq 1 ]]; }

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
  # ── Fedora / атомарные / игровые ──
  *Bazzite*)     add gamer 14 "Bazzite"; add fresh_witness 8 "атомарная ОС" ;;
  *Nobara*)      add gamer 12 "Nobara"; add fresh_witness 6 "Fedora для игр" ;;
  *Silverblue*|*Kinoite*) add fresh_witness 12 "атомарная Fedora"; add devops 6 "immutable" ;;
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
  *NixOS*)       add devops 14 "декларативность NixOS"; add fresh_witness 10 "NixOS" ;;
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

if safe_ge "$RAM" 64; then add devops 10 "64+ GB RAM"; add gamer 6 "флагманское железо"
elif safe_ge "$RAM" 32; then add devops 6 "32+ GB RAM"; add gamer 5 "мощное железо"
elif safe_le "$RAM" 4;  then add minimalist 10 "≤4 GB RAM"
elif safe_le "$RAM" 8;  then add minimalist 5 "8 GB RAM"
fi

if safe_ge "$CPU" 16; then add devops 8 "16+ ядер"; add gamer 4 "много потоков"
elif safe_ge "$CPU" 8; then add devops 4 "8+ ядер"
fi

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
  add sysadmin 8 "systemd"
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

# Защитные механизмы → Системный администратор (не «паранойя»)
has fail2ban && add sysadmin 6 "fail2ban"
if has ufw || has firewalld || has nft; then add sysadmin 3 "файрвол"; fi
if [[ -d /sys/fs/selinux ]]; then add sysadmin 5 "SELinux"
elif [[ -d /sys/kernel/security/apparmor ]]; then add sysadmin 4 "AppArmor"
fi

# ──────────────────────────────────────────────
# Контейнеры, оркестрация, IaC → DevOps (эфемерное, декларативное)
# ──────────────────────────────────────────────

has docker     && add devops 10 "Docker"
has podman     && { add devops 10 "Podman"; add anonymous 3 "rootless-контейнеры"; }
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
if has htop || has btop || has glances; then add sysadmin 4 "мониторинг процессов"; fi
has prometheus && { add sysadmin 6 "Prometheus"; add devops 3 "метрики"; }
has grafana    && { add sysadmin 6 "Grafana"; add devops 3 "дашборды"; }
# Бэкапы — визитная карточка классического админа
if has borg || has restic || has rsnapshot || has duplicity; then add sysadmin 8 "бэкапы (borg/restic)"; fi
has rsync     && add sysadmin 3 "rsync"
has smartctl  && add sysadmin 4 "S.M.A.R.T.-мониторинг"
if has crontab && [[ -d /var/spool/cron || -f /etc/crontab ]]; then add sysadmin 4 "cron-задачи"; fi
has logrotate && add sysadmin 3 "logrotate"

# ──────────────────────────────────────────────
# Разработка → Программист
# ──────────────────────────────────────────────

dev_count=0
for lang in python3 python node deno bun rustc go java ruby php elixir julia lua perl zig haskell scala kotlin; do
  if has "$lang"; then
    add programmer 3 "$lang"
    dev_count=$(( dev_count + 1 ))
  fi
done

if safe_ge "$dev_count" 6; then add programmer 12 "полиглот (6+ языков)"
elif safe_ge "$dev_count" 3; then add programmer 5 "несколько языков"
fi

if has gcc || has clang; then add programmer 4 "C/C++"; add old_hacker 4 "gcc/clang"; fi
if has make || has cmake; then add programmer 4 "сборочные системы"; fi
has rustc && { add fresh_witness 8 "Rust"; add programmer 4 "memory-safety"; }
has go    && { add devops 5 "Go"; add programmer 3 "Go"; }
has git   && add programmer 5 "Git"
has docker-compose && add devops 5 "Compose"

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
has librewolf       && { add anonymous 6 "LibreWolf"; add ricer 3 "приватный форк"; }
has brave           && add anonymous 4 "Brave"
has mullvad-browser && add anonymous 8 "Mullvad Browser"

# ──────────────────────────────────────────────
# Пакетные менеджеры
# ──────────────────────────────────────────────

has flatpak && add minimalist 3 "Flatpak"
has snap    && add programmer 2 "Snap"
has nix     && { add devops 8 "Nix"; add fresh_witness 4 "декларативные пакеты"; }
has brew    && add programmer 4 "Homebrew"
if has yay || has paru; then add old_hacker 4 "AUR-хелпер"; fi
if has guix; then add old_hacker 8 "GNU Guix"; add fresh_witness 6 "функциональный пакетинг"; fi

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

if [[ -d "${HOME:-}" ]]; then
  DOTS=$(find "$HOME" -maxdepth 1 -name ".*" ! -name ".." 2>/dev/null | wc -l)
  if safe_gt "$DOTS" 30; then add ricer 10 "30+ dotfiles"
  elif safe_gt "$DOTS" 20; then add ricer 7 "много dotfiles"
  elif safe_le "$DOTS" 5; then add minimalist 8 "минимум dotfiles"
  fi
fi

if [[ -d "${HOME:-}/dotfiles/.git" ]] || [[ -d "${HOME:-}/.dotfiles/.git" ]]; then
  add ricer 8 "dotfiles в Git"
  add programmer 4 "управление конфигами"
fi
if has stow || has chezmoi; then add ricer 5 "менеджер dotfiles"; add devops 3 "воспроизводимые конфиги"; fi

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

# ──────────────────────────────────────────────
# Оформление
# ──────────────────────────────────────────────

BOLD=$'\033[1m'; DIM=$'\033[2m'; RESET=$'\033[0m'
GREEN=$'\033[32m'; YELLOW=$'\033[33m'; CYAN=$'\033[36m'

make_bar() {
  local p=$1 filled=$(( $1 / 5 )) i=0 b=""
  while (( i < filled )); do b+="█"; i=$(( i + 1 )); done
  while (( i < 20 ));     do b+="░"; i=$(( i + 1 )); done
  printf '%s' "$b"
}

# ──────────────────────────────────────────────
# Сборка прокручиваемого тела (BODY) для интерактивного просмотра
# ──────────────────────────────────────────────

declare -a WRAPPED=()

# wrap_into TEXT WIDTH — переносит текст по словам, кладёт строки в массив WRAPPED
wrap_into() {
  WRAPPED=()
  local text=$1 width=$2
  local -a words=()
  read -ra words <<< "$text"
  local cur="" w
  for w in "${words[@]}"; do
    if [[ -z $cur ]]; then
      cur=$w
    elif (( ${#cur} + 1 + ${#w} > width )); then
      WRAPPED+=("$cur")
      cur=$w
    else
      cur="$cur $w"
    fi
  done
  [[ -n $cur ]] && WRAPPED+=("$cur")
}

# ──────────────────────────────────────────────
# Статический вывод (для пайпов / не-TTY)
# ──────────────────────────────────────────────

print_static() {
  echo
  echo "${BOLD}╔══════════════════════════════════════════╗${RESET}"
  echo "${BOLD}║     Linux Psychological Profiler v4.0    ║${RESET}"
  echo "${BOLD}╚══════════════════════════════════════════╝${RESET}"
  echo "${DIM}  Дистрибутив : ${DISTRO}${RESET}"
  echo "${DIM}  Ядро        : $(uname -r 2>/dev/null || echo '?')${RESET}"
  echo "${DIM}  Init        : $(ps -p 1 -o comm= 2>/dev/null || echo '?')${RESET}"
  echo

  local maxlen=0 key label len
  for key in "${sorted_keys[@]}"; do
    label="${LABEL[$key]}"; len=${#label}
    (( len > maxlen )) && maxlen=$len
  done

  local p color pad
  for key in "${sorted_keys[@]}"; do
    p=${norm_score[$key]}; label="${LABEL[$key]}"; len=${#label}
    if safe_ge "$p" 80; then color=$GREEN
    elif safe_ge "$p" 50; then color=$YELLOW
    else color=$DIM; fi
    pad=$(( maxlen - len ))
    printf "${color}%s%*s${RESET}  %3d%%  ${color}%s${RESET}\n" "$label" "$pad" "" "$p" "$(make_bar "$p")"
  done

  local W=${sorted_keys[0]} S=${sorted_keys[1]} T=${sorted_keys[2]}
  echo
  echo "${BOLD}${GREEN}▶ Ты — ${LABEL[$W]}${RESET}"
  echo "  ${DESCRIPTION[$W]}"
  echo "${BOLD}Что повлияло:${RESET} ${DIM}${reasons[$W]}${RESET}"
  echo
  echo "${DIM}  2. ${LABEL[$S]} (${norm_score[$S]}%) — ${reasons[$S]:-—}${RESET}"
  echo "${DIM}  3. ${LABEL[$T]} (${norm_score[$T]}%) — ${reasons[$T]:-—}${RESET}"
  echo
}

# ──────────────────────────────────────────────
# Интерактивный просмотрщик критериев
# ──────────────────────────────────────────────

REPLY_KEY=""
read_key() {
  local k="" rest="" tilde=""
  IFS= read -rsn1 k
  if [[ $k == $'\x1b' ]]; then
    IFS= read -rsn2 -t 0.05 rest 2>/dev/null
    k+=$rest
    if [[ $rest == '[5' || $rest == '[6' ]]; then
      IFS= read -rsn1 -t 0.05 tilde 2>/dev/null
      k+=$tilde
    fi
  fi
  REPLY_KEY=$k
}

_view_cleanup() {
  tput cnorm 2>/dev/null
  tput rmcup 2>/dev/null
}

render_frame() {
  local sel=$1 i sk lbl p pad used l
  tput cup 0 0 2>/dev/null
  tput ed   2>/dev/null

  # Шапка — та же, что в статическом выводе; остаётся на месте
  printf '%s\n' "${BOLD}╔══════════════════════════════════════════╗${RESET}"
  printf '%s\n' "${BOLD}║     Linux Psychological Profiler v4.0    ║${RESET}"
  printf '%s\n' "${BOLD}╚══════════════════════════════════════════╝${RESET}"
  printf '%s\n' "${DIM}  Дистрибутив : ${DISTRO}${RESET}"
  printf '%s\n' "${DIM}  Ядро        : ${KERNEL}  ·  Init : ${INIT1}${RESET}"
  printf '\n'

  # Список классов статичен; меняется только зелёный маркер выделения
  for i in "${!sorted_keys[@]}"; do
    sk=${sorted_keys[i]}
    lbl=${LABEL[$sk]}
    p=${norm_score[$sk]}
    pad=$(( MAXLEN - ${#lbl} ))
    if (( i == sel )); then
      printf "${GREEN}${BOLD}▶ %s%*s  %3d%%  %s${RESET}\n" "$lbl" "$pad" "" "$p" "$(make_bar "$p")"
    else
      printf "${DIM}  %s%*s  %3d%%  %s${RESET}\n" "$lbl" "$pad" "" "$p" "$(make_bar "$p")"
    fi
  done

  printf '%s\n' "${DIM}  ────────────────────────────────────────────${RESET}"

  # Нижняя панель обновляется под выбранный класс (высота фиксирована)
  sk=${sorted_keys[sel]}
  used=0
  printf "${BOLD}${GREEN}▶ %s — %d%%${RESET}\n" "${LABEL[$sk]}" "${norm_score[$sk]}"; used=$(( used + 1 ))
  wrap_into "${DESCRIPTION[$sk]}" "$WRAP_W"
  for l in "${WRAPPED[@]}"; do printf "  %s\n" "$l"; used=$(( used + 1 )); done
  printf '\n'; used=$(( used + 1 ))
  printf "${BOLD}  Что повлияло:${RESET}\n"; used=$(( used + 1 ))
  wrap_into "${reasons[$sk]:-—}" "$WRAP_W"
  for l in "${WRAPPED[@]}"; do printf "${DIM}  %s${RESET}\n" "$l"; used=$(( used + 1 )); done
  while (( used < DETAIL_H )); do printf '\n'; used=$(( used + 1 )); done

  printf '%s\n' "${DIM}  ────────────────────────────────────────────${RESET}"
  printf '%s' "${DIM}  ${BOLD}↑/↓${RESET}${DIM} или ${BOLD}j/k${RESET}${DIM} — двигать маркер · ${BOLD}g/G${RESET}${DIM} — первый/последний · ${BOLD}${YELLOW}q${RESET}${DIM} — выход${RESET}"
}

# Глобальные параметры раскладки интерактивного вида
MAXLEN=0
WRAP_W=72
DETAIL_H=0
KERNEL=""
INIT1=""

interactive_view() {
  local cols sk d c total l last sel=0

  KERNEL=$(uname -r 2>/dev/null || echo '?')
  INIT1=$(ps -p 1 -o comm= 2>/dev/null || echo '?')

  cols=$(tput cols 2>/dev/null || echo 80)
  WRAP_W=$(( cols - 4 ))
  (( WRAP_W > 76 )) && WRAP_W=76
  (( WRAP_W < 30 )) && WRAP_W=30

  # Ширина колонки меток
  MAXLEN=0
  for sk in "${sorted_keys[@]}"; do
    (( ${#LABEL[$sk]} > MAXLEN )) && MAXLEN=${#LABEL[$sk]}
  done

  # Фиксированная высота нижней панели = максимум по всем классам,
  # чтобы раскладка не «прыгала» при смене выделения
  DETAIL_H=0
  for sk in "${sorted_keys[@]}"; do
    wrap_into "${DESCRIPTION[$sk]}" "$WRAP_W"; d=${#WRAPPED[@]}
    wrap_into "${reasons[$sk]:-—}" "$WRAP_W"; c=${#WRAPPED[@]}
    total=$(( 1 + d + 1 + 1 + c ))   # заголовок + описание + пусто + «Что повлияло:» + критерии
    (( total > DETAIL_H )) && DETAIL_H=$total
  done

  last=$(( ${#sorted_keys[@]} - 1 ))

  tput smcup 2>/dev/null; tput civis 2>/dev/null
  trap '_view_cleanup' EXIT INT TERM

  while true; do
    render_frame "$sel"
    read_key
    case "$REPLY_KEY" in
      q|Q)          break ;;
      j|$'\x1b[B')  (( sel < last )) && sel=$(( sel + 1 )) ;;
      k|$'\x1b[A')  (( sel > 0 ))    && sel=$(( sel - 1 )) ;;
      g)            sel=0 ;;
      G)            sel=$last ;;
    esac
  done

  _view_cleanup
  trap - EXIT INT TERM

  # Оповещение о выходе
  local W=${sorted_keys[0]}
  echo
  echo "${BOLD}${GREEN}▶ Твой профиль: ${LABEL[$W]}${RESET}"
  echo "${DIM}  ${DESCRIPTION[$W]}${RESET}"
  echo "${DIM}  Просмотр закрыт по «q». До встречи! 🐧${RESET}"
  echo
}

# ──────────────────────────────────────────────
# Точка входа
# ──────────────────────────────────────────────

# В настоящем терминале сразу открываем интерактивный режим;
# для пайпа / не-TTY печатаем статическую сводку.
if [[ -t 0 && -t 1 ]] && has tput; then
  interactive_view
else
  print_static
fi
