package panel

import (
	"fmt"
	"strings"

	"github.com/vexxuh/cy-monitor/internal/format"
	"github.com/vexxuh/cy-monitor/internal/widget"
)

// RenderMemory renders the Memory panel content.
func RenderMemory(state PanelState, width int) string {
	mem := state.Data.Memory
	contentW := width - 4
	if contentW < 10 {
		contentW = 10
	}

	barW := contentW - 8
	if barW < 6 {
		barW = 6
	}

	var lines []string

	// RAM bar
	ramPct := 0.0
	if mem.RAMTotal > 0 {
		ramPct = float64(mem.RAMUsed) / float64(mem.RAMTotal) * 100
	}
	lines = append(lines, fmt.Sprintf("%s %s %s",
		LabelStyle.Render("RAM"),
		usageColor(ramPct).Render(widget.ProgressBar(ramPct, barW)),
		ValueStyle.Render(format.Percent(ramPct, true))))

	// RAM used/total
	lines = append(lines, fmt.Sprintf("%s %s %s",
		ValueStyle.Render(format.Bytes(float64(mem.RAMUsed), true)),
		DimStyle.Render("/"),
		DimStyle.Render(format.Bytes(float64(mem.RAMTotal), true))))

	// Swap bar
	swapPct := 0.0
	if mem.SwapTotal > 0 {
		swapPct = float64(mem.SwapUsed) / float64(mem.SwapTotal) * 100
	}
	lines = append(lines, fmt.Sprintf("%s %s %s",
		LabelStyle.Render("Swap"),
		usageColor(swapPct).Render(widget.ProgressBar(swapPct, barW-1)),
		ValueStyle.Render(format.Percent(swapPct, true))))

	// Swap used/total
	lines = append(lines, fmt.Sprintf("%s %s %s",
		ValueStyle.Render(format.Bytes(float64(mem.SwapUsed), true)),
		DimStyle.Render("/"),
		DimStyle.Render(format.Bytes(float64(mem.SwapTotal), true))))

	// Cached + Buffers
	lines = append(lines, fmt.Sprintf("%s %s  %s %s",
		LabelStyle.Render("Cached:"),
		ValueStyle.Render(format.Bytes(float64(mem.Cached), true)),
		LabelStyle.Render("Buf:"),
		ValueStyle.Render(format.Bytes(float64(mem.Buffers), true))))

	return strings.Join(lines, "\n")
}
