package widget

import "strings"

// ProgressBar renders a text progress bar like "[####----]".
// Percent is clamped to [0, 100].
func ProgressBar(percent float64, width int) string {
	if percent < 0 {
		percent = 0
	}
	if percent > 100 {
		percent = 100
	}
	if width < 2 {
		width = 2
	}

	inner := width - 2 // subtract brackets
	filled := int(percent / 100 * float64(inner))
	if filled > inner {
		filled = inner
	}

	return "[" + strings.Repeat("#", filled) + strings.Repeat("-", inner-filled) + "]"
}
