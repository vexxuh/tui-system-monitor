package panel

import (
	"github.com/vexxuh/cy-monitor/internal/model"
)

// PanelState holds all data needed for panel rendering.
type PanelState struct {
	Data     model.AllData
	Config   model.AppConfig
	History  HistoryState
	Focused  int // 1-6
	Expanded int // 0 = none, 1-6 = expanded panel
}

// HistoryState holds ring buffers for sparkline history.
type HistoryState struct {
	CPU   model.RingBuffer
	NetRx model.RingBuffer
	NetTx model.RingBuffer
	Temps map[string]model.RingBuffer
}

// NewHistoryState creates a new history state with the given buffer size.
func NewHistoryState(size int) HistoryState {
	return HistoryState{
		CPU:   model.NewRingBuffer(size),
		NetRx: model.NewRingBuffer(size),
		NetTx: model.NewRingBuffer(size),
		Temps: make(map[string]model.RingBuffer),
	}
}

// PushData updates history buffers with new data, returning a new HistoryState.
func (h HistoryState) PushData(data model.AllData, historySize int) HistoryState {
	newH := HistoryState{
		CPU:   h.CPU.Push(data.CPU.Usage),
		NetRx: h.NetRx.Push(data.Network.RxSpeed),
		NetTx: h.NetTx.Push(data.Network.TxSpeed),
		Temps: make(map[string]model.RingBuffer),
	}

	// Copy existing temp histories
	for k, v := range h.Temps {
		newH.Temps[k] = v
	}

	// Push new temp values
	for _, sensor := range data.Temps {
		buf, ok := newH.Temps[sensor.Label]
		if !ok {
			buf = model.NewRingBuffer(historySize)
		}
		newH.Temps[sensor.Label] = buf.Push(sensor.TempC)
	}

	return newH
}

// PanelName returns the display name for a panel index (1-6).
func PanelName(idx int) string {
	names := [6]string{"CPU", "Memory", "Temps", "Network", "Disk", "System"}
	if idx >= 1 && idx <= 6 {
		return names[idx-1]
	}
	return ""
}
