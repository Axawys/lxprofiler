// Package lxprofiler хранит встроенные в бинарник ресурсы (корень модуля).
package lxprofiler

import _ "embed"

// ChangelogMD — содержимое CHANGELOG.md, вшитое в бинарник на этапе сборки.
// Показывается командой `lxprofile -c` и в уведомлении об обновлении.
//
//go:embed CHANGELOG.md
var ChangelogMD string
