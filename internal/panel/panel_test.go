package panel

import (
	"strings"
	"testing"

	"github.com/vexxuh/cy-monitor/internal/model"
)

func testState() PanelState {
	return PanelState{
		Data: model.AllData{
			CPU: model.CPUData{
				Usage: 42.5,
				Cores: []float64{30, 55, 40, 65},
				Freq:  3600,
				Load:  [3]float64{0.42, 0.38, 0.30},
			},
			Memory: model.MemoryData{
				RAMUsed:   4 * 1024 * 1024 * 1024,
				RAMTotal:  16 * 1024 * 1024 * 1024,
				SwapUsed:  512 * 1024 * 1024,
				SwapTotal: 2 * 1024 * 1024 * 1024,
				Cached:    2 * 1024 * 1024 * 1024,
				Buffers:   100 * 1024 * 1024,
			},
			Temps: []model.TempSensor{
				{Label: "CPU", TempC: 55},
				{Label: "GPU", TempC: 70},
			},
			Network: model.NetworkData{
				Interface: "eth0",
				RxSpeed:   1e6,
				TxSpeed:   500e3,
				RxTotal:   1e9,
				TxTotal:   200e6,
			},
			Disk: model.DiskData{
				Mounts: []model.DiskMount{
					{Mount: "/", Total: 500e9, Used: 250e9, Percent: 50},
					{Mount: "/home", Total: 1e12, Used: 400e9, Percent: 40},
				},
				IORead:  10e6,
				IOWrite: 2e6,
				Device:  "nvme0n1",
			},
			Battery: model.BatteryData{
				Percent:    85,
				Status:     "Charging",
				Power:      15.2,
				Voltage:    12.5,
				Technology: "Li-ion",
			},
			System: model.SystemData{
				Uptime: 9000,
				Procs:  312,
				Load:   [3]float64{0.42, 0.38, 0.30},
			},
		},
		Config:  model.DefaultConfig(),
		History: NewHistoryState(60),
		Focused: 1,
	}
}

func TestRenderCPU(t *testing.T) {
	s := testState()
	out := RenderCPU(s, 40)
	if !strings.Contains(out, "42%") {
		t.Error("should contain CPU usage percentage")
	}
	if !strings.Contains(out, "3.6 GHz") {
		t.Error("should contain frequency")
	}
	if !strings.Contains(out, "Load:") {
		t.Error("should contain load averages")
	}
}

func TestRenderMemory(t *testing.T) {
	s := testState()
	out := RenderMemory(s, 40)
	if !strings.Contains(out, "RAM") {
		t.Error("should contain RAM label")
	}
	if !strings.Contains(out, "Swap") {
		t.Error("should contain Swap label")
	}
	if !strings.Contains(out, "Cached:") {
		t.Error("should contain Cached")
	}
}

func TestRenderTemps(t *testing.T) {
	s := testState()
	out := RenderTemps(s, 40)
	if !strings.Contains(out, "CPU") {
		t.Error("should contain CPU sensor")
	}
	if !strings.Contains(out, "55°C") {
		t.Error("should contain temperature value")
	}
}

func TestRenderTempsEmpty(t *testing.T) {
	s := testState()
	s.Data.Temps = nil
	out := RenderTemps(s, 40)
	if out != "No sensors found" {
		t.Errorf("expected 'No sensors found', got %q", out)
	}
}

func TestRenderNetwork(t *testing.T) {
	s := testState()
	out := RenderNetwork(s, 40)
	if !strings.Contains(out, "eth0") {
		t.Error("should contain interface name")
	}
	if !strings.Contains(out, "Total:") {
		t.Error("should contain total")
	}
}

func TestRenderDisk(t *testing.T) {
	s := testState()
	out := RenderDisk(s, 40)
	if !strings.Contains(out, "nvme0n1") {
		t.Error("should contain device name")
	}
	if !strings.Contains(out, "/") {
		t.Error("should contain mount point")
	}
}

func TestRenderSystem(t *testing.T) {
	s := testState()
	out := RenderSystem(s, 40)
	if !strings.Contains(out, "Charging") {
		t.Error("should contain battery status")
	}
	if !strings.Contains(out, "Uptime:") {
		t.Error("should contain uptime")
	}
	if !strings.Contains(out, "Procs:") {
		t.Error("should contain procs")
	}
}

func TestRenderSystemNoBattery(t *testing.T) {
	s := testState()
	s.Data.Battery.Percent = -1
	out := RenderSystem(s, 40)
	if !strings.Contains(out, "AC Power") {
		t.Error("should show AC Power when no battery")
	}
}

func TestPanelName(t *testing.T) {
	names := []string{"CPU", "Memory", "Temps", "Network", "Disk", "System"}
	for i, want := range names {
		got := PanelName(i + 1)
		if got != want {
			t.Errorf("PanelName(%d) = %q, want %q", i+1, got, want)
		}
	}
	if got := PanelName(0); got != "" {
		t.Errorf("PanelName(0) = %q, want empty", got)
	}
}

func TestHistoryPushData(t *testing.T) {
	h := NewHistoryState(60)
	s := testState()
	h2 := h.PushData(s.Data, 60)

	if h2.CPU.Len() != 1 {
		t.Errorf("CPU history should have 1 entry, got %d", h2.CPU.Len())
	}
	if h2.NetRx.Len() != 1 {
		t.Errorf("NetRx history should have 1 entry, got %d", h2.NetRx.Len())
	}
	if _, ok := h2.Temps["CPU"]; !ok {
		t.Error("should have CPU temp history")
	}
	// Original unchanged
	if h.CPU.Len() != 0 {
		t.Error("original history should be unchanged")
	}
}
