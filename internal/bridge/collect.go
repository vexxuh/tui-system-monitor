package bridge

import (
	"github.com/vexxuh/cy-monitor/internal/model"
	lua "github.com/yuin/gopher-lua"
)

// CollectAll calls src.collectors.collect_all() and returns structured data.
func (b *Bridge) CollectAll() (model.AllData, error) {
	var data model.AllData

	err := b.Call(func(L *lua.LState) error {
		if err := L.DoString(`
			local collectors = require("src.collectors")
			return collectors.collect_all()
		`); err != nil {
			return err
		}

		tbl := L.CheckTable(-1)
		L.Pop(1)

		data.CPU = parseCPU(tbl.RawGetString("cpu"))
		data.Memory = parseMemory(tbl.RawGetString("memory"))
		data.Temps = parseTemps(tbl.RawGetString("temps"))
		data.Network = parseNetwork(tbl.RawGetString("network"))
		data.Disk = parseDisk(tbl.RawGetString("disk"))
		data.Battery = parseBattery(tbl.RawGetString("battery"))
		data.System = parseSystem(tbl.RawGetString("system"))

		return nil
	})

	return data, err
}

func parseCPU(lv lua.LValue) model.CPUData {
	tbl, ok := lv.(*lua.LTable)
	if !ok {
		return model.CPUData{}
	}

	cpu := model.CPUData{
		Usage: getNumber(tbl, "usage"),
		Freq:  getNumber(tbl, "freq"),
	}

	if loadTbl, ok := tbl.RawGetString("load").(*lua.LTable); ok {
		cpu.Load = getLoad(loadTbl)
	}

	if coresTbl, ok := tbl.RawGetString("cores").(*lua.LTable); ok {
		coresTbl.ForEach(func(_, v lua.LValue) {
			if n, ok := v.(lua.LNumber); ok {
				cpu.Cores = append(cpu.Cores, float64(n))
			}
		})
	}

	return cpu
}

func parseMemory(lv lua.LValue) model.MemoryData {
	tbl, ok := lv.(*lua.LTable)
	if !ok {
		return model.MemoryData{}
	}
	return model.MemoryData{
		RAMUsed:   int64(getNumber(tbl, "ram_used")),
		RAMTotal:  int64(getNumber(tbl, "ram_total")),
		SwapUsed:  int64(getNumber(tbl, "swap_used")),
		SwapTotal: int64(getNumber(tbl, "swap_total")),
		Cached:    int64(getNumber(tbl, "cached")),
		Buffers:   int64(getNumber(tbl, "buffers")),
	}
}

func parseTemps(lv lua.LValue) []model.TempSensor {
	tbl, ok := lv.(*lua.LTable)
	if !ok {
		return nil
	}

	var sensors []model.TempSensor
	tbl.ForEach(func(_, v lua.LValue) {
		if st, ok := v.(*lua.LTable); ok {
			sensors = append(sensors, model.TempSensor{
				Label: getString(st, "label"),
				TempC: getNumber(st, "temp_c"),
			})
		}
	})
	return sensors
}

func parseNetwork(lv lua.LValue) model.NetworkData {
	tbl, ok := lv.(*lua.LTable)
	if !ok {
		return model.NetworkData{}
	}
	return model.NetworkData{
		Interface: getString(tbl, "interface"),
		RxSpeed:   getNumber(tbl, "rx_speed"),
		TxSpeed:   getNumber(tbl, "tx_speed"),
		RxTotal:   int64(getNumber(tbl, "rx_total")),
		TxTotal:   int64(getNumber(tbl, "tx_total")),
	}
}

func parseDisk(lv lua.LValue) model.DiskData {
	tbl, ok := lv.(*lua.LTable)
	if !ok {
		return model.DiskData{}
	}

	disk := model.DiskData{
		IORead:  getNumber(tbl, "io_read"),
		IOWrite: getNumber(tbl, "io_write"),
		Device:  getString(tbl, "device"),
	}

	if mountsTbl, ok := tbl.RawGetString("mounts").(*lua.LTable); ok {
		mountsTbl.ForEach(func(_, v lua.LValue) {
			if mt, ok := v.(*lua.LTable); ok {
				disk.Mounts = append(disk.Mounts, model.DiskMount{
					Mount:   getString(mt, "mount"),
					Total:   int64(getNumber(mt, "total")),
					Used:    int64(getNumber(mt, "used")),
					Percent: getNumber(mt, "percent"),
				})
			}
		})
	}

	return disk
}

func parseBattery(lv lua.LValue) model.BatteryData {
	tbl, ok := lv.(*lua.LTable)
	if !ok {
		return model.BatteryData{}
	}
	return model.BatteryData{
		Percent:    getNumber(tbl, "percent"),
		Status:     getString(tbl, "status"),
		Power:      getNumber(tbl, "power"),
		Voltage:    getNumber(tbl, "voltage"),
		Technology: getString(tbl, "technology"),
	}
}

func parseSystem(lv lua.LValue) model.SystemData {
	tbl, ok := lv.(*lua.LTable)
	if !ok {
		return model.SystemData{}
	}

	sys := model.SystemData{
		Uptime: getNumber(tbl, "uptime"),
		Procs:  int(getNumber(tbl, "procs")),
	}

	if loadTbl, ok := tbl.RawGetString("load").(*lua.LTable); ok {
		sys.Load = getLoad(loadTbl)
	}

	return sys
}

// Helper functions for extracting values from Lua tables.

func getNumber(tbl *lua.LTable, key string) float64 {
	v := tbl.RawGetString(key)
	if n, ok := v.(lua.LNumber); ok {
		return float64(n)
	}
	return 0
}

func getString(tbl *lua.LTable, key string) string {
	v := tbl.RawGetString(key)
	if s, ok := v.(lua.LString); ok {
		return string(s)
	}
	return ""
}

func getLoad(tbl *lua.LTable) [3]float64 {
	var load [3]float64
	for i := range 3 {
		v := tbl.RawGetInt(i + 1)
		if n, ok := v.(lua.LNumber); ok {
			load[i] = float64(n)
		}
	}
	return load
}
