package panel

import (
	"fmt"
	"strings"

	"github.com/vexxuh/cy-monitor/internal/format"
	"github.com/vexxuh/cy-monitor/internal/widget"
)

// RenderSystem renders the System/Battery panel content.
func RenderSystem(state PanelState, width int) string {
	bat := state.Data.Battery
	sys := state.Data.System
	contentW := width - 4
	if contentW < 10 {
		contentW = 10
	}

	barW := contentW - 14
	if barW < 6 {
		barW = 6
	}

	var lines []string

	// Battery section
	if bat.Percent >= 0 {
		barColor := usageColor(100 - bat.Percent) // invert: low battery = red
		statusColor := GreenStyle
		if bat.Status == "Discharging" {
			statusColor = YellowStyle
		}
		lines = append(lines, fmt.Sprintf("%s %s  %s",
			barColor.Render(widget.ProgressBar(bat.Percent, barW)),
			ValueStyle.Render(format.Percent(bat.Percent, true)),
			statusColor.Render(bat.Status)))

		var details []string
		if bat.Voltage > 0 {
			details = append(details, ValueStyle.Render(fmt.Sprintf("%.1fV", bat.Voltage)))
		}
		if bat.Power > 0 {
			details = append(details, ValueStyle.Render(fmt.Sprintf("%.1fW", bat.Power)))
		}
		if bat.Technology != "" {
			details = append(details, DimStyle.Render(bat.Technology))
		}
		if len(details) > 0 {
			lines = append(lines, strings.Join(details, "  "))
		}
	} else {
		lines = append(lines, DimStyle.Render("AC Power (no battery)"))
	}

	// Uptime
	lines = append(lines, fmt.Sprintf("%s %s",
		LabelStyle.Render("Uptime:"),
		ValueStyle.Render(format.Duration(sys.Uptime, sys.Uptime > 0))))

	// Procs + Load
	lines = append(lines, fmt.Sprintf("%s %s  %s %s",
		LabelStyle.Render("Procs:"),
		ValueStyle.Render(fmt.Sprintf("%d", sys.Procs)),
		LabelStyle.Render("Load:"),
		ValueStyle.Render(fmt.Sprintf("%.2f", sys.Load[0]))))

	return strings.Join(lines, "\n")
}
