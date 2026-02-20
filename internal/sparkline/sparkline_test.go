package sparkline

import (
	"math"
	"strings"
	"testing"
	"unicode/utf8"
)

func TestRenderAscending(t *testing.T) {
	values := []float64{0, 25, 50, 75, 100}
	got := Render(values, 5)
	if utf8.RuneCountInString(got) != 5 {
		t.Errorf("expected 5 runes, got %d", utf8.RuneCountInString(got))
	}
	runes := []rune(got)
	if runes[0] != '▁' {
		t.Errorf("first rune should be ▁, got %c", runes[0])
	}
	if runes[4] != '█' {
		t.Errorf("last rune should be █, got %c", runes[4])
	}
}

func TestRenderConstant(t *testing.T) {
	values := []float64{50, 50, 50}
	got := Render(values, 3)
	runes := []rune(got)
	for i, r := range runes {
		if r != '▄' {
			t.Errorf("index %d: expected ▄, got %c", i, r)
		}
	}
}

func TestRenderEmpty(t *testing.T) {
	got := Render(nil, 10)
	if got != strings.Repeat(" ", 10) {
		t.Errorf("empty render = %q", got)
	}
}

func TestResampleIdentity(t *testing.T) {
	values := []float64{1, 2, 3}
	got := Resample(values, 3)
	for i, v := range got {
		if v != values[i] {
			t.Errorf("index %d: got %f, want %f", i, v, values[i])
		}
	}
}

func TestResampleDownsample(t *testing.T) {
	values := []float64{0, 25, 50, 75, 100}
	got := Resample(values, 3)
	if len(got) != 3 {
		t.Fatalf("expected len 3, got %d", len(got))
	}
	if got[0] != 0 {
		t.Errorf("first = %f", got[0])
	}
	if got[2] != 100 {
		t.Errorf("last = %f", got[2])
	}
}

func TestResampleUpsample(t *testing.T) {
	values := []float64{0, 100}
	got := Resample(values, 3)
	if len(got) != 3 {
		t.Fatalf("expected len 3, got %d", len(got))
	}
	if math.Abs(got[1]-50) > 0.01 {
		t.Errorf("midpoint = %f, want 50", got[1])
	}
}

func TestResampleEmpty(t *testing.T) {
	got := Resample(nil, 5)
	if got != nil {
		t.Errorf("expected nil, got %v", got)
	}
}
