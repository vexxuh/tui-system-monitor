package format

import (
	"testing"
)

func TestBytes(t *testing.T) {
	tests := []struct {
		value float64
		valid bool
		want  string
	}{
		{1.5e9, true, "1.5G"},
		{2.5e6, true, "2.5M"},
		{1500, true, "1.5K"},
		{500, true, "500B"},
		{0, true, "0B"},
		{-5, true, "0B"},
		{0, false, "--"},
	}
	for _, tt := range tests {
		got := Bytes(tt.value, tt.valid)
		if got != tt.want {
			t.Errorf("Bytes(%v, %v) = %q, want %q", tt.value, tt.valid, got, tt.want)
		}
	}
}

func TestBytesSpeed(t *testing.T) {
	if got := BytesSpeed(1e6, true); got != "1.0M/s" {
		t.Errorf("BytesSpeed(1e6) = %q", got)
	}
	if got := BytesSpeed(0, false); got != "-- /s" {
		t.Errorf("BytesSpeed invalid = %q", got)
	}
}

func TestTemp(t *testing.T) {
	if got := Temp(46, true, "C"); got != "46°C" {
		t.Errorf("Temp C = %q", got)
	}
	if got := Temp(100, true, "F"); got != "212°F" {
		t.Errorf("Temp F = %q", got)
	}
	if got := Temp(0, false, "C"); got != "--" {
		t.Errorf("Temp invalid = %q", got)
	}
}

func TestDuration(t *testing.T) {
	tests := []struct {
		secs  float64
		valid bool
		want  string
	}{
		{90000, true, "1d 1h 0m"},
		{9000, true, "2h 30m"},
		{300, true, "5m"},
		{0, true, "--"},
		{-10, true, "--"},
		{0, false, "--"},
	}
	for _, tt := range tests {
		got := Duration(tt.secs, tt.valid)
		if got != tt.want {
			t.Errorf("Duration(%v, %v) = %q, want %q", tt.secs, tt.valid, got, tt.want)
		}
	}
}

func TestFreq(t *testing.T) {
	if got := Freq(4200, true); got != "4.2 GHz" {
		t.Errorf("Freq GHz = %q", got)
	}
	if got := Freq(800, true); got != "800 MHz" {
		t.Errorf("Freq MHz = %q", got)
	}
	if got := Freq(0, false); got != "--" {
		t.Errorf("Freq invalid = %q", got)
	}
}

func TestPercent(t *testing.T) {
	if got := Percent(42, true); got != "42%" {
		t.Errorf("Percent = %q", got)
	}
	if got := Percent(0, false); got != "--%" {
		t.Errorf("Percent invalid = %q", got)
	}
}

func TestPadRight(t *testing.T) {
	if got := PadRight("hi", 5); got != "hi   " {
		t.Errorf("PadRight pad = %q", got)
	}
	if got := PadRight("hello world", 5); got != "hello" {
		t.Errorf("PadRight truncate = %q", got)
	}
	if got := PadRight("exact", 5); got != "exact" {
		t.Errorf("PadRight exact = %q", got)
	}
}

func TestTruncate(t *testing.T) {
	if got := Truncate("hello", 10); got != "hello" {
		t.Errorf("Truncate short = %q", got)
	}
	if got := Truncate("hello world", 5); got != "hell…" {
		t.Errorf("Truncate long = %q", got)
	}
	if got := Truncate("ab", 1); got != "…" {
		t.Errorf("Truncate min = %q", got)
	}
	if got := Truncate("ab", 0); got != "" {
		t.Errorf("Truncate zero = %q", got)
	}
}
