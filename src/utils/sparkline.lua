local M = {}

-- Unicode block characters U+2581 through U+2588
local BLOCKS = {
    "▁", -- U+2581
    "▂", -- U+2582
    "▃", -- U+2583
    "▄", -- U+2584
    "▅", -- U+2585
    "▆", -- U+2586
    "▇", -- U+2587
    "█", -- U+2588
}

function M.render(values, width)
    if not values or #values == 0 then
        return string.rep(" ", width or 20)
    end

    width = width or #values

    -- Resample values to fit width if needed
    local sampled = M.resample(values, width)

    local min_val = math.huge
    local max_val = -math.huge
    for _, v in ipairs(sampled) do
        if v < min_val then min_val = v end
        if v > max_val then max_val = v end
    end

    local range = max_val - min_val
    if range == 0 then
        range = 1
    end

    local result = {}
    for _, v in ipairs(sampled) do
        local normalized = (v - min_val) / range
        local index = math.floor(normalized * 7) + 1
        if index > 8 then index = 8 end
        if index < 1 then index = 1 end
        result[#result + 1] = BLOCKS[index]
    end

    return table.concat(result)
end

function M.resample(values, target_len)
    if #values == target_len then
        return values
    end

    if #values == 0 then
        return {}
    end

    local result = {}
    local ratio = #values / target_len

    for i = 1, target_len do
        local src_idx = (i - 1) * ratio + 1
        local lower = math.floor(src_idx)
        local upper = math.min(lower + 1, #values)
        local frac = src_idx - lower

        if lower == upper or upper > #values then
            result[i] = values[lower]
        else
            result[i] = values[lower] * (1 - frac) + values[upper] * frac
        end
    end

    return result
end

return M
