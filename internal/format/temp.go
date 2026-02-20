package format

import "fmt"

// Temp formats a Celsius temperature, converting to Fahrenheit if unit is "F".
// Returns "--" for invalid input.
func Temp(celsius float64, valid bool, unit string) string {
	if !valid {
		return "--"
	}
	if unit == "F" {
		f := celsius*9.0/5.0 + 32
		return fmt.Sprintf("%.0f°F", f)
	}
	return fmt.Sprintf("%.0f°C", celsius)
}
