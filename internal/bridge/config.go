package bridge

import (
	"github.com/vexxuh/cy-monitor/internal/model"
	lua "github.com/yuin/gopher-lua"
)

// LoadConfig calls src.config.load() via Lua and returns an AppConfig.
func (b *Bridge) LoadConfig() (model.AppConfig, error) {
	cfg := model.DefaultConfig()

	err := b.Call(func(L *lua.LState) error {
		if err := L.DoString(`
			local config = require("src.config")
			return config.load()
		`); err != nil {
			return err
		}

		tbl := L.CheckTable(-1)
		L.Pop(1)

		if v := tbl.RawGetString("refresh_interval"); v.Type() == lua.LTNumber {
			cfg.RefreshInterval = int(lua.LVAsNumber(v))
		}
		if v := tbl.RawGetString("temp_unit"); v.Type() == lua.LTString {
			cfg.TempUnit = lua.LVAsString(v)
		}
		if v := tbl.RawGetString("history_size"); v.Type() == lua.LTNumber {
			cfg.HistorySize = int(lua.LVAsNumber(v))
		}

		return nil
	})

	return cfg, err
}
