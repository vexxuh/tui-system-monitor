package format

import "fmt"

// Duration formats seconds as "Xd Xh Xm".
// Returns "--" for invalid, zero, or negative values.
func Duration(secs float64, valid bool) string {
	if !valid || secs <= 0 {
		return "--"
	}

	total := int(secs)
	days := total / 86400
	hours := (total % 86400) / 3600
	minutes := (total % 3600) / 60

	switch {
	case days > 0:
		return fmt.Sprintf("%dd %dh %dm", days, hours, minutes)
	case hours > 0:
		return fmt.Sprintf("%dh %dm", hours, minutes)
	default:
		return fmt.Sprintf("%dm", minutes)
	}
}
