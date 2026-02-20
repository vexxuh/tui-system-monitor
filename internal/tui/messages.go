package tui

import (
	"github.com/vexxuh/cy-monitor/internal/model"
)

// tickMsg signals that a data collection cycle completed.
type tickMsg struct {
	data model.AllData
}

// errMsg signals a collection error.
type errMsg struct {
	err error
}
