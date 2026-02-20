package format

import "unicode/utf8"

// PadRight pads or truncates a string to exactly width characters.
func PadRight(s string, width int) string {
	n := utf8.RuneCountInString(s)
	if n >= width {
		runes := []rune(s)
		return string(runes[:width])
	}
	result := make([]byte, 0, len(s)+width-n)
	result = append(result, s...)
	for range width - n {
		result = append(result, ' ')
	}
	return string(result)
}

// Truncate truncates a string and appends "…" if it exceeds maxWidth.
func Truncate(s string, maxWidth int) string {
	if maxWidth <= 0 {
		return ""
	}
	runes := []rune(s)
	if len(runes) <= maxWidth {
		return s
	}
	if maxWidth <= 1 {
		return "…"
	}
	return string(runes[:maxWidth-1]) + "…"
}
