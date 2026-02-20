package.path = package.path .. ";./?.lua;./?/init.lua"

describe("sparkline", function()
    local sparkline = require("src.utils.sparkline")

    describe("render", function()
        it("renders ascending values", function()
            local result = sparkline.render({1, 2, 3, 4, 5, 6, 7, 8})
            assert.is_string(result)
            assert.are.equal(8, utf8.len(result))
        end)

        it("renders constant values as same block", function()
            local result = sparkline.render({5, 5, 5, 5})
            -- All same value => all same block char
            local first_char = result:sub(1, utf8.offset(result, 2) - 1)
            for _, c in utf8.codes(result) do
                assert.are.equal(utf8.codepoint(first_char), c)
            end
        end)

        it("handles empty values", function()
            local result = sparkline.render({}, 10)
            assert.are.equal(10, #result)
        end)

        it("handles nil values", function()
            local result = sparkline.render(nil, 10)
            assert.are.equal(10, #result)
        end)

        it("renders min at lowest block, max at highest", function()
            local result = sparkline.render({0, 100})
            local chars = {}
            for _, c in utf8.codes(result) do
                chars[#chars + 1] = c
            end
            -- First should be lowest block (U+2581 = 9601)
            assert.are.equal(9601, chars[1])
            -- Last should be highest block (U+2588 = 9608)
            assert.are.equal(9608, chars[2])
        end)

        it("resamples to fit width", function()
            local result = sparkline.render({1, 2, 3, 4, 5, 6, 7, 8, 9, 10}, 5)
            assert.are.equal(5, utf8.len(result))
        end)
    end)

    describe("resample", function()
        it("returns same values when lengths match", function()
            local result = sparkline.resample({1, 2, 3}, 3)
            assert.are.same({1, 2, 3}, result)
        end)

        it("downsamples correctly", function()
            local result = sparkline.resample({1, 2, 3, 4}, 2)
            assert.are.equal(2, #result)
        end)

        it("upsamples with interpolation", function()
            local result = sparkline.resample({0, 10}, 3)
            assert.are.equal(3, #result)
            assert.are.equal(0, result[1])
        end)

        it("handles empty input", function()
            local result = sparkline.resample({}, 5)
            assert.are.same({}, result)
        end)
    end)
end)
