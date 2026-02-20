local M = {}

function M.bytes(value)
    if not value or type(value) ~= "number" then
        return "--"
    end
    if value < 0 then value = 0 end
    if value >= 1024 * 1024 * 1024 then
        return string.format("%.1fG", value / (1024 * 1024 * 1024))
    elseif value >= 1024 * 1024 then
        return string.format("%.1fM", value / (1024 * 1024))
    elseif value >= 1024 then
        return string.format("%.1fK", value / 1024)
    else
        return string.format("%dB", math.floor(value))
    end
end

function M.bytes_speed(value)
    if not value or type(value) ~= "number" then
        return "-- /s"
    end
    return M.bytes(value) .. "/s"
end

function M.temp(celsius, unit)
    if not celsius or type(celsius) ~= "number" then
        return "--"
    end
    unit = unit or "C"
    if unit == "F" then
        return string.format("%.0f\194\176F", celsius * 9 / 5 + 32)
    end
    return string.format("%.0f\194\176C", celsius)
end

function M.duration(secs)
    if not secs or type(secs) ~= "number" or secs < 0 then
        return "--"
    end
    secs = math.floor(secs)
    local days = math.floor(secs / 86400)
    local hours = math.floor((secs % 86400) / 3600)
    local mins = math.floor((secs % 3600) / 60)

    if days > 0 then
        return string.format("%dd %dh %dm", days, hours, mins)
    elseif hours > 0 then
        return string.format("%dh %dm", hours, mins)
    else
        return string.format("%dm", mins)
    end
end

function M.freq(mhz)
    if not mhz or type(mhz) ~= "number" then
        return "--"
    end
    if mhz >= 1000 then
        return string.format("%.1f GHz", mhz / 1000)
    end
    return string.format("%.0f MHz", mhz)
end

function M.percent(value)
    if not value or type(value) ~= "number" then
        return "--%"
    end
    return string.format("%.0f%%", value)
end

function M.pad_right(str, width)
    str = str or ""
    local len = utf8 and utf8.len(str) or #str
    if len >= width then
        return str:sub(1, width)
    end
    return str .. string.rep(" ", width - len)
end

function M.truncate(str, max_width)
    if not str then return "" end
    local len = utf8 and utf8.len(str) or #str
    if len <= max_width then
        return str
    end
    return str:sub(1, max_width - 1) .. "…"
end

return M
