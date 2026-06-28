#!/usr/bin/env bash
# install.sh — установка lxprofiler одной командой.
#
#   curl -fsSL https://raw.githubusercontent.com/Axawys/lxprofiler/main/install.sh | bash
#
# Клонирует репозиторий в ~/.local/share/lxprofiler и создаёт симлинк
# ~/.local/bin/lxprofile. Повторный запуск обновляет установку.
#
# Переопределяемые переменные окружения:
#   LXPROFILE_HOME — куда клонировать (по умолчанию ~/.local/share/lxprofiler)
#   LXPROFILE_BIN  — куда положить symlink (по умолчанию ~/.local/bin)

set -euo pipefail

REPO_URL="${LXPROFILE_REPO:-https://github.com/Axawys/lxprofiler}"
INSTALL_DIR="${LXPROFILE_HOME:-$HOME/.local/share/lxprofiler}"
BIN_DIR="${LXPROFILE_BIN:-$HOME/.local/bin}"

info()  { printf '\033[36m::\033[0m %s\n' "$*"; }
ok()    { printf '\033[32m✓\033[0m %s\n' "$*"; }
err()   { printf '\033[31m✗\033[0m %s\n' "$*" >&2; }

if ! command -v git >/dev/null 2>&1; then
  err "Для установки нужен git. Установите его и повторите."
  exit 1
fi

mkdir -p "$BIN_DIR"

if [[ -d "$INSTALL_DIR/.git" ]]; then
  info "Найдена установка в $INSTALL_DIR — обновляю…"
  git -C "$INSTALL_DIR" pull --ff-only
else
  if [[ -e "$INSTALL_DIR" ]]; then
    err "Каталог $INSTALL_DIR существует, но это не git-репозиторий."
    err "Удалите его или задайте LXPROFILE_HOME и повторите."
    exit 1
  fi
  info "Клонирую $REPO_URL в $INSTALL_DIR…"
  git clone --depth 1 "$REPO_URL" "$INSTALL_DIR"
fi

chmod +x "$INSTALL_DIR/lxprofile"
ln -sf "$INSTALL_DIR/lxprofile" "$BIN_DIR/lxprofile"
ok "Команда lxprofile установлена → $BIN_DIR/lxprofile"

# Проверяем, что BIN_DIR в PATH
case ":$PATH:" in
  *":$BIN_DIR:"*)
    ok "Готово! Запустите: lxprofile"
    ;;
  *)
    printf '\n'
    err "Каталог $BIN_DIR не в \$PATH."
    printf '   Добавьте в ~/.bashrc или ~/.zshrc строку:\n\n'
    printf '     export PATH="%s:$PATH"\n\n' "$BIN_DIR"
    printf '   затем перезапустите терминал и выполните: lxprofile\n'
    ;;
esac
