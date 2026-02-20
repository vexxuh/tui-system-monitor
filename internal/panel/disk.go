package panel

import (
	"fmt"
	"strings"

	"github.com/vexxuh/cy-monitor/internal/format"
	"github.com/vexxuh/cy-monitor/internal/widget"
)

// RenderDisk renders the Disk panel content.
func RenderDisk(state PanelState, width int) string {
	disk := state.Data.Disk
	contentW := width - 4
	if contentW < 10 {
		contentW = 10
	}

	barW := contentW - 16
	if barW < 6 {
		barW = 6
	}

	var lines []string

	// Device name
	device := disk.Device
	if device == "" {
		device = "unknown"
	}
	lines = append(lines, BoldStyle.Render(device))

	// Mount points
	for _, m := range disk.Mounts {
		mount := LabelStyle.Render(format.PadRight(m.Mount, 8))
		bar := usageColor(m.Percent).Render(widget.ProgressBar(m.Percent, barW-8))
		pct := ValueStyle.Render(format.Percent(m.Percent, true))
		lines = append(lines, fmt.Sprintf("%s%s %s", mount, bar, pct))
	}

	// I/O speeds
	lines = append(lines, fmt.Sprintf("%s %s  %s %s",
		LabelStyle.Render("R:"),
		GreenStyle.Render(format.BytesSpeed(disk.IORead, true)),
		LabelStyle.Render("W:"),
		YellowStyle.Render(format.BytesSpeed(disk.IOWrite, true))))

	return strings.Join(lines, "\n")
}
