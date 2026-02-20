package tui

import (
	"fmt"
	"time"

	"github.com/charmbracelet/lipgloss"
)

var statusStyle = lipgloss.NewStyle().
	Background(lipgloss.Color("4")).
	Foreground(lipgloss.Color("15")).
	Padding(0, 1)

var errStyle = lipgloss.NewStyle().
	Background(lipgloss.Color("1")).
	Foreground(lipgloss.Color("15")).
	Bold(true).
	Padding(0, 1)

// renderStatusBar produces the bottom status bar.
func renderStatusBar(m Model) string {
	now := time.Now().Format("15:04:05")

	status := fmt.Sprintf(
		"%s | %s | q:quit r:refresh c/f:%s 1-6:expand Tab:focus",
		now,
		m.platform,
		m.config.TempUnit,
	)

	bar := statusStyle.Width(m.width).Render(status)

	if m.lastErr != nil {
		errMsg := errStyle.Render(fmt.Sprintf(" ERR: %v ", m.lastErr))
		bar = statusStyle.Width(m.width - lipgloss.Width(errMsg)).Render(status) + errMsg
	}

	return bar
}
