#!/usr/bin/env bash
# lib/stats.sh — статистика истории команд

# ── Режим «Забавная статистика» ────────────────────────────────
STATS_OK=0
TOTAL_CMDS=0; UNIQ_CMDS=0; TOP_CMD=""; TOP_CNT=0
FF_COUNT=0; SUDO_COUNT=0; UPD_COUNT=0; RMRF_COUNT=0; TYPO_COUNT=0
VIM_COUNT=0; NVIM_COUNT=0; NANO_COUNT=0; EMACS_COUNT=0; MICRO_COUNT=0; SPAN_SEC=0

human_interval() {   # секунды → «раз в …»
  local s=$1
  if   (( s <= 0 ));     then printf '—'
  elif (( s >= 86400 )); then printf 'раз в %d дн.'  $(( s / 86400 ))
  elif (( s >= 3600 ));  then printf 'раз в %d ч.'   $(( s / 3600 ))
  elif (( s >= 60 ));    then printf 'раз в %d мин.' $(( s / 60 ))
  else                        printf 'раз в %d сек.' "$s"
  fi
}
freq() {   # count → интервал по охвату истории
  local c=$1
  (( c <= 0 || SPAN_SEC <= 0 )) && { printf '—'; return; }
  human_interval $(( SPAN_SEC / c ))
}

compute_stats() {
  local raw ts tsmin tsmax topline
  raw=$(
    { [[ -r ~/.bash_history ]] && grep -av '^#' ~/.bash_history
      [[ -r ~/.zsh_history  ]] && sed -E 's/^: [0-9]+:[0-9]+;//' ~/.zsh_history
      [[ -r ~/.local/share/fish/fish_history ]] && sed -nE 's/^- cmd: //p' ~/.local/share/fish/fish_history
    } 2>/dev/null
  )
  # первое слово команды, sudo/doas разворачиваем до реальной команды
  CMD_LIST=$(awk '{c=$1; if(c=="sudo"||c=="doas")c=$2; print c}' <<< "$raw" | grep -vE '^[[:space:]]*$')
  TOTAL_CMDS=$(grep -c . <<< "$CMD_LIST"); TOTAL_CMDS=${TOTAL_CMDS:-0}
  (( TOTAL_CMDS == 0 )) && { STATS_OK=0; return; }
  STATS_OK=1

  UNIQ_CMDS=$(sort -u <<< "$CMD_LIST" | grep -c .)
  topline=$(sort <<< "$CMD_LIST" | uniq -c | sort -rn | head -1)
  TOP_CNT=$(awk '{print $1}' <<< "$topline")
  TOP_CMD=$(awk '{print $2}' <<< "$topline")

  FF_COUNT=$(grep -cxE "${FF_MATCH_RE:-fastfetch|neofetch|screenfetch|pfetch|hyfetch}" <<< "$CMD_LIST")
  SUDO_COUNT=$(awk '{print $1}' <<< "$raw" | grep -cxE 'sudo|doas')
  UPD_COUNT=$(grep -ciE '(pacman[^|]*-S(yu|yyu|u)|(^|[; ])(yay|paru)([[:space:]]|$)|apt(-get)?[[:space:]]+(update|upgrade|full-upgrade)|dnf[[:space:]]+(update|upgrade)|zypper[[:space:]]+(up|dup|patch)|nixos-rebuild|emerge[^|]*(-u|world)|flatpak[[:space:]]+update|rpm-ostree[[:space:]]+upgrade|xbps-install[^|]*-Su)' <<< "$raw")
  RMRF_COUNT=$(grep -ciE 'rm[[:space:]]+-[a-zA-Z]*[rf][a-zA-Z]*' <<< "$raw")
  TYPO_COUNT=$(grep -cxE 'sl|gti|claer|grpe|exti|pythno|sudp|suod|cd\.\.\.?' <<< "$CMD_LIST")
  VIM_COUNT=$(grep -cxE 'vim|vi' <<< "$CMD_LIST")
  NVIM_COUNT=$(grep -cxE 'nvim' <<< "$CMD_LIST")
  NANO_COUNT=$(grep -cxE 'nano' <<< "$CMD_LIST")
  EMACS_COUNT=$(grep -cxE 'emacs|emacsclient' <<< "$CMD_LIST")
  MICRO_COUNT=$(grep -cxE 'micro' <<< "$CMD_LIST")

  ts=$(
    { [[ -r ~/.zsh_history  ]] && sed -nE 's/^: ([0-9]+):.*/\1/p' ~/.zsh_history
      [[ -r ~/.bash_history ]] && grep -aoE '^#[0-9]{9,}' ~/.bash_history | tr -d '#'
      [[ -r ~/.local/share/fish/fish_history ]] && sed -nE 's/^[[:space:]]*when: ([0-9]+).*/\1/p' ~/.local/share/fish/fish_history
    } 2>/dev/null | sort -n
  )
  tsmin=$(grep -m1 . <<< "$ts"); tsmax=$(tail -1 <<< "$ts")
  if [[ -n $tsmin && -n $tsmax ]] && (( tsmax > tsmin )); then
    SPAN_SEC=$(( tsmax - tsmin ))
  elif safe_gt "${INSTALL_EPOCH:-0}" 0; then
    SPAN_SEC=$(( $(date +%s) - INSTALL_EPOCH ))
  else
    SPAN_SEC=0
  fi
}

# Интервал между вызовами в секундах (-1 если не определить)
_ivl() { local c=$1; (( c <= 0 || SPAN_SEC <= 0 )) && { echo -1; return; }; echo $(( SPAN_SEC / c )); }

_ff_quip()   { (( FF_COUNT == 0 )) && { printf 'ни разу, аскет'; return; }
               local i; i=$(_ivl "$FF_COUNT")
               if   (( i < 0 ));        then printf 'частоту не определить';
               elif (( i >= 2592000 )); then printf 'очень редко';
               elif (( i >= 604800 ));  then printf 'иногда любуешься системой';
               elif (( i >= 86400 ));   then printf 'почти ежедневный ритуал';
               elif (( i >= 3600 ));    then printf 'по нескольку раз в день';
               else printf 'это уже зависимость))'; fi; }
_upd_quip()  { (( UPD_COUNT == 0 )) && { printf 'ни разу — смело'; return; }
               local i; i=$(_ivl "$UPD_COUNT")
               if   (( i < 0 ));        then printf 'частоту не определить';
               elif (( i >= 2592000 )); then printf 'обновляешься редко — стабильность важнее';
               elif (( i >= 604800 ));  then printf 'апдейт по выходным';
               elif (( i >= 86400 ));   then printf 'держишь систему свежей';
               else printf 'апдейт — это медитация'; fi; }
_rmrf_quip() { (( RMRF_COUNT == 0 )) && { printf 'аккуратно'; return; }
               local i; i=$(_ivl "$RMRF_COUNT")
               if   (( i < 0 ));        then printf 'частоту не определить';
               elif (( i >= 2592000 )); then printf 'редко, но метко';
               elif (( i >= 604800 ));  then printf 'бывает';
               elif (( i >= 86400 ));   then printf 'живёшь опасно';
               else printf 'как ты ещё жив?'; fi; }
_editor_win(){
  local names=(vim nvim nano emacs micro)
  local counts=("$VIM_COUNT" "$NVIM_COUNT" "$NANO_COUNT" "$EMACS_COUNT" "$MICRO_COUNT")
  local i total=0 bestn=-1 ties=0 best=""
  for i in "${!counts[@]}"; do total=$(( total + counts[i] )); (( counts[i] > bestn )) && bestn=${counts[i]}; done
  (( total == 0 )) && { printf 'все мимо — GUI?'; return; }
  for i in "${!counts[@]}"; do (( counts[i] == bestn )) && { ties=$(( ties + 1 )); best=${names[i]}; }; done
  (( ties > 1 )) && { printf 'ничья'; return; }
  printf 'победил %s' "$best"
}
_top_quip()  { local i; i=$(_ivl "$TOP_CNT")
               if   (( i < 0 ));      then printf 'частоту не определить';
               elif (( i >= 86400 )); then printf 'заходит нечасто';
               elif (( i >= 3600 ));  then printf 'крепкая привычка';
               elif (( i >= 600 ));   then printf 'мышечная память';
               else printf 'набита вслепую'; fi; }
_sudo_quip() { (( SUDO_COUNT == 0 )) && { printf 'живёшь без рута'; return; }
               local i; i=$(_ivl "$SUDO_COUNT")
               if   (( i < 0 ));      then printf 'частоту не определить';
               elif (( i >= 86400 )); then printf 'рут по праздникам';
               elif (( i >= 3600 ));  then printf 'уверенно у руля';
               else printf 'практически root'; fi; }

render_stats() {
  render_header "Забавная статистика"
  if (( STATS_OK == 0 )); then
    printf '%s\n' "  ${DIM}История команд пуста или недоступна.${RESET}"
    printf '%s\n' "  ${DIM}Подсказка: включи ${BOLD}HISTTIMEFORMAT${RESET}${DIM} — и время будет точнее.${RESET}"
    render_footer
    return
  fi
  local span_d=$(( SPAN_SEC / 86400 ))
  printf '%s\n' "  В истории ${BOLD}${TOTAL_CMDS}${RESET} команд, ${BOLD}${UNIQ_CMDS}${RESET} уникальных${DIM} (охват ~${span_d} дн.)${RESET}"
  printf '\n'
  printf '%s\n' "  Любимая команда: ${BOLD}${GREEN}${TOP_CMD}${RESET} — ${BOLD}${TOP_CNT}×${RESET} ${DIM}($(freq "$TOP_CNT")) — $(_top_quip)${RESET}"
  printf '%s\n' "  fastfetch/neofetch${FF_ALIAS_LABEL}: ${BOLD}${FF_COUNT}×${RESET} ${DIM}($(freq "$FF_COUNT")) — $(_ff_quip)${RESET}"
  printf '%s\n' "  Обновления: ${BOLD}${UPD_COUNT}×${RESET} ${DIM}($(freq "$UPD_COUNT")) — $(_upd_quip)${RESET}"
  printf '%s\n' "  sudo/doas: ${BOLD}${SUDO_COUNT}×${RESET} ${DIM}($(freq "$SUDO_COUNT")) — $(_sudo_quip)${RESET}"
  printf '%s\n' "  rm -rf: ${BOLD}${RMRF_COUNT}×${RESET} ${DIM}— $(_rmrf_quip)${RESET}"
  printf '%s\n' "  Опечаток поймано: ${BOLD}${TYPO_COUNT}${RESET}${DIM} (sl, gti, claer, cd..…)${RESET}"
  printf '%s\n' "  Редактор-война: ${DIM}vim${RESET} ${VIM_COUNT} : ${DIM}nvim${RESET} ${NVIM_COUNT} : ${DIM}nano${RESET} ${NANO_COUNT} : ${DIM}emacs${RESET} ${EMACS_COUNT} : ${DIM}micro${RESET} ${MICRO_COUNT}  ${DIM}→ $(_editor_win)${RESET}"
  render_footer
}
