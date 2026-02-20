package panel

import (
	"fmt"
	"strings"

	"github.com/vexxuh/cy-monitor/internal/format"
	"github.com/vexxuh/cy-monitor/internal/sparkline"
)

// RenderTemps renders the Temperature panel content.
func RenderTemps(state PanelState, width int) string {
	temps := state.Data.Temps
	contentW := width - 4
	if contentW < 10 {
		contentW = 10
	}

	if len(temps) == 0 {
		return DimStyle.Render("No sensors found")
	}

	sparkW := contentW - 20
	if sparkW < 3 {
		sparkW = 3
	}

	var lines []string
	for _, sensor := range temps {
		label := LabelStyle.Render(format.PadRight(sensor.Label, 10))
		temp := tempColor(sensor.TempC).Render(format.PadRight(
			format.Temp(sensor.TempC, true, state.Config.TempUnit), 8))

		hist, ok := state.History.Temps[sensor.Label]
		spark := ""
		if ok {
			spark = CyanStyle.Render(sparkline.Render(hist.GetAll(), sparkW))
		} else {
			spark = strings.Repeat(" ", sparkW)
		}

		lines = append(lines, fmt.Sprintf("%s%s%s", label, temp, spark))
	}

	return strings.Join(lines, "\n")
}
