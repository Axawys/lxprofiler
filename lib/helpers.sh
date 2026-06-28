#!/usr/bin/env bash
# lib/helpers.sh — счётчики и вспомогательные функции (требует lib/data.sh)

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
