package.path = package.path .. ";./?.lua;./?/init.lua"

describe("format", function()
    local format = require("src.utils.format")

    describe("bytes", function()
        it("formats gigabytes", function()
            assert.are.equal("1.5G", format.bytes(1536 * 1024 * 1024))
        end)

        it("formats megabytes", function()
            assert.are.equal("512.0M", format.bytes(512 * 1024 * 1024))
        end)

        it("formats kilobytes", function()
            assert.are.equal("100.0K", format.bytes(100 * 1024))
        end)

        it("formats bytes", function()
            assert.are.equal("500B", format.bytes(500))
        end)

        it("handles zero", function()
            assert.are.equal("0B", format.bytes(0))
        end)

        it("handles nil", function()
            assert.are.equal("--", format.bytes(nil))
        end)

        it("clamps negative to zero", function()
            assert.are.equal("0B", format.bytes(-100))
        end)
    end)

    describe("bytes_speed", function()
        it("appends /s", function()
            assert.are.equal("1.0M/s", format.bytes_speed(1024 * 1024))
        end)

        it("handles nil", function()
            assert.are.equal("-- /s", format.bytes_speed(nil))
        end)
    end)

    describe("temp", function()
        it("formats celsius", function()
            assert.are.equal("46°C", format.temp(46))
        end)

        it("formats fahrenheit", function()
            assert.are.equal("115°F", format.temp(46, "F"))
        end)

        it("handles nil", function()
            assert.are.equal("--", format.temp(nil))
        end)
    end)

    describe("duration", function()
        it("formats days hours minutes", function()
            assert.are.equal("1d 2h 30m", format.duration(86400 + 7200 + 1800))
        end)

        it("formats hours and minutes", function()
            assert.are.equal("2h 30m", format.duration(9000))
        end)

        it("formats minutes only", function()
            assert.are.equal("5m", format.duration(300))
        end)

        it("handles zero", function()
            assert.are.equal("0m", format.duration(0))
        end)

        it("handles nil", function()
            assert.are.equal("--", format.duration(nil))
        end)

        it("handles negative", function()
            assert.are.equal("--", format.duration(-10))
        end)
    end)

    describe("freq", function()
        it("formats GHz", function()
            assert.are.equal("4.2 GHz", format.freq(4200))
        end)

        it("formats MHz", function()
            assert.are.equal("800 MHz", format.freq(800))
        end)

        it("handles nil", function()
            assert.are.equal("--", format.freq(nil))
        end)
    end)

    describe("percent", function()
        it("formats percentage", function()
            assert.are.equal("42%", format.percent(42))
        end)

        it("handles nil", function()
            assert.are.equal("--%", format.percent(nil))
        end)
    end)

    describe("pad_right", function()
        it("pads short strings", function()
            assert.are.equal("abc       ", format.pad_right("abc", 10))
        end)

        it("truncates long strings", function()
            assert.are.equal("abcde", format.pad_right("abcdefghij", 5))
        end)
    end)
end)
