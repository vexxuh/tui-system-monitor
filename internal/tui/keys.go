package tui

import (
	tea "github.com/charmbracelet/bubbletea"
)

// handleKey processes a key press and returns the updated model.
func handleKey(m Model, msg tea.KeyMsg) (Model, tea.Cmd) {
	switch msg.String() {
	case "q", "ctrl+c":
		return m, tea.Quit

	case "r":
		return m, m.collectCmd()

	case "c", "f":
		if m.config.TempUnit == "C" {
			m.config.TempUnit = "F"
		} else {
			m.config.TempUnit = "C"
		}
		return m, nil

	case "tab":
		m.focused++
		if m.focused > 6 {
			m.focused = 1
		}
		return m, nil

	case "esc":
		m.expanded = 0
		return m, nil

	case "1", "2", "3", "4", "5", "6":
		idx := int(msg.String()[0] - '0')
		if m.expanded == idx {
			m.expanded = 0
		} else {
			m.expanded = idx
		}
		return m, nil
	}

	return m, nil
}
