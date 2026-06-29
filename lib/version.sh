#!/usr/bin/env bash
# lib/version.sh — версия lxprofiler
#
# !!! НЕ ЗАБУДЬ обновить LXPROFILE_VERSION и добавить запись в CHANGELOG.md
#     перед коммитом с новой функциональностью.
#     Схема semver: MAJOR.MINOR.PATCH
#       MAJOR — несовместимые изменения, MINOR — новые фичи, PATCH — багфиксы.
#     Это число показывается в `lxprofile --version` (вместе с git-хешем),
#     а CHANGELOG.md — в уведомлении об обновлении.

LXPROFILE_VERSION="4.10.1"
