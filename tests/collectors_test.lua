package.path = package.path .. ";./?.lua;./?/init.lua"

describe("collectors", function()
    local collectors = require("src.collectors")

    it("detects platform", function()
        assert.is_string(collectors.platform)
        assert.is_not.are.equal("", collectors.platform)
    end)

    it("collects all data", function()
        local data = collectors.collect_all()
        assert.is_table(data)

        -- CPU
        assert.is_table(data.cpu)
        assert.is_number(data.cpu.usage)
        assert.is_table(data.cpu.cores)
        assert.is_number(data.cpu.freq)
        assert.is_table(data.cpu.load)
        assert.are.equal(3, #data.cpu.load)

        -- Memory
        assert.is_table(data.memory)
        assert.is_number(data.memory.ram_used)
        assert.is_number(data.memory.ram_total)
        assert.is_number(data.memory.swap_used)
        assert.is_number(data.memory.swap_total)
        assert.is_number(data.memory.cached)
        assert.is_number(data.memory.buffers)

        -- Temps
        assert.is_table(data.temps)

        -- Network
        assert.is_table(data.network)
        assert.is_string(data.network.interface)
        assert.is_number(data.network.rx_speed)
        assert.is_number(data.network.tx_speed)
        assert.is_number(data.network.rx_total)
        assert.is_number(data.network.tx_total)

        -- Disk
        assert.is_table(data.disk)
        assert.is_table(data.disk.mounts)
        assert.is_number(data.disk.io_read)
        assert.is_number(data.disk.io_write)
        assert.is_string(data.disk.device)

        -- Battery
        assert.is_table(data.battery)
        assert.is_number(data.battery.percent)
        assert.is_string(data.battery.status)

        -- System
        assert.is_table(data.system)
        assert.is_number(data.system.uptime)
        assert.is_number(data.system.procs)
        assert.is_table(data.system.load)
    end)

    it("returns sensible memory values", function()
        local data = collectors.collect_all()
        assert.is_true(data.memory.ram_total > 0, "ram_total should be > 0")
        assert.is_true(data.memory.ram_used >= 0, "ram_used should be >= 0")
        assert.is_true(data.memory.ram_used <= data.memory.ram_total,
            "ram_used should be <= ram_total")
    end)

    it("returns sensible system values", function()
        local data = collectors.collect_all()
        assert.is_true(data.system.uptime > 0, "uptime should be > 0")
        assert.is_true(data.system.procs > 0, "procs should be > 0")
    end)
end)
