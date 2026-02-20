package format

import "fmt"

// Freq formats a frequency in MHz as "X.X GHz" or "X MHz".
// Returns "--" for invalid input.
func Freq(mhz float64, valid bool) string {
	if !valid {
		return "--"
	}
	if mhz >= 1000 {
		return fmt.Sprintf("%.1f GHz", mhz/1000)
	}
	return fmt.Sprintf("%.0f MHz", mhz)
}
