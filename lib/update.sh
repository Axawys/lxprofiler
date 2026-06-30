#!/usr/bin/env bash
# lib/update.sh — обновление и проверка новых версий.
# Требует переменные LXPROFILE_ROOT, LXPROFILE_LIB и LXPROFILE_VERSION,
# которые задаёт точка входа lxprofile до подключения этого модуля.

# Файл-маркер «пользователь отказался от обновления».
# Пока он существует, предложение обновиться не показывается; снимается только
# ручным обновлением (do_update). Содержит версию, от которой отказались.
LXPROFILE_STATE_DIR="${XDG_STATE_HOME:-$HOME/.local/state}/lxprofiler"
LXPROFILE_DECLINED_FILE="$LXPROFILE_STATE_DIR/update_declined"
LXPROFILE_LX_MARKER="$LXPROFILE_STATE_DIR/lx_setup_done"

# Короткие команды: lx (= lxprofile) и её «слитные» формы с флагом.
# lxu == lx -u, lxs == lx -s, lxv == lx -v, lxh == lx -h, lxc == lx -c,
# lxrm == lx --rm (разбор имени вызова — в точке входа lxprofile).
LXPROFILE_SHORT_CMDS=(lx lxu lxs lxv lxh lxc lxrm)

# Занято ли короткое имя $1? $2 — путь нашего будущего симлинка.
# Занятым считаем: чужой бинарник/симлинк в PATH (не указывающий на наш lxprofile)
# ИЛИ алиас/функция с этим именем в shell-конфигах (их command -v из скрипта не видит).
_short_taken() {
  local name=$1 link=$2 found tgt our f re
  our=$(readlink -f "$LXPROFILE_ROOT/lxprofile" 2>/dev/null)
  found=$(command -v "$name" 2>/dev/null || true)
  if [[ -n $found && $found != "$link" ]]; then
    tgt=$(readlink -f "$found" 2>/dev/null)
    [[ $tgt != "$our" ]] && return 0   # чужая команда с таким именем
  fi
  re="^[[:space:]]*alias[[:space:]]+${name}[[:space:]]*=|^[[:space:]]*alias[[:space:]]+${name}[[:space:]]+|(^|[[:space:]])function[[:space:]]+${name}([[:space:]]|\\(|\$)|(^|[[:space:]])${name}[[:space:]]*\\(\\)"
  for f in "${HOME:-}/.bashrc" "${HOME:-}/.bash_aliases" "${HOME:-}/.zshrc" \
           "${HOME:-}/.zsh_aliases" "${HOME:-}/.aliases" "${HOME:-}/.profile" \
           "${HOME:-}/.config/fish/config.fish"; do
    [[ -r $f ]] || continue
    grep -qiE "$re" "$f" 2>/dev/null && return 0
  done
  return 1
}

# Создаёт короткие команды (симлинки на lxprofile) для каждого свободного имени
# из LXPROFILE_SHORT_CMDS. Дорогая часть выполняется один раз (маркер
# LXPROFILE_LX_MARKER); $1=force — перепроверить принудительно (при установке
# и обновлении). Сообщения печатаются только про основную команду lx.
ensure_lx() {
  local force="${1:-}" marker="$LXPROFILE_LX_MARKER"
  [[ -z $force && -e $marker ]] && return 0
  mkdir -p "$LXPROFILE_STATE_DIR" 2>/dev/null

  local lxp bin name link our tgt verbose=""
  [[ -e $marker ]] || verbose=1
  lxp=$(command -v lxprofile 2>/dev/null)
  if [[ -n $lxp ]]; then bin=$(dirname "$lxp"); else bin="${LXPROFILE_BIN:-${HOME:-}/.local/bin}"; fi
  our=$(readlink -f "$LXPROFILE_ROOT/lxprofile" 2>/dev/null)

  for name in "${LXPROFILE_SHORT_CMDS[@]}"; do
    link="$bin/$name"

    # уже наш симлинк — ничего не делаем
    if [[ -L $link ]]; then
      tgt=$(readlink -f "$link" 2>/dev/null)
      [[ $tgt == "$our" ]] && continue
    fi

    if _short_taken "$name" "$link"; then
      [[ $name == lx && -n $verbose ]] && printf 'Короткая команда lx уже занята — используйте lxprofile.\n' >&2
    else
      mkdir -p "$bin" 2>/dev/null
      if ln -sf "$LXPROFILE_ROOT/lxprofile" "$link" 2>/dev/null; then
        [[ $name == lx && -n $verbose ]] && printf 'Создана короткая команда: lx (= lxprofile) и слитные формы lxu/lxs/lxc…\n' >&2
      fi
    fi
  done
  : >"$marker" 2>/dev/null
}

# Достаёт номер версии из текста lib/version.sh (с stdin)
_extract_version() {
  grep -m1 '^LXPROFILE_VERSION=' | sed -E 's/.*"([^"]+)".*/\1/'
}

# _version_gt A B → истина, если версия A строго новее B (semver через sort -V)
_version_gt() {
  [[ "$1" != "$2" ]] && [[ "$(printf '%s\n%s\n' "$1" "$2" | sort -V | tail -n1)" == "$1" ]]
}

# Показывает changelog из локального CHANGELOG.md:
#   $1 пуст      → последние 5 версий;
#   $1 = версия  → запись этой версии (или ошибка со списком доступных).
do_changes() {
  local want="${1:-}" file="$LXPROFILE_ROOT/CHANGELOG.md"
  if [[ ! -r $file ]]; then
    printf 'CHANGELOG.md не найден в %s\n' "$LXPROFILE_ROOT" >&2
    return 1
  fi
  if [[ -n $want ]]; then
    want=${want#v}
    local out
    out=$(awk -v v="$want" '
      $1=="##" && $2==v { p=1; print; next }
      $1=="##" && p     { p=0 }
      p { print }
    ' "$file")
    if [[ -z ${out//[[:space:]]/} ]]; then
      printf 'Версия %s не найдена. Доступные версии:\n' "$want" >&2
      grep -oE '^## [0-9]+\.[0-9]+\.[0-9]+' "$file" | sed 's/^## /  /' >&2
      return 1
    fi
    printf '%s\n' "$out"
  else
    awk '
      /^## [0-9]+\.[0-9]+\.[0-9]+/ { n++; if (n > 5) exit }
      n >= 1 { print }
    ' "$file"
  fi
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
  # Перепроверяем короткую команду lx (вдруг освободилась/установка переехала)
  ensure_lx force
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

# Печатает записи changelog для всех версий новее $1 (вашей) из CHANGELOG.md
# в репозитории. Origin уже обновлён предыдущим fetch'ем.
_changelog_since() {
  local current=$1 branch text line ver inblock=0 any=0
  branch=$(git -C "$LXPROFILE_ROOT" rev-parse --abbrev-ref HEAD 2>/dev/null)
  [[ -z $branch || $branch == HEAD ]] && branch=main
  text=$(git -C "$LXPROFILE_ROOT" show "origin/$branch:CHANGELOG.md" 2>/dev/null) || return 0
  [[ -z $text ]] && return 0
  while IFS= read -r line; do
    if [[ $line =~ ^##[[:space:]]+v?([0-9]+\.[0-9]+\.[0-9]+) ]]; then
      ver=${BASH_REMATCH[1]}
      if _version_gt "$ver" "$current"; then
        inblock=1
        (( any == 0 )) && printf 'Что нового:\n'
        any=1
        printf '  %s:\n' "$ver"
      else
        inblock=0
      fi
    elif (( inblock )) && [[ -n ${line//[[:space:]]/} ]]; then
      printf '    %s\n' "$line"
    fi
  done <<< "$text"
}

# Предлагает обновиться до версии $1.
# y → обновить и перезапуститься; n → больше не предлагать до ручного обновления.
_offer_update() {
  local remote_ver=$1
  printf 'Доступна новая версия: \033[31m%s\033[0m -> \033[32m%s\033[0m\n' "$LXPROFILE_VERSION" "$remote_ver"
  _changelog_since "$LXPROFILE_VERSION"
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

# Файл с результатом фоновой проверки и pid фонового процесса
LXPROFILE_UPD_RESULT=""
LXPROFILE_UPD_PID=""

# Запускает проверку обновлений В ФОНЕ, чтобы не тормозить открытие утилиты.
# Результат (номер новой версии) кладётся в файл, который читает finish_update_check.
start_update_check() {
  # Только в интерактивном терминале и только из git-установки
  [[ -t 0 && -t 1 ]] || return 0
  command -v git >/dev/null 2>&1 || return 0
  git -C "$LXPROFILE_ROOT" rev-parse --git-dir >/dev/null 2>&1 || return 0
  # Пользователь ранее отказался — молчим до ручного обновления
  [[ -e $LXPROFILE_DECLINED_FILE ]] && return 0

  LXPROFILE_UPD_RESULT=$(mktemp 2>/dev/null) || { LXPROFILE_UPD_RESULT=""; return 0; }
  {
    rv=$(_remote_version) || exit 0
    [[ -z $rv ]] && exit 0
    if _version_gt "$rv" "$LXPROFILE_VERSION"; then
      printf '%s\n' "$rv" > "$LXPROFILE_UPD_RESULT"
    fi
  } >/dev/null 2>&1 &
  LXPROFILE_UPD_PID=$!
}

# Дожидается фоновой проверки и, если вышла новая версия, предлагает обновиться.
# Вызывается ПОСЛЕ закрытия интерактивного просмотра.
finish_update_check() {
  [[ -n $LXPROFILE_UPD_PID ]] || return 0
  wait "$LXPROFILE_UPD_PID" 2>/dev/null
  local remote_ver=""
  [[ -r $LXPROFILE_UPD_RESULT ]] && remote_ver=$(<"$LXPROFILE_UPD_RESULT")
  rm -f "$LXPROFILE_UPD_RESULT" 2>/dev/null
  [[ -z $remote_ver ]] && return 0
  _offer_update "$remote_ver"
}
