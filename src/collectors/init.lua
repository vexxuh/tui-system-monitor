local M = {}

local function detect_platform()
    local f = io.open("/proc/version", "r")
    if f then
        f:close()
        return "linux"
    end
    local handle = io.popen("uname -s 2>/dev/null")
    if handle then
        local result = handle:read("*l")
        handle:close()
        if result and result:match("Darwin") then
            return "macos"
        end
    end
    return "unknown"
end

M.platform = detect_platform()

local collector
if M.platform == "linux" then
    collector = require("src.collectors.linux")
elseif M.platform == "macos" then
    collector = require("src.collectors.macos")
end

function M.collect_all()
    if not collector then
        return {
            cpu = { usage = 0, cores = {}, freq = 0, load = { 0, 0, 0 } },
            memory = { ram_used = 0, ram_total = 0, swap_used = 0, swap_total = 0, cached = 0, buffers = 0 },
            temps = {},
            network = { interface = "none", rx_speed = 0, tx_speed = 0, rx_total = 0, tx_total = 0 },
            disk = { mounts = {}, io_read = 0, io_write = 0, device = "unknown" },
            battery = { percent = 0, status = "Unknown", power = 0, voltage = 0, technology = "" },
            system = { uptime = 0, procs = 0, load = { 0, 0, 0 } },
        }
    end
    return collector.collect_all()
end

return M
