local M = {}
M.__index = M

function M.new(size)
    return setmetatable({
        _size = size or 60,
        _data = {},
    }, M)
end

function M:push(value)
    local data = {}
    local start = 1
    if #self._data >= self._size then
        start = 2
    end
    for i = start, #self._data do
        data[#data + 1] = self._data[i]
    end
    data[#data + 1] = value
    return setmetatable({
        _size = self._size,
        _data = data,
    }, M)
end

function M:get_all()
    local result = {}
    for i, v in ipairs(self._data) do
        result[i] = v
    end
    return result
end

return M
