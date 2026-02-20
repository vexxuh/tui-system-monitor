package panel

import (
	"fmt"
	"strings"

	"github.com/vexxuh/cy-monitor/internal/format"
	"github.com/vexxuh/cy-monitor/internal/sparkline"
)

// RenderNetwork renders the Network panel content.
func RenderNetwork(state PanelState, width int) string {
	net := state.Data.Network
	contentW := width - 4
	if contentW < 10 {
		contentW = 10
	}

	sparkW := contentW - 18
	if sparkW < 3 {
		sparkW = 3
	}

	var lines []string

	// Interface name
	lines = append(lines, BoldStyle.Render(net.Interface))

	// Download speed + sparkline
	dlSpeed := format.PadRight(format.BytesSpeed(net.RxSpeed, true), 14)
	rxSpark := GreenStyle.Render(sparkline.Render(state.History.NetRx.GetAll(), sparkW))
	lines = append(lines, GreenStyle.Render("↓ ")+ValueStyle.Render(dlSpeed)+rxSpark)

	// Upload speed + sparkline
	ulSpeed := format.PadRight(format.BytesSpeed(net.TxSpeed, true), 14)
	txSpark := CyanStyle.Render(sparkline.Render(state.History.NetTx.GetAll(), sparkW))
	lines = append(lines, CyanStyle.Render("↑ ")+ValueStyle.Render(ulSpeed)+txSpark)

	// Totals
	lines = append(lines, fmt.Sprintf("%s %s %s  %s %s %s",
		LabelStyle.Render("Total:"),
		GreenStyle.Render("↓"),
		ValueStyle.Render(format.Bytes(float64(net.RxTotal), true)),
		CyanStyle.Render("↑"),
		ValueStyle.Render(format.Bytes(float64(net.TxTotal), true)),
		""))

	return strings.Join(lines, "\n")
}
