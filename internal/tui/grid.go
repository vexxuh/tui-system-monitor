package tui

import (
	"github.com/charmbracelet/lipgloss"
	"github.com/vexxuh/cy-monitor/internal/panel"
	"github.com/vexxuh/cy-monitor/internal/widget"
)

// renderGrid produces the 2x3 panel grid or a single expanded panel.
func renderGrid(m Model) string {
	state := panel.PanelState{
		Data:     m.data,
		Config:   m.config,
		History:  m.history,
		Focused:  m.focused,
		Expanded: m.expanded,
	}

	if m.expanded >= 1 && m.expanded <= 6 {
		return renderExpanded(state, m.width, m.height-1) // -1 for status bar
	}

	return renderNormal(state, m.width, m.height-1)
}

func renderExpanded(state panel.PanelState, width, height int) string {
	content := renderPanel(state.Expanded, state, width)
	return widget.TitledBorder(
		panel.PanelName(state.Expanded),
		content,
		width, height,
		true,
	)
}

func renderNormal(state panel.PanelState, width, height int) string {
	cols := 3
	rows := 2
	cellW := width / cols
	cellH := height / rows

	panelIdx := 1
	var rowViews []string

	for r := range rows {
		var colViews []string
		for c := range cols {
			w := cellW
			h := cellH
			// Last column gets remainder width
			if c == cols-1 {
				w = width - cellW*(cols-1)
			}
			// Last row gets remainder height
			if r == rows-1 {
				h = height - cellH*(rows-1)
			}

			content := renderPanel(panelIdx, state, w)
			bordered := widget.TitledBorder(
				panel.PanelName(panelIdx),
				content,
				w, h,
				state.Focused == panelIdx,
			)
			colViews = append(colViews, bordered)
			panelIdx++
		}
		rowViews = append(rowViews, lipgloss.JoinHorizontal(lipgloss.Top, colViews...))
	}

	return lipgloss.JoinVertical(lipgloss.Left, rowViews...)
}

func renderPanel(idx int, state panel.PanelState, width int) string {
	switch idx {
	case 1:
		return panel.RenderCPU(state, width)
	case 2:
		return panel.RenderMemory(state, width)
	case 3:
		return panel.RenderTemps(state, width)
	case 4:
		return panel.RenderNetwork(state, width)
	case 5:
		return panel.RenderDisk(state, width)
	case 6:
		return panel.RenderSystem(state, width)
	default:
		return ""
	}
}
