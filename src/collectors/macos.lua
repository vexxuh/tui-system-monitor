local M = {}

local prev_cpu_ticks = nil
local prev_net_stats = nil
local prev_net_time = nil

local function exec(cmd)
    local handle = io.popen(cmd .. " 2>/dev/null")
    if not handle then return nil end
    local result = handle:read("*a")
    handle:close()
    return result
end

local function exec_line(cmd)
    local handle = io.popen(cmd .. " 2>/dev/null")
    if not handle then return nil end
    local result = handle:read("*l")
    handle:close()
    return result
end

local function sysctl_int(key)
    local line = exec_line("sysctl -n " .. key)
    return tonumber(line) or 0
end

local function parse_cpu()
    -- CPU ticks: user, system, idle, nice
    local ticks_str = exec_line("sysctl -n kern.cp_time")
    local usage = 0
    if ticks_str then
        local parts = {}
        for n in ticks_str:gmatch("%d+") do
            parts[#parts + 1] = tonumber(n)
        end
        if #parts >= 4 and prev_cpu_ticks then
            local d_user = parts[1] - prev_cpu_ticks[1]
            local d_sys = parts[2] - prev_cpu_ticks[2]
            local d_idle = parts[3] - prev_cpu_ticks[3]
            local d_nice = parts[4] - prev_cpu_ticks[4]
            local d_total = d_user + d_sys + d_idle + d_nice
            if d_total > 0 then
                usage = ((d_user + d_sys + d_nice) / d_total) * 100
            end
        end
        if #parts >= 4 then
            prev_cpu_ticks = parts
        end
    end

    local freq = sysctl_int("hw.cpufrequency") / 1e6  -- Hz to MHz
    if freq == 0 then
        freq = sysctl_int("hw.cpufrequency_max") / 1e6
    end

    local load = { 0, 0, 0 }
    local load_str = exec_line("sysctl -n vm.loadavg")
    if load_str then
        local l1, l5, l15 = load_str:match("([%d%.]+)%s+([%d%.]+)%s+([%d%.]+)")
        if l1 then
            load = { tonumber(l1), tonumber(l5), tonumber(l15) }
        end
    end

    local ncpu = sysctl_int("hw.ncpu")
    local cores = {}
    for i = 1, ncpu do
        cores[i] = usage  -- macOS doesn't easily expose per-core without powermetrics
    end

    return {
        usage = usage,
        cores = cores,
        freq = freq,
        load = load,
    }
end

local function parse_memory()
    local page_size = sysctl_int("hw.pagesize")
    if page_size == 0 then page_size = 4096 end
    local total = sysctl_int("hw.memsize")

    local vm_output = exec("vm_stat")
    local pages = {}
    if vm_output then
        for key, val in vm_output:gmatch('"?([^":\n]+)"?:%s+(%d+)') do
            pages[key:match("^%s*(.-)%s*$")] = tonumber(val)
        end
    end

    local active = (pages["Pages active"] or 0) * page_size
    local inactive = (pages["Pages inactive"] or 0) * page_size
    local speculative = (pages["Pages speculative"] or 0) * page_size
    local wired = (pages["Pages wired down"] or 0) * page_size
    local compressed = (pages["Pages occupied by compressor"] or 0) * page_size
    local cached = inactive + speculative

    local used = active + wired + compressed

    -- macOS swap
    local swap_output = exec("sysctl -n vm.swapusage")
    local swap_total = 0
    local swap_used = 0
    if swap_output then
        local st = swap_output:match("total%s*=%s*([%d%.]+)M")
        local su = swap_output:match("used%s*=%s*([%d%.]+)M")
        if st then swap_total = tonumber(st) * 1024 * 1024 end
        if su then swap_used = tonumber(su) * 1024 * 1024 end
    end

    return {
        ram_used = used,
        ram_total = total,
        swap_used = swap_used,
        swap_total = swap_total,
        cached = cached,
        buffers = 0,
    }
end

local function parse_temps()
    -- macOS temp reading is limited without sudo powermetrics
    -- Try osx-cpu-temp if available
    local temp_str = exec_line("osx-cpu-temp 2>/dev/null")
    local temps = {}
    if temp_str then
        local temp_c = tonumber(temp_str:match("([%d%.]+)"))
        if temp_c then
            temps[#temps + 1] = { label = "CPU", temp_c = temp_c }
        end
    end
    return temps
end

local function parse_network()
    local output = exec("netstat -ib")
    if not output then
        return { interface = "none", rx_speed = 0, tx_speed = 0, rx_total = 0, tx_total = 0 }
    end

    local now = os.time()
    local current_stats = {}
    local best_iface = nil
    local best_rx = 0

    for line in output:gmatch("[^\n]+") do
        local iface, rx, tx = line:match("^(%S+)%s+%d+%s+%S+%s+%S+%s+(%d+)%s+%S+%s+(%d+)")
        if iface and iface ~= "lo0" and not iface:match("^lo") then
            rx = tonumber(rx)
            tx = tonumber(tx)
            if rx and tx then
                current_stats[iface] = { rx = rx, tx = tx }
                if rx > best_rx then
                    best_rx = rx
                    best_iface = iface
                end
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
    local mounts = {}
    local handle = io.popen("df -b / /System/Volumes/Data 2>/dev/null")
    if handle then
        local output = handle:read("*a")
        handle:close()
        local first = true
        for line in output:gmatch("[^\n]+") do
            if first then
                first = false
            else
                local _, blocks, used, _, pcent, target =
                    line:match("^(%S+)%s+(%d+)%s+(%d+)%s+(%d+)%s+(%d+)%%%s+(%S+)")
                if target then
                    mounts[#mounts + 1] = {
                        mount = target,
                        total = tonumber(blocks) * 512,
                        used = tonumber(used) * 512,
                        percent = tonumber(pcent),
                    }
                end
            end
        end
    end

    -- iostat for disk I/O
    local io_read = 0
    local io_write = 0
    local device = "disk0"
    local io_output = exec("iostat -d -c 1 2>/dev/null")
    if io_output then
        local last_line = nil
        for line in io_output:gmatch("[^\n]+") do
            if line:match("%d") then last_line = line end
        end
        if last_line then
            local kb_read, kb_write = last_line:match("([%d%.]+)%s+([%d%.]+)$")
            if kb_read then
                io_read = tonumber(kb_read) * 1024
                io_write = tonumber(kb_write) * 1024
            end
        end
    end

    return {
        mounts = mounts,
        io_read = io_read,
        io_write = io_write,
        device = device,
    }
end

local function parse_battery()
    local output = exec("pmset -g batt")
    if not output then
        return { percent = -1, status = "No battery", power = 0, voltage = 0, technology = "Li-ion" }
    end

    local percent = tonumber(output:match("(%d+)%%")) or -1
    local status = "Unknown"
    if output:match("charging") then
        status = "Charging"
    elseif output:match("discharging") then
        status = "Discharging"
    elseif output:match("charged") then
        status = "Full"
    end

    return {
        percent = percent,
        status = status,
        power = 0,
        voltage = 0,
        technology = "Li-ion",
    }
end

local function parse_system()
    local load = { 0, 0, 0 }
    local load_str = exec_line("sysctl -n vm.loadavg")
    if load_str then
        local l1, l5, l15 = load_str:match("([%d%.]+)%s+([%d%.]+)%s+([%d%.]+)")
        if l1 then
            load = { tonumber(l1), tonumber(l5), tonumber(l15) }
        end
    end

    -- Uptime
    local uptime = 0
    local boot_str = exec_line("sysctl -n kern.boottime")
    if boot_str then
        local boot_sec = tonumber(boot_str:match("sec%s*=%s*(%d+)"))
        if boot_sec then
            uptime = os.time() - boot_sec
        end
    end

    -- Process count
    local procs = 0
    local ps_output = exec("ps -e | wc -l")
    if ps_output then
        procs = math.max(0, (tonumber(ps_output:match("%d+")) or 1) - 1)
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
