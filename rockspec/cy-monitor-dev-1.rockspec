package = "cy-monitor"
version = "dev-1"
source = {
    url = "git+https://github.com/vexxuh/cy-monitor.git"
}
description = {
    summary = "Terminal-based hardware monitor",
    detailed = [[
        A cross-platform terminal hardware monitor using Lua + lcurses.
        Features 2x3 dashboard grid with CPU, memory, temperatures, network,
        disk, and system/battery panels. Reads from /proc + /sys on Linux,
        shell commands on macOS. No network dependencies.
    ]],
    license = "MIT"
}
dependencies = {
    "lua >= 5.3",
    "lcurses >= 9.0",
}
build = {
    type = "builtin",
    modules = {},
    install = {
        bin = {
            ["cy-monitor"] = "src/main.lua"
        }
    }
}
