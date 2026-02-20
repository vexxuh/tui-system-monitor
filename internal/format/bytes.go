package format

import "fmt"

// Bytes formats a byte count as a human-readable string (B/K/M/G).
// Negative values are clamped to 0. Returns "--" for invalid input.
func Bytes(value float64, valid bool) string {
	if !valid {
		return "--"
	}
	if value < 0 {
		value = 0
	}

	switch {
	case value >= 1e9:
		return fmt.Sprintf("%.1fG", value/1e9)
	case value >= 1e6:
		return fmt.Sprintf("%.1fM", value/1e6)
	case value >= 1e3:
		return fmt.Sprintf("%.1fK", value/1e3)
	default:
		return fmt.Sprintf("%.0fB", value)
	}
}

// BytesSpeed formats a byte-per-second value as "X/s".
func BytesSpeed(value float64, valid bool) string {
	if !valid {
		return "-- /s"
	}
	return Bytes(value, true) + "/s"
}
