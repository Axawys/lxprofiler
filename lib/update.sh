#!/usr/bin/env bash
# lib/update.sh — обновление и проверка новых версий.
# Требует переменные LXPROFILE_ROOT, LXPROFILE_LIB и LXPROFILE_VERSION,
# которые задаёт точка входа lxprofile до подключения этого модуля.

# Файл-маркер «пользователь отказался от обновления».
# Пока он существует, предложение обновиться не показывается; снимается только
# ручным обновлением (do_update). Содержит версию, от которой отказались.
LXPROFILE_STATE_DIR="${XDG_STATE_HOME:-$HOME/.local/state}/lxprofiler"
LXPROFILE_DECLINED_FILE="$LXPROFILE_STATE_DIR/update_declined"

# Достаёт номер версии из текста lib/version.sh (с stdin)
_extract_version() {
  grep -m1 '^LXPROFILE_VERSION=' | sed -E 's/.*"([^"]+)".*/\1/'
}

# _version_gt A B → истина, если версия A строго новее B (semver через sort -V)
_version_gt() {
  [[ "$1" != "$2" ]] && [[ "$(printf '%s\n%s\n' "$1" "$2" | sort -V | tail -n1)" == "$1" ]]
}

do_update() {
  if ! command -v git >/dev/null 2>&1; then
    printf 'Для обновления нужен git, но он не найден.\n' >&2
    return 1
  fi
  if ! git -C "$LXPROFILE_ROOT" rev-parse --git-dir >/dev/null 2>&1; then
    printf 'Каталог установки не является git-репозиторием:\n  %s\n' "$LXPROFILE_ROOT" >&2
    printf 'Переустановите утилиту через install.sh.\n' >&2
    return 1
  fi
  printf 'Обновляю lxprofile в %s …\n' "$LXPROFILE_ROOT"
  if ! git -C "$LXPROFILE_ROOT" pull --ff-only; then
    printf 'Не удалось обновить (возможны локальные изменения).\n' >&2
    return 1
  fi
  # Ручное обновление снимает «отказ» — после него снова можно предлагать новинки
  rm -f "$LXPROFILE_DECLINED_FILE" 2>/dev/null
  # Перечитываем версию после обновления
  source "$LXPROFILE_LIB/version.sh"
  printf 'Готово. Текущая версия: %s\n' "$LXPROFILE_VERSION"
}

# Узнаёт версию в репозитории (origin). Печатает её в stdout, либо ничего.
_remote_version() {
  local branch
  branch=$(git -C "$LXPROFILE_ROOT" rev-parse --abbrev-ref HEAD 2>/dev/null)
  [[ -z $branch || $branch == HEAD ]] && branch=main

  # Тихо тянем сведения с origin; с таймаутом, чтобы не висеть оффлайн
  if command -v timeout >/dev/null 2>&1; then
    timeout 5 git -C "$LXPROFILE_ROOT" fetch --quiet origin "$branch" 2>/dev/null || return 1
  else
    git -C "$LXPROFILE_ROOT" fetch --quiet origin "$branch" 2>/dev/null || return 1
  fi
  git -C "$LXPROFILE_ROOT" show "origin/$branch:lib/version.sh" 2>/dev/null | _extract_version
}

# Проверяет наличие новой версии и при необходимости предлагает обновиться.
# y → обновить и перезапуститься; n → больше не предлагать до ручного обновления.
check_for_update() {
  # Только в интерактивном терминале и только из git-установки
  [[ -t 0 && -t 1 ]] || return 0
  command -v git >/dev/null 2>&1 || return 0
  git -C "$LXPROFILE_ROOT" rev-parse --git-dir >/dev/null 2>&1 || return 0

  local remote_ver
  remote_ver=$(_remote_version) || return 0
  [[ -z $remote_ver ]] && return 0

  # Новее ли версия в репозитории?
  _version_gt "$remote_ver" "$LXPROFILE_VERSION" || return 0

  # Пользователь ранее отказался — молчим до ручного обновления
  [[ -e $LXPROFILE_DECLINED_FILE ]] && return 0

  printf 'Доступна новая версия lxprofile: %s (у вас %s).\n' "$remote_ver" "$LXPROFILE_VERSION"
  printf 'Обновить сейчас? [y/N] '
  local ans=""
  read -r ans </dev/tty 2>/dev/null || ans=""
  case "$ans" in
    y|Y|yes|Yes|да|Да|д|Д)
      if do_update; then
        printf 'Перезапуск…\n'
        exec "$LXPROFILE_ROOT/lxprofile"
      fi
      ;;
    *)
      mkdir -p "$LXPROFILE_STATE_DIR" 2>/dev/null
      printf '%s\n' "$remote_ver" > "$LXPROFILE_DECLINED_FILE" 2>/dev/null
      printf 'Хорошо, больше не буду предлагать — пока не обновитесь вручную (lxprofile --update).\n'
      ;;
  esac
}
