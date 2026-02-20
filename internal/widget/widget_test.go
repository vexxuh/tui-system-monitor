package widget

import (
	"strings"
	"testing"
)

func TestProgressBar(t *testing.T) {
	tests := []struct {
		percent float64
		width   int
		want    string
	}{
		{50, 12, "[#####-----]"},
		{0, 12, "[----------]"},
		{100, 12, "[##########]"},
		{-10, 12, "[----------]"},
		{110, 12, "[##########]"},
	}
	for _, tt := range tests {
		got := ProgressBar(tt.percent, tt.width)
		if got != tt.want {
			t.Errorf("ProgressBar(%v, %v) = %q, want %q", tt.percent, tt.width, got, tt.want)
		}
	}
}

func TestTitledBorder(t *testing.T) {
	result := TitledBorder("Test", "hello", 20, 5, false)
	if !strings.Contains(result, "Test") {
		t.Error("title not found in bordered output")
	}
	if !strings.Contains(result, "hello") {
		t.Error("content not found in bordered output")
	}
	if !strings.Contains(result, "╭") {
		t.Error("top-left corner not found")
	}
	if !strings.Contains(result, "╯") {
		t.Error("bottom-right corner not found")
	}
}

func TestTitledBorderFocused(t *testing.T) {
	result := TitledBorder("Focus", "content", 20, 5, true)
	if !strings.Contains(result, "Focus") {
		t.Error("title not found in focused bordered output")
	}
}

func TestTitledBorderNoTitle(t *testing.T) {
	result := TitledBorder("", "hello", 20, 5, false)
	lines := strings.Split(result, "\n")
	if len(lines) != 5 {
		t.Errorf("expected 5 lines, got %d", len(lines))
	}
}
