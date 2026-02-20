package tui

import (
	"runtime"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/vexxuh/cy-monitor/internal/bridge"
	"github.com/vexxuh/cy-monitor/internal/model"
	"github.com/vexxuh/cy-monitor/internal/panel"
)

// Model is the Bubble Tea model for the TUI.
type Model struct {
	bridge   *bridge.Bridge
	config   model.AppConfig
	data     model.AllData
	history  panel.HistoryState
	focused  int
	expanded int
	width    int
	height   int
	platform string
	lastErr  error
}

// New creates a new TUI Model.
func New(b *bridge.Bridge, cfg model.AppConfig) Model {
	platform := runtime.GOOS
	return Model{
		bridge:   b,
		config:   cfg,
		history:  panel.NewHistoryState(cfg.HistorySize),
		focused:  1,
		expanded: 0,
		platform: platform,
	}
}

// Init implements tea.Model.
func (m Model) Init() tea.Cmd {
	return m.collectCmd()
}

// Update implements tea.Model.
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		return handleKey(m, msg)

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil

	case tickMsg:
		m.data = msg.data
		m.history = m.history.PushData(msg.data, m.config.HistorySize)
		m.lastErr = nil
		return m, m.scheduleCollect()

	case errMsg:
		m.lastErr = msg.err
		return m, m.scheduleCollect()
	}

	return m, nil
}

// View implements tea.Model.
func (m Model) View() string {
	if m.width == 0 || m.height == 0 {
		return "Initializing..."
	}

	grid := renderGrid(m)
	bar := renderStatusBar(m)

	return grid + "\n" + bar
}

// collectCmd returns a command that collects system data via the Lua bridge.
func (m Model) collectCmd() tea.Cmd {
	return func() tea.Msg {
		data, err := m.bridge.CollectAll()
		if err != nil {
			return errMsg{err: err}
		}
		return tickMsg{data: data}
	}
}

// scheduleCollect waits for the refresh interval, then fires a collection.
func (m Model) scheduleCollect() tea.Cmd {
	interval := time.Duration(m.config.RefreshInterval) * time.Second
	return tea.Tick(interval, func(_ time.Time) tea.Msg {
		data, err := m.bridge.CollectAll()
		if err != nil {
			return errMsg{err: err}
		}
		return tickMsg{data: data}
	})
}
