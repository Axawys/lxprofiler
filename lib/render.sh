#!/usr/bin/env bash
# lib/render.sh — оформление, статический и интерактивный вывод

# ──────────────────────────────────────────────
# Оформление
# ──────────────────────────────────────────────

BOLD=$'\033[1m'; DIM=$'\033[2m'; RESET=$'\033[0m'
GREEN=$'\033[32m'; YELLOW=$'\033[33m'; CYAN=$'\033[36m'

# Значок пингвина показываем только в UTF-8 локали — иначе вместо 🐧 будут
# «кракозябры». В не-UTF-8 окружении заголовки идут без эмодзи.
if [[ "${LANG:-}${LC_ALL:-}${LC_CTYPE:-}" == *[Uu][Tt][Ff]* ]]; then
  PENGUIN="🐧 "
else
  PENGUIN=""
fi

# Ширина полоски заполненности. По умолчанию 20 (как и раньше); на узких
# терминалах (Termux) interactive_view уменьшает её, не трогая прочие системы.
BAR_W=20

make_bar() {
  local p=$1 w=${BAR_W:-20} filled=$(( p * ${BAR_W:-20} / 100 )) i=0 b=""
  while (( i < filled )); do b+="█"; i=$(( i + 1 )); done
  while (( i < w ));      do b+="░"; i=$(( i + 1 )); done
  printf '%s' "$b"
}

# «Сломанная» полоска для таинственных классов (ширины BAR_W)
make_broken_bar() {
  local w=${BAR_W:-20}
  if (( w == 20 )); then printf '%s' "█▒ █░▓  ░▒█ ░ ▓█░ ▓?"; return; fi
  local pat="█▒ █░▓ ░▒█ ░ ▓█░ ▓" s=""
  while (( ${#s} < w )); do s+=$pat; done
  s=${s:0:w}
  (( w >= 1 )) && s="${s:0:w-1}?"   # последний символ — вопрос, для «сломанности»
  printf '%s' "$s"
}

# Маскирует название вопросительными знаками (той же длины — раскладка не плывёт)
mask_label() {
  local n=${#1} s=""
  while (( n-- > 0 )); do s+="?"; done
  printf '%s' "$s"
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

  local p color pad pf bar
  for key in "${sorted_keys[@]}"; do
    p=${norm_score[$key]}; label="${LABEL[$key]}"
    if [[ -n "${MYSTERY[$key]:-}" ]]; then
      # в статике навести курсор нельзя — имя остаётся закрытым
      color=$DIM; pf="???"; bar=$(make_broken_bar); label=$(mask_label "$label")
    else
      if safe_ge "$p" 80; then color=$GREEN
      elif safe_ge "$p" 50; then color=$YELLOW
      else color=$DIM; fi
      pf=$(printf '%3d' "$p"); bar=$(make_bar "$p")
    fi
    pad=$(( maxlen - ${#label} ))
    printf "${color}%s%*s${RESET}  %s%%  ${color}%s${RESET}\n" "$label" "$pad" "" "$pf" "$bar"
  done

  local W=${sorted_keys[0]} S=${sorted_keys[1]} T=${sorted_keys[2]}
  echo
  echo "${BOLD}${GREEN}▶ Ты — ${LABEL[$W]}${RESET}"
  echo "  $(describe "$W" "${norm_score[$W]}")"
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
  # Чистим терминал после выхода из интерактивного просмотра
  tput clear 2>/dev/null || clear 2>/dev/null
}

# ── Общие части кадра ──────────────────────────────────────────
render_header() {
  printf '%s\n\n' "${BOLD}${CYAN}  ${PENGUIN}$1${RESET}"
}

render_footer() {
  printf '%s\n' "${DIM}  ────────────────────────────────────────────${RESET}"
  printf '%s' "${DIM}  ${BOLD}↑↓${RESET}${DIM} — листать · ${BOLD}${CYAN}←→${RESET}${DIM} — режим [${MODE_NAME}] · ${BOLD}${YELLOW}q${RESET}${DIM} — выход${RESET}"
}

# ── Режим «Список» ─────────────────────────────────────────────
render_list() {
  local sel=$1 i sk lbl p pad used l
  render_header "Профиль архетипов"

  local pf bar
  for i in "${!sorted_keys[@]}"; do
    sk=${sorted_keys[i]}; lbl=${LABEL[$sk]}; p=${norm_score[$sk]}
    if [[ -n "${MYSTERY[$sk]:-}" ]] && (( i != sel )); then
      # пока не наведён — имя скрыто вопросами, процент «???», полоска сломана
      pf="???"; bar=$(make_broken_bar); lbl=$(mask_label "$lbl")
    else
      # обычный класс или таинственный под маркером — нормальный вид
      pf=$(printf '%3d' "$p"); bar=$(make_bar "$p")
    fi
    pad=$(( MAXLEN - ${#lbl} ))
    if (( i == sel )); then
      printf "${GREEN}${BOLD}▶ %s%*s  %s%%  %s${RESET}\n" "$lbl" "$pad" "" "$pf" "$bar"
    else
      printf "${DIM}  %s%*s  %s%%  %s${RESET}\n" "$lbl" "$pad" "" "$pf" "$bar"
    fi
  done

  printf '%s\n' "${DIM}  ────────────────────────────────────────────${RESET}"

  sk=${sorted_keys[sel]}; used=0
  # Выбранный элемент всегда «раскрыт» — показываем настоящий процент
  printf "${BOLD}${GREEN}▶ %s — %d%%${RESET}\n" "${LABEL[$sk]}" "${norm_score[$sk]}"; used=$(( used + 1 ))
  wrap_into "$(describe "$sk" "${norm_score[$sk]}")" "$WRAP_W"
  for l in "${WRAPPED[@]}"; do printf "  %s\n" "$l"; used=$(( used + 1 )); done
  printf '\n'; used=$(( used + 1 ))
  printf "${BOLD}  Что повлияло:${RESET}\n"; used=$(( used + 1 ))
  wrap_into "${reasons[$sk]:-—}" "$WRAP_W"
  for l in "${WRAPPED[@]}"; do printf "${DIM}  %s${RESET}\n" "$l"; used=$(( used + 1 )); done
  while (( used < DETAIL_H )); do printf '\n'; used=$(( used + 1 )); done

  render_footer
}

# ── Диспетчер кадра ────────────────────────────────────────────
render_frame() {
  local sel=$1
  tput cup 0 0 2>/dev/null
  tput ed   2>/dev/null
  case "$VIEW_MODE" in
    compass) render_compass "$sel" ;;
    stats)   render_stats ;;
    *)       render_list "$sel" ;;
  esac
}

# Глобальные параметры раскладки интерактивного вида
MAXLEN=0
WRAP_W=72
DETAIL_H=0
KERNEL=""
INIT1=""
VIEW_MODE="list"
MODE_NAME="список"

# Циклическое перелистывание режимов (для m и vim-style h/l)
MODES=(list compass stats)
MODE_NAMES=(список компас статистика)
MODE_IDX=0
cycle_mode() {   # $1 = +1 (вперёд) / -1 (назад)
  local n=${#MODES[@]}
  MODE_IDX=$(( (MODE_IDX + $1 + n) % n ))
  VIEW_MODE=${MODES[MODE_IDX]}
  MODE_NAME=${MODE_NAMES[MODE_IDX]}
}

interactive_view() {
  local cols sk d c total l last sel=0

  KERNEL=$(uname -r 2>/dev/null || echo '?')
  INIT1=$(ps -p 1 -o comm= 2>/dev/null || echo '?')
  compute_compass
  compute_stats

  cols=$(tput cols 2>/dev/null || echo 80)
  WRAP_W=$(( cols - 4 ))
  (( WRAP_W > 76 )) && WRAP_W=76
  (( WRAP_W < 30 )) && WRAP_W=30

  # Ширина колонки меток
  MAXLEN=0
  for sk in "${sorted_keys[@]}"; do
    (( ${#LABEL[$sk]} > MAXLEN )) && MAXLEN=${#LABEL[$sk]}
  done

  # Ширина полоски. На прочих системах — прежние 20. На Termux (узкий экран
  # телефона) ужимаем под ширину терминала, чтобы строки не переносились.
  BAR_W=20
  if (( ${META_ANDROID:-0} )); then
    BAR_W=$(( cols - MAXLEN - 12 ))   # «▶ » + отступы + «100%» + отступы
    (( BAR_W > 20 )) && BAR_W=20
    (( BAR_W < 6 ))  && BAR_W=6
  fi

  # Фиксированная высота нижней панели = максимум по всем классам,
  # чтобы раскладка не «прыгала» при смене выделения
  DETAIL_H=0
  for sk in "${sorted_keys[@]}"; do
    wrap_into "$(describe "$sk" "${norm_score[$sk]}")" "$WRAP_W"; d=${#WRAPPED[@]}
    wrap_into "${reasons[$sk]:-—}" "$WRAP_W"; c=${#WRAPPED[@]}
    total=$(( 1 + d + 1 + 1 + c ))   # заголовок + описание + пусто + «Что повлияло:» + критерии
    (( total > DETAIL_H )) && DETAIL_H=$total
  done

  last=$(( ${#sorted_keys[@]} - 1 ))

  tput smcup 2>/dev/null; tput civis 2>/dev/null
  # На Ctrl+C/SIGTERM восстанавливаем терминал и выходим аккуратно
  trap '_view_cleanup' EXIT
  trap '_view_cleanup; exit 130' INT TERM

  while true; do
    render_frame "$sel"
    read_key
    case "$REPLY_KEY" in
      q|Q)          break ;;
      j|$'\x1b[B')  (( sel < last )) && sel=$(( sel + 1 )) ;;
      k|$'\x1b[A')  (( sel > 0 ))    && sel=$(( sel - 1 )) ;;
      g)            sel=0 ;;
      G)            sel=$last ;;
      m|M|l|$'\x1b[C')  cycle_mode 1 ;;   # → / l — следующий режим
      h|$'\x1b[D')      cycle_mode -1 ;;  # ← / h — предыдущий режим
    esac
  done

  _view_cleanup
  trap - EXIT INT TERM
}
