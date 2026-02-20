package.path = package.path .. ";./?.lua;./?/init.lua"

describe("history ring buffer", function()
    local history = require("src.collectors.history")

    it("creates empty buffer", function()
        local h = history.new(5)
        assert.are.same({}, h:get_all())
    end)

    it("pushes values immutably", function()
        local h1 = history.new(3)
        local h2 = h1:push(10)
        local h3 = h2:push(20)

        assert.are.same({}, h1:get_all())
        assert.are.same({10}, h2:get_all())
        assert.are.same({10, 20}, h3:get_all())
    end)

    it("respects max size", function()
        local h = history.new(3)
        h = h:push(1)
        h = h:push(2)
        h = h:push(3)
        h = h:push(4)
        assert.are.same({2, 3, 4}, h:get_all())
    end)

    it("get_all returns a copy", function()
        local h = history.new(5)
        h = h:push(10)
        local all = h:get_all()
        all[1] = 999
        assert.are.equal(10, h:get_all()[1])
    end)
end)
