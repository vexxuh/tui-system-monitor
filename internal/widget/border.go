package widget

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// borderColor returns the lipgloss style for coloring border elements.
func borderColor(focused bool) lipgloss.Style {
	if focused {
		return lipgloss.NewStyle().Foreground(lipgloss.Color("12")) // bright blue
	}
	return lipgloss.NewStyle().Foreground(lipgloss.Color("8")) // gray
}

// titleStyle returns the lipgloss style for the title text.
func titleStyle(focused bool) lipgloss.Style {
	s := lipgloss.NewStyle().Bold(true)
	if focused {
		s = s.Foreground(lipgloss.Color("12"))
	} else {
		s = s.Foreground(lipgloss.Color("7"))
	}
	return s
}

// TitledBorder wraps content in a hand-built bordered box with a title.
// This avoids post-processing lipgloss output which contains ANSI codes.
func TitledBorder(title, content string, width, height int, focused bool) string {
	if width < 4 {
		width = 4
	}
	if height < 3 {
		height = 3
	}

	bc := borderColor(focused)
	ts := titleStyle(focused)
	innerW := width - 2 // subtract left+right border columns

	// Build top line: ╭─ Title ─────╮
	topLeft := bc.Render("╭")
	topRight := bc.Render("╮")
	var topMiddle string
	if title != "" {
		titleStr := ts.Render(" " + title + " ")
		titleVisualLen := len(" "+title+" ") // ASCII title, no wide chars
		remaining := innerW - titleVisualLen
		if remaining < 0 {
			remaining = 0
		}
		leftDash := bc.Render("─")
		rightDashes := bc.Render(strings.Repeat("─", remaining))
		topMiddle = leftDash + titleStr + rightDashes
	} else {
		topMiddle = bc.Render(strings.Repeat("─", innerW))
	}
	topLine := topLeft + topMiddle + topRight

	// Build bottom line: ╰──────────────╯
	botLeft := bc.Render("╰")
	botRight := bc.Render("╯")
	botMiddle := bc.Render(strings.Repeat("─", innerW))
	botLine := botLeft + botMiddle + botRight

	// Build content lines, pad to innerW
	side := bc.Render("│")
	contentLines := strings.Split(content, "\n")
	innerH := height - 2 // subtract top+bottom border rows

	var rows []string
	rows = append(rows, topLine)
	for i := range innerH {
		var line string
		if i < len(contentLines) {
			line = contentLines[i]
		}
		padded := padToWidth(line, innerW)
		rows = append(rows, side+padded+side)
	}
	rows = append(rows, botLine)

	return strings.Join(rows, "\n")
}

// padToWidth pads a string with spaces to the given visual width.
// This is a simple byte-length pad that works for ASCII content.
func padToWidth(s string, width int) string {
	// Use lipgloss width measurement for accurate visual width
	vis := lipgloss.Width(s)
	if vis >= width {
		return s
	}
	return s + strings.Repeat(" ", width-vis)
}
