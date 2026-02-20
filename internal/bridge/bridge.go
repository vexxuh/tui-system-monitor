package bridge

import (
	"sync"

	lua "github.com/yuin/gopher-lua"
)

// Bridge wraps a gopher-lua LState with a mutex for concurrent safety.
// A single LState is reused across ticks to preserve module-level delta state
// (e.g. prev_cpu_stats in linux.lua).
type Bridge struct {
	mu sync.Mutex
	l  *lua.LState
}

// New creates a Bridge with a new LState and configures package.path
// to include the project's src/ directory.
func New(projectRoot string) (*Bridge, error) {
	L := lua.NewState()

	// Set package.path so Lua can find src.collectors, src.utils, etc.
	pkg := L.GetGlobal("package")
	if tbl, ok := pkg.(*lua.LTable); ok {
		currentPath := lua.LVAsString(tbl.RawGetString("path"))
		newPath := projectRoot + "/src/?.lua;" +
			projectRoot + "/src/?/init.lua;" +
			projectRoot + "/?.lua;" +
			projectRoot + "/?/init.lua;" +
			currentPath
		tbl.RawSetString("path", lua.LString(newPath))
	}

	return &Bridge{l: L}, nil
}

// Close shuts down the LState.
func (b *Bridge) Close() {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.l.Close()
}

// Call executes a Lua function call under the mutex.
// fn should use b.L directly.
func (b *Bridge) Call(fn func(*lua.LState) error) error {
	b.mu.Lock()
	defer b.mu.Unlock()
	return fn(b.l)
}
