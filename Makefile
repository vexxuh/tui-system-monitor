GO      := go
LUA     := lua
BUSTED  := busted
BINARY  := cy-monitor
GOFLAGS :=

# Luarocks paths for busted tests
LUAROCKS := luarocks
export PATH      := $(shell $(LUAROCKS) path --lr-bin 2>/dev/null):$(PATH)
export LUA_PATH  := $(shell $(LUAROCKS) path --lr-path 2>/dev/null);$(LUA_PATH)
export LUA_CPATH := $(shell $(LUAROCKS) path --lr-cpath 2>/dev/null);$(LUA_CPATH)

.PHONY: help build run test test-go test-lua check clean install

## help: Show this help text (default)
help:
	@echo "cy-monitor — Terminal Hardware Monitor"
	@echo ""
	@echo "Usage:"
	@echo "  make build      Build the Go binary"
	@echo "  make run        Build and run the monitor"
	@echo "  make test       Run all tests (Go + Lua)"
	@echo "  make test-go    Run Go tests only"
	@echo "  make test-lua   Run Lua tests only (busted)"
	@echo "  make check      Verify Lua modules load"
	@echo "  make install    Install Lua test deps (busted)"
	@echo "  make clean      Remove binary and config"
	@echo "  make help       Show this help text"

## build: Compile the Go binary
build:
	$(GO) build $(GOFLAGS) -o $(BINARY) ./cmd/cy-monitor

## run: Build and launch the monitor
run: build
	./$(BINARY)

## test: Run all tests
test: test-go test-lua

## test-go: Run Go unit tests
test-go:
	$(GO) test ./... -v

## test-lua: Run Lua tests via busted
test-lua:
	$(BUSTED)

## check: Verify Lua modules load without starting the TUI
check:
	@$(LUA) -e 'package.path="./?/init.lua;./?.lua;"..package.path;for _,m in ipairs({"src.utils.format","src.utils.sparkline","src.config","src.collectors","src.collectors.history"})do local ok,err=pcall(require,m);if ok then print("  OK  "..m)else print("  FAIL "..m..": "..err);os.exit(1)end end;print("All modules loaded.")'

## install: Install Lua test dependencies
install:
	@echo "Installing Lua test deps..."
	@$(LUAROCKS) install --local busted > /dev/null 2>&1 || { echo "FAILED: busted"; exit 1; }
	@echo "Done."
	@echo "Go deps managed by go.mod (run 'go mod tidy')"

## clean: Remove binary and user config
clean:
	rm -f $(BINARY)
	rm -rf $(HOME)/.config/cy-monitor
