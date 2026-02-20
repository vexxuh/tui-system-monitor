local M = {}

-- Persistent state for delta calculations
local prev_cpu_stats = nil
local prev_net_stats = nil
local prev_net_time = nil
local prev_disk_stats = nil
local prev_disk_time = nil

local function read_file(path)
    local f = io.open(path, "r")
    if not f then return nil end
    local content = f:read("*a")
    f:close()
    return content
end

local function read_first_line(path)
    local f = io.open(path, "r")
    if not f then return nil end
    local line = f:read("*l")
    f:close()
    return line
end

local function parse_cpu()
    local content = read_file("/proc/stat")
    if not content then
        return { usage = 0, cores = {}, freq = 0, load = { 0, 0, 0 } }
    end

    local cores = {}
    local total_usage = 0

    -- Parse aggregate and per-core lines
    local cpu_lines = {}
    for line in content:gmatch("[^\n]+") do
        if line:match("^cpu") then
            cpu_lines[#cpu_lines + 1] = line
        end
    end

    local current_stats = {}
    for _, line in ipairs(cpu_lines) do
        local name, user, nice, system, idle, iowait, irq, softirq, steal =
            line:match("^(%S+)%s+(%d+)%s+(%d+)%s+(%d+)%s+(%d+)%s+(%d+)%s+(%d+)%s+(%d+)%s+(%d+)")
        if name then
            local total = user + nice + system + idle + iowait + irq + softirq + steal
            local busy = total - idle - iowait
            current_stats[name] = { busy = busy, total = total }
        end
    end

    if prev_cpu_stats then
        for name, cur in pairs(current_stats) do
            local prev = prev_cpu_stats[name]
            if prev then
                local d_total = cur.total - prev.total
                local d_busy = cur.busy - prev.busy
                local usage = 0
                if d_total > 0 then
                    usage = (d_busy / d_total) * 100
                end
                if name == "cpu" then
                    total_usage = usage
                else
                    local core_num = tonumber(name:match("cpu(%d+)"))
                    if core_num then
                        cores[core_num + 1] = usage
                    end
                end
            end
        end
    end
    prev_cpu_stats = current_stats

    -- CPU frequency
    local freq = 0
    local freq_content = read_file("/proc/cpuinfo")
    if freq_content then
        local count = 0
        local sum = 0
        for mhz in freq_content:gmatch("cpu MHz%s*:%s*([%d%.]+)") do
            sum = sum + tonumber(mhz)
            count = count + 1
        end
        if count > 0 then
            freq = sum / count
        end
    end

    -- Try scaling_cur_freq as fallback
    if freq == 0 then
        local f_line = read_first_line("/sys/devices/system/cpu/cpu0/cpufreq/scaling_cur_freq")
        if f_line then
            freq = (tonumber(f_line) or 0) / 1000
        end
    end

    -- Load average
    local load = { 0, 0, 0 }
    local loadavg = read_first_line("/proc/loadavg")
    if loadavg then
        local l1, l5, l15 = loadavg:match("([%d%.]+)%s+([%d%.]+)%s+([%d%.]+)")
        if l1 then
            load = { tonumber(l1), tonumber(l5), tonumber(l15) }
        end
    end

    return {
        usage = total_usage,
        cores = cores,
        freq = freq,
        load = load,
    }
end

local function parse_memory()
    local content = read_file("/proc/meminfo")
    if not content then
        return { ram_used = 0, ram_total = 0, swap_used = 0, swap_total = 0, cached = 0, buffers = 0 }
    end

    local values = {}
    for key, val in content:gmatch("(%S+):%s+(%d+)") do
        values[key] = tonumber(val) * 1024  -- Convert kB to bytes
    end

    local total = values["MemTotal"] or 0
    local free = values["MemFree"] or 0
    local buffers = values["Buffers"] or 0
    local cached = values["Cached"] or 0
    local sreclaimable = values["SReclaimable"] or 0
    local shmem = values["Shmem"] or 0

    -- Used = Total - Free - Buffers - Cached - SReclaimable + Shmem (htop formula)
    local used = total - free - buffers - cached - sreclaimable + shmem

    local swap_total = values["SwapTotal"] or 0
    local swap_free = values["SwapFree"] or 0

    return {
        ram_used = used,
        ram_total = total,
        swap_used = swap_total - swap_free,
        swap_total = swap_total,
        cached = cached + sreclaimable,
        buffers = buffers,
    }
end

local function parse_temps()
    local temps = {}
    local hwmon_base = "/sys/class/hwmon"

    -- Try to enumerate hwmon devices
    local handle = io.popen('ls "' .. hwmon_base .. '" 2>/dev/null')
    if not handle then return temps end

    local devices = handle:read("*a")
    handle:close()

    for device in devices:gmatch("%S+") do
        local dev_path = hwmon_base .. "/" .. device
        local name_line = read_first_line(dev_path .. "/name")
        local label = name_line or device

        -- Read temp inputs
        local i = 1
        while true do
            local temp_path = dev_path .. "/temp" .. i .. "_input"
            local temp_line = read_first_line(temp_path)
            if not temp_line then break end

            local temp_c = (tonumber(temp_line) or 0) / 1000
            local temp_label_path = dev_path .. "/temp" .. i .. "_label"
            local temp_label = read_first_line(temp_label_path) or (label .. " #" .. i)

            if temp_c > 0 and temp_c < 150 then
                temps[#temps + 1] = {
                    label = temp_label,
                    temp_c = temp_c,
                }
            end
            i = i + 1
        end
    end

    return temps
end

local function parse_network()
    local content = read_file("/proc/net/dev")
    if not content then
        return { interface = "none", rx_speed = 0, tx_speed = 0, rx_total = 0, tx_total = 0 }
    end

    local now = os.time()
    local current_stats = {}
    local best_iface = nil
    local best_rx = 0

    for line in content:gmatch("[^\n]+") do
        local iface, rx_bytes, _, _, _, _, _, _, _, tx_bytes =
            line:match("^%s*(%S+):%s*(%d+)%s+(%d+)%s+(%d+)%s+(%d+)%s+(%d+)%s+(%d+)%s+(%d+)%s+(%d+)%s+(%d+)")
        if iface and iface ~= "lo" then
            rx_bytes = tonumber(rx_bytes)
            tx_bytes = tonumber(tx_bytes)
            current_stats[iface] = { rx = rx_bytes, tx = tx_bytes }
            if rx_bytes > best_rx then
                best_rx = rx_bytes
                best_iface = iface
            end
        end
    end

    if not best_iface then
        return { interface = "none", rx_speed = 0, tx_speed = 0, rx_total = 0, tx_total = 0 }
    end

    local rx_speed = 0
    local tx_speed = 0
    if prev_net_stats and prev_net_time then
        local dt = now - prev_net_time
        if dt > 0 then
            local prev = prev_net_stats[best_iface]
            local cur = current_stats[best_iface]
            if prev and cur then
                rx_speed = (cur.rx - prev.rx) / dt
                tx_speed = (cur.tx - prev.tx) / dt
            end
        end
    end
    prev_net_stats = current_stats
    prev_net_time = now

    local cur = current_stats[best_iface]
    return {
        interface = best_iface,
        rx_speed = math.max(0, rx_speed),
        tx_speed = math.max(0, tx_speed),
        rx_total = cur.rx,
        tx_total = cur.tx,
    }
end

local function parse_disk()
    -- Mount points via df
    local mounts = {}
    local handle = io.popen("df -B1 --output=target,size,used,avail,pcent 2>/dev/null")
    if handle then
        local output = handle:read("*a")
        handle:close()
        local first = true
        for line in output:gmatch("[^\n]+") do
            if first then
                first = false
            else
                local target, size, used, _, pcent =
                    line:match("^(%S+)%s+(%d+)%s+(%d+)%s+(%d+)%s+(%d+)%%")
                if target and (target == "/" or target:match("^/boot") or target:match("^/home")) then
                    mounts[#mounts + 1] = {
                        mount = target,
                        total = tonumber(size),
                        used = tonumber(used),
                        percent = tonumber(pcent),
                    }
                end
            end
        end
    end

    -- Disk I/O from /proc/diskstats
    local device = "unknown"
    local io_read = 0
    local io_write = 0
    local content = read_file("/proc/diskstats")
    if content then
        local now = os.time()
        local current_disk = {}
        for line in content:gmatch("[^\n]+") do
            local dev, rd_sectors, wr_sectors =
                line:match("%s+%d+%s+%d+%s+(%S+)%s+%d+%s+%d+%s+(%d+)%s+%d+%s+%d+%s+%d+%s+(%d+)")
            if dev and (dev:match("^nvme%d+n%d+$") or dev:match("^sd%a$") or dev:match("^vd%a$")) then
                current_disk[dev] = {
                    rd = tonumber(rd_sectors) * 512,
                    wr = tonumber(wr_sectors) * 512,
                }
                if device == "unknown" then device = dev end
            end
        end

        if prev_disk_stats and prev_disk_time then
            local dt = now - prev_disk_time
            if dt > 0 and current_disk[device] and prev_disk_stats[device] then
                io_read = (current_disk[device].rd - prev_disk_stats[device].rd) / dt
                io_write = (current_disk[device].wr - prev_disk_stats[device].wr) / dt
            end
        end
        prev_disk_stats = current_disk
        prev_disk_time = now
    end

    return {
        mounts = mounts,
        io_read = math.max(0, io_read),
        io_write = math.max(0, io_write),
        device = device,
    }
end

local function parse_battery()
    local base = "/sys/class/power_supply"
    local bat_path = nil

    local handle = io.popen('ls "' .. base .. '" 2>/dev/null')
    if handle then
        local output = handle:read("*a")
        handle:close()
        for dev in output:gmatch("%S+") do
            local type_line = read_first_line(base .. "/" .. dev .. "/type")
            if type_line and type_line:match("Battery") then
                bat_path = base .. "/" .. dev
                break
            end
        end
    end

    if not bat_path then
        return { percent = -1, status = "No battery", power = 0, voltage = 0, technology = "" }
    end

    local capacity = tonumber(read_first_line(bat_path .. "/capacity") or "0")
    local status = read_first_line(bat_path .. "/status") or "Unknown"
    local voltage = (tonumber(read_first_line(bat_path .. "/voltage_now") or "0")) / 1e6
    local power = (tonumber(read_first_line(bat_path .. "/power_now") or "0")) / 1e6
    local tech = read_first_line(bat_path .. "/technology") or ""

    return {
        percent = capacity,
        status = status,
        power = power,
        voltage = voltage,
        technology = tech,
    }
end

local function parse_system()
    local load = { 0, 0, 0 }
    local loadavg = read_first_line("/proc/loadavg")
    if loadavg then
        local l1, l5, l15 = loadavg:match("([%d%.]+)%s+([%d%.]+)%s+([%d%.]+)")
        if l1 then
            load = { tonumber(l1), tonumber(l5), tonumber(l15) }
        end
    end

    local uptime = 0
    local uptime_line = read_first_line("/proc/uptime")
    if uptime_line then
        uptime = tonumber(uptime_line:match("([%d%.]+)")) or 0
    end

    -- Process count from /proc/loadavg (4th field: running/total)
    local procs = 0
    if loadavg then
        local total_procs = loadavg:match("%d+/(%d+)")
        if total_procs then
            procs = tonumber(total_procs) or 0
        end
    end

    return {
        uptime = uptime,
        procs = procs,
        load = load,
    }
end

function M.collect_all()
    return {
        cpu = parse_cpu(),
        memory = parse_memory(),
        temps = parse_temps(),
        network = parse_network(),
        disk = parse_disk(),
        battery = parse_battery(),
        system = parse_system(),
    }
end

return M
