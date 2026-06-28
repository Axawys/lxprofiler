#!/usr/bin/env bash
# lib/compass.sh — линуксоидный компас (векторы, расчёт, отрисовка)

# ── Линуксоидный компас: векторы классов ──────────────────────
#   x: − = новаторы (лево), + = традиции (право)
#   y: + = контроль/DIY (верх), − = удобство/из-коробки (низ)
declare -A VX=(
  [devops]=-20 [programmer]=-10 [sysadmin]=50  [minimalist]=0
  [old_hacker]=70 [ricer]=-30 [gamer]=-10 [anonymous]=0
  [pentester]=10 [import_substituted]=50 [fresh_witness]=-90 [atomic]=-60
)
declare -A VY=(
  [devops]=-20 [programmer]=30 [sysadmin]=40  [minimalist]=60
  [old_hacker]=70 [ricer]=70 [gamer]=-60 [anonymous]=40
  [pentester]=50 [import_substituted]=-10 [fresh_witness]=-10 [atomic]=-50
)

AGG_CX=0
AGG_CY=0
compute_compass() {
  local key sumx=0 sumy=0 sumw=0 w
  for key in "${!score[@]}"; do
    w=${score[$key]}
    (( w <= 0 )) && continue
    sumx=$(( sumx + w * ${VX[$key]} ))
    sumy=$(( sumy + w * ${VY[$key]} ))
    sumw=$(( sumw + w ))
  done
  (( sumw == 0 )) && sumw=1
  AGG_CX=$(( sumx / sumw ))
  AGG_CY=$(( sumy / sumw ))
}

compass_quadrant() {   # CX CY → название квадранта
  local cx=$1 cy=$2 h v
  if   (( cx <= -15 )); then h=N
  elif (( cx >=  15 )); then h=T
  else h=C; fi
  if   (( cy >=  15 )); then v=U
  elif (( cy <= -15 )); then v=D
  else v=C; fi
  case "$h$v" in
    NU) printf '%s' "Лаборатория — DIY-новатор (Arch/NixOS/tiling)" ;;
    TU) printf '%s' "Цитадель Unix — всё руками, старая школа" ;;
    ND) printf '%s' "Гладкое будущее — новое и из коробки (atomic/Bazzite)" ;;
    TD) printf '%s' "Тёплая гавань — стабильно и удобно (Ubuntu/Mint)" ;;
    CU) printf '%s' "Инженер-середняк — по взглядам центрист, но всё руками" ;;
    CD) printf '%s' "Прагматик — посередине по взглядам, ценит удобство" ;;
    NC) printf '%s' "Новатор-центрист — за свежее, баланс DIY и удобства" ;;
    TC) printf '%s' "Традиционалист-центрист — проверенное, баланс DIY и удобства" ;;
    *)  printf '%s' "Центрист — сбалансированный линуксоид" ;;
  esac
}

# ── Режим «Координаты» (линуксоидный компас) ───────────────────
render_compass() {
  local sel=$1 selkey=${sorted_keys[sel]}
  render_header "Линуксоидный компас"

  local GW=49 GH=15 cols rows
  cols=$(tput cols 2>/dev/null || echo 80)
  rows=$(tput lines 2>/dev/null || echo 24)
  (( GW > cols - 4 ))  && GW=$(( cols - 4 ))
  (( GW < 21 ))        && GW=21
  (( GW % 2 == 0 ))    && GW=$(( GW - 1 ))
  (( GH > rows - 14 )) && GH=$(( rows - 14 ))
  (( GH < 7 ))         && GH=7
  (( GH % 2 == 0 ))    && GH=$(( GH - 1 ))
  local ccol=$(( (GW - 1) / 2 )) crow=$(( (GH - 1) / 2 ))

  # Единственная точка — твой центр масс (общие взгляды пользователя)
  local -A CELL=()
  local ac=$(( (AGG_CX + 100) * (GW - 1) / 200 ))
  local ar=$(( (100 - AGG_CY) * (GH - 1) / 200 ))
  (( ac < 0 )) && ac=0; (( ac > GW - 1 )) && ac=$(( GW - 1 ))
  (( ar < 0 )) && ar=0; (( ar > GH - 1 )) && ar=$(( GH - 1 ))
  CELL["$ar,$ac"]="${CYAN}${BOLD}●${RESET}"

  printf '   %s\n' "${BOLD}▲ КОНТРОЛЬ (всё руками)${RESET}"
  local r c cell line
  for (( r=0; r<GH; r++ )); do
    line="   "
    for (( c=0; c<GW; c++ )); do
      cell=${CELL["$r,$c"]:-}
      if [[ -z $cell ]]; then
        if   (( r == crow && c == ccol )); then cell="${DIM}┼${RESET}"
        elif (( r == crow ));              then cell="${DIM}─${RESET}"
        elif (( c == ccol ));              then cell="${DIM}│${RESET}"
        else cell=" "; fi
      fi
      line+="$cell"
    done
    printf '%s\n' "$line"
  done
  printf '   %s\n' "${BOLD}▼ УДОБСТВО (из коробки)${RESET}"
  printf '   %s\n' "${DIM}◄ новаторы$(printf '%*s' $(( GW - 20 )) '')традиции ►${RESET}"

  # Подписи знаков координат
  local sx=$AGG_CX sy=$AGG_CY
  [[ ${sx:0:1} != - ]] && sx="+$sx"
  [[ ${sy:0:1} != - ]] && sy="+$sy"

  printf '\n'
  printf '%s\n' "  ${CYAN}${BOLD}●${RESET} ${BOLD}ты:${RESET} $(compass_quadrant "$AGG_CX" "$AGG_CY")"
  printf '%s\n' "  ${DIM}   координаты:  новат↔трад ${sx}   ·   контроль↔удоб ${sy}${RESET}"

  render_footer
}
