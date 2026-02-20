package sparkline

import "strings"

// blocks are the Unicode block characters U+2581 through U+2588.
var blocks = []rune{'▁', '▂', '▃', '▄', '▅', '▆', '▇', '█'}

// Render produces a sparkline string of the given width from values.
// Empty input returns spaces of the given width.
func Render(values []float64, width int) string {
	if width <= 0 {
		width = 20
	}
	if len(values) == 0 {
		return strings.Repeat(" ", width)
	}

	resampled := Resample(values, width)

	min, max := resampled[0], resampled[0]
	for _, v := range resampled[1:] {
		if v < min {
			min = v
		}
		if v > max {
			max = v
		}
	}

	runes := make([]rune, len(resampled))
	for i, v := range resampled {
		if max == min {
			runes[i] = blocks[3] // middle block for flat line
		} else {
			norm := (v - min) / (max - min)
			idx := int(norm * 7)
			if idx > 7 {
				idx = 7
			}
			runes[i] = blocks[idx]
		}
	}
	return string(runes)
}

// Resample resamples values to targetLen using linear interpolation.
func Resample(values []float64, targetLen int) []float64 {
	if len(values) == 0 {
		return nil
	}
	if len(values) == targetLen {
		out := make([]float64, targetLen)
		copy(out, values)
		return out
	}

	out := make([]float64, targetLen)
	srcLen := len(values)

	for i := range targetLen {
		var srcIdx float64
		if targetLen == 1 {
			srcIdx = float64(srcLen-1) / 2.0
		} else {
			srcIdx = float64(i) * float64(srcLen-1) / float64(targetLen-1)
		}

		lo := int(srcIdx)
		if lo >= srcLen-1 {
			out[i] = values[srcLen-1]
			continue
		}

		frac := srcIdx - float64(lo)
		out[i] = values[lo]*(1-frac) + values[lo+1]*frac
	}

	return out
}
