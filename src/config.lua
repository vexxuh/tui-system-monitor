local M = {}

local CONFIG_DIR = os.getenv("HOME") .. "/.config/cy-monitor"
local CONFIG_FILE = CONFIG_DIR .. "/config.lua"

local function deep_copy(tbl)
    if type(tbl) ~= "table" then
        return tbl
    end
    local copy = {}
    for k, v in pairs(tbl) do
        copy[k] = deep_copy(v)
    end
    return copy
end

local function merge(base, overrides)
    local result = deep_copy(base)
    for k, v in pairs(overrides) do
        if type(v) == "table" and type(result[k]) == "table" then
            result[k] = merge(result[k], v)
        else
            result[k] = deep_copy(v)
        end
    end
    return result
end

local function load_sandboxed(path)
    local safe_env = {
        math = math,
        string = string,
        table = table,
        tonumber = tonumber,
        tostring = tostring,
        type = type,
        pairs = pairs,
        ipairs = ipairs,
    }
    local fn, err = loadfile(path, "t", safe_env)
    if not fn then return nil end
    local ok, result = pcall(fn)
    if ok and type(result) == "table" then
        return result
    end
    return nil
end

function M.load()
    local defaults = require("config.default")
    local config = deep_copy(defaults)

    local user_config = load_sandboxed(CONFIG_FILE)
    if user_config then
        config = merge(config, user_config)
    end

    return config
end

return M
