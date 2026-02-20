package panel

import "github.com/charmbracelet/lipgloss"

// Color palette for panel rendering.
var (
	LabelStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("4"))  // blue
	ValueStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("15")) // bright white
	DimStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("8"))  // gray
	GreenStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("2"))  // green
	RedStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("1"))  // red
	YellowStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("3")) // yellow
	CyanStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("6"))  // cyan
	BoldStyle  = lipgloss.NewStyle().Bold(true)
)

// tempColor returns a style based on temperature thresholds.
func tempColor(celsius float64) lipgloss.Style {
	switch {
	case celsius >= 80:
		return RedStyle
	case celsius >= 60:
		return YellowStyle
	default:
		return GreenStyle
	}
}

// usageColor returns a style based on usage percentage.
func usageColor(pct float64) lipgloss.Style {
	switch {
	case pct >= 90:
		return RedStyle
	case pct >= 70:
		return YellowStyle
	default:
		return GreenStyle
	}
}
