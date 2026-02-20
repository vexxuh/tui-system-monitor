package panel

import (
	"fmt"
	"strings"

	"github.com/vexxuh/cy-monitor/internal/format"
	"github.com/vexxuh/cy-monitor/internal/sparkline"
	"github.com/vexxuh/cy-monitor/internal/widget"
)

// RenderCPU renders the CPU panel content.
func RenderCPU(state PanelState, width int) string {
	cpu := state.Data.CPU
	contentW := width - 4
	if contentW < 10 {
		contentW = 10
	}

	var lines []string

	// Sparkline
	hist := state.History.CPU.GetAll()
	spark := CyanStyle.Render(sparkline.Render(hist, contentW))
	lines = append(lines, spark)

	// Usage + frequency
	usageStr := usageColor(cpu.Usage).Render(format.Percent(cpu.Usage, true))
	freqStr := ValueStyle.Render(format.Freq(cpu.Freq, cpu.Freq > 0))
	lines = append(lines, usageStr+DimStyle.Render(" @ ")+freqStr)

	// Per-core bars, 2 per row
	barW := (contentW - 2) / 2
	if barW < 8 {
		barW = 8
	}
	for i := 0; i < len(cpu.Cores); i += 2 {
		left := DimStyle.Render(fmt.Sprintf("C%d", i)) + usageColor(cpu.Cores[i]).Render(widget.ProgressBar(cpu.Cores[i], barW-3))
		row := left
		if i+1 < len(cpu.Cores) {
			right := " " + DimStyle.Render(fmt.Sprintf("C%d", i+1)) + usageColor(cpu.Cores[i+1]).Render(widget.ProgressBar(cpu.Cores[i+1], barW-3))
			row += right
		}
		lines = append(lines, row)
	}

	// Load averages
	lines = append(lines, LabelStyle.Render("Load: ")+ValueStyle.Render(
		fmt.Sprintf("%.2f  %.2f  %.2f", cpu.Load[0], cpu.Load[1], cpu.Load[2])))

	return strings.Join(lines, "\n")
}
