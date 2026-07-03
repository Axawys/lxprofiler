#!/usr/bin/env bash
# install.sh — установка lxprofiler (Go) одной командой.
#
#   curl -fsSL https://raw.githubusercontent.com/Axawys/lxprofiler/main/install.sh | bash
#
# Скачивает готовый бинарник последнего релиза с GitHub Releases под вашу
# ОС/архитектуру и кладёт его в ~/.local/bin/lxprofile. Ни Go, ни git не нужны.
#
# Переопределяемые переменные окружения:
#   LXPROFILE_BIN     — куда положить бинарник (по умолчанию ~/.local/bin)
#   LXPROFILE_VERSION — какую версию ставить (по умолчанию последняя)

set -euo pipefail

REPO="${LXPROFILE_REPO:-Axawys/lxprofiler}"
BIN_DIR="${LXPROFILE_BIN:-$HOME/.local/bin}"

info() { printf '\033[36m::\033[0m %s\n' "$*"; }
ok()   { printf '\033[32m✓\033[0m %s\n' "$*"; }
err()  { printf '\033[31m✗\033[0m %s\n' "$*" >&2; }

for tool in curl uname; do
  command -v "$tool" >/dev/null 2>&1 || { err "нужен $tool, но он не найден."; exit 1; }
done

# ── Миграция со старой bash-версии ─────────────────────────────
# Раньше lxprofiler ставился как git-клон в ~/.local/share/lxprofiler со
# симлинками на bash-скрипт. Убираем его начисто, чтобы Go-бинарник и короткие
# команды не конфликтовали со старыми симлинками. Идемпотентно.
migrate_from_bash() {
  local share="$HOME/.local/share/lxprofiler"
  local marker="${XDG_STATE_HOME:-$HOME/.local/state}/lxprofiler/lx_setup_done"
  local is_bash=0 name link

  [[ -d "$share/.git" ]] && is_bash=1
  # Наш lxprofile-симлинк ведёт внутрь старого share? — тоже bash-инсталл.
  if [[ -L "$BIN_DIR/lxprofile" ]]; then
    case "$(readlink -f "$BIN_DIR/lxprofile" 2>/dev/null)" in
      "$share"/*) is_bash=1 ;;
    esac
  fi
  [[ $is_bash -eq 1 ]] || return 0

  info "Обнаружена старая bash-версия — убираю её…"
  # Снимаем симлинки, указывающие в старый каталог (пока он ещё на месте).
  for name in lxprofile lx lxu lxs lxv lxh lxc lxrm; do
    link="$BIN_DIR/$name"
    [[ -L "$link" ]] || continue
    case "$(readlink -f "$link" 2>/dev/null)" in
      "$share"/*) rm -f "$link" ;;
    esac
  done
  rm -rf "$share"
  rm -f "$marker"   # чтобы Go пересоздал короткие команды на себя
  ok "Старая bash-версия удалена."
}
migrate_from_bash

# ── ОС и архитектура → имя ассета ──────────────────────────────
os=$(uname -s | tr '[:upper:]' '[:lower:]')
case "$os" in
  linux|darwin) : ;;
  *) err "неподдерживаемая ОС: $os"; exit 1 ;;
esac

arch=$(uname -m)
case "$arch" in
  x86_64|amd64)  arch=amd64 ;;
  aarch64|arm64) arch=arm64 ;;
  *) err "неподдерживаемая архитектура: $arch"; exit 1 ;;
esac
asset="lxprofile-${os}-${arch}"

# ── Версия релиза ──────────────────────────────────────────────
if [[ -n "${LXPROFILE_VERSION:-}" ]]; then
  tag="v${LXPROFILE_VERSION#v}"
else
  info "Узнаю последнюю версию…"
  # Без api.github.com (лимит 60/час): берём тег из редиректа веб-страницы
  # github.com/<repo>/releases/latest → .../releases/tag/vX.Y.Z.
  tag=$(curl -fsSL -o /dev/null -w '%{url_effective}' "https://github.com/$REPO/releases/latest" \
        | sed -nE 's#.*/releases/tag/([^/?#]+).*#\1#p')
  [[ -n "$tag" ]] || { err "не удалось узнать последнюю версию (ещё нет релизов?)."; exit 1; }
fi

url="https://github.com/$REPO/releases/download/$tag/$asset"

# ── Скачивание ─────────────────────────────────────────────────
mkdir -p "$BIN_DIR"
info "Скачиваю $asset ($tag)…"
tmp=$(mktemp)
trap 'rm -f "$tmp"' EXIT
if ! curl -fSL "$url" -o "$tmp"; then
  err "не удалось скачать $url"
  err "проверьте, что релиз $tag содержит ассет $asset."
  exit 1
fi
chmod +x "$tmp"
mv "$tmp" "$BIN_DIR/lxprofile"
trap - EXIT
ok "Установлено → $BIN_DIR/lxprofile ($tag)"

# ── Короткие команды lx/lxu/… (создаёт сам бинарник при первом запуске) ──
"$BIN_DIR/lxprofile" -s >/dev/null 2>&1 || true
[[ -L "$BIN_DIR/lx" ]] && ok "Короткая команда: lx (= lxprofile); слитные формы: lxu, lxs, lxc…"

# ── Проверка PATH ──────────────────────────────────────────────
case ":$PATH:" in
  *":$BIN_DIR:"*) ok "Готово! Запустите: lxprofile" ;;
  *)
    printf '\n'
    err "Каталог $BIN_DIR не в \$PATH."
    printf '   Добавьте в ~/.bashrc или ~/.zshrc строку:\n\n'
    printf '     export PATH="%s:$PATH"\n\n' "$BIN_DIR"
    printf '   затем перезапустите терминал и выполните: lxprofile\n'
    ;;
esac
