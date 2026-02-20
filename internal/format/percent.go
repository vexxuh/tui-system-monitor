package format

import "fmt"

// Percent formats a value as "X%". Returns "--%" for invalid input.
func Percent(value float64, valid bool) string {
	if !valid {
		return "--%"
	}
	return fmt.Sprintf("%.0f%%", value)
}
