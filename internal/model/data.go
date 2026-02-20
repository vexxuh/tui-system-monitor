package model

// CPUData holds CPU usage, per-core usage, frequency, and load averages.
type CPUData struct {
	Usage float64
	Cores []float64
	Freq  float64
	Load  [3]float64
}

// MemoryData holds RAM and swap usage in bytes.
type MemoryData struct {
	RAMUsed   int64
	RAMTotal  int64
	SwapUsed  int64
	SwapTotal int64
	Cached    int64
	Buffers   int64
}

// TempSensor holds a single temperature sensor reading.
type TempSensor struct {
	Label string
	TempC float64
}

// NetworkData holds interface stats and transfer speeds.
type NetworkData struct {
	Interface string
	RxSpeed   float64
	TxSpeed   float64
	RxTotal   int64
	TxTotal   int64
}

// DiskMount holds mount point usage.
type DiskMount struct {
	Mount   string
	Total   int64
	Used    int64
	Percent float64
}

// DiskData holds mount points and I/O speeds.
type DiskData struct {
	Mounts  []DiskMount
	IORead  float64
	IOWrite float64
	Device  string
}

// BatteryData holds battery/power supply info.
type BatteryData struct {
	Percent    float64
	Status     string
	Power      float64
	Voltage    float64
	Technology string
}

// SystemData holds system uptime, process count, and load.
type SystemData struct {
	Uptime float64
	Procs  int
	Load   [3]float64
}

// AllData is the complete snapshot returned by a collect cycle.
type AllData struct {
	CPU     CPUData
	Memory  MemoryData
	Temps   []TempSensor
	Network NetworkData
	Disk    DiskData
	Battery BatteryData
	System  SystemData
}
