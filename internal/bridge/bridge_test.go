package bridge

import (
	"os"
	"path/filepath"
	"testing"

	lua "github.com/yuin/gopher-lua"
)

func projectRoot(t *testing.T) string {
	t.Helper()
	// Walk up from the test file to find go.mod
	dir, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	for {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			return dir
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			t.Fatal("could not find project root")
		}
		dir = parent
	}
}

func TestNewBridge(t *testing.T) {
	root := projectRoot(t)
	b, err := New(root)
	if err != nil {
		t.Fatalf("New() error: %v", err)
	}
	defer b.Close()

	// Verify we can execute Lua
	err = b.Call(func(L *lua.LState) error {
		return L.DoString(`x = 1 + 1`)
	})
	if err != nil {
		t.Fatalf("Call() error: %v", err)
	}
}

func TestLoadConfig(t *testing.T) {
	root := projectRoot(t)
	b, err := New(root)
	if err != nil {
		t.Fatalf("New() error: %v", err)
	}
	defer b.Close()

	cfg, err := b.LoadConfig()
	if err != nil {
		t.Fatalf("LoadConfig() error: %v", err)
	}

	if cfg.RefreshInterval <= 0 {
		t.Errorf("refresh_interval should be > 0, got %d", cfg.RefreshInterval)
	}
	if cfg.TempUnit != "C" && cfg.TempUnit != "F" {
		t.Errorf("temp_unit should be C or F, got %q", cfg.TempUnit)
	}
	if cfg.HistorySize <= 0 {
		t.Errorf("history_size should be > 0, got %d", cfg.HistorySize)
	}
}

func TestCollectAll(t *testing.T) {
	root := projectRoot(t)
	b, err := New(root)
	if err != nil {
		t.Fatalf("New() error: %v", err)
	}
	defer b.Close()

	data, err := b.CollectAll()
	if err != nil {
		t.Fatalf("CollectAll() error: %v", err)
	}

	if data.Memory.RAMTotal <= 0 {
		t.Error("ram_total should be > 0")
	}
	if data.Memory.RAMUsed > data.Memory.RAMTotal {
		t.Error("ram_used should be <= ram_total")
	}
	if data.System.Uptime <= 0 {
		t.Error("uptime should be > 0")
	}
	if data.System.Procs <= 0 {
		t.Error("procs should be > 0")
	}
}
