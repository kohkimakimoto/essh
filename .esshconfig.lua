require "drivers"

task "example-remote-lua" {
    driver = "essh-lua",
    description = "Run a remote lua script",
    script_file = "https://raw.githubusercontent.com/sevir/essh/refs/heads/master/test/test.lua"
}

task "example-http-lua" {
    driver = "essh-lua",
    description = "Run a remote lua script using http module",
    script = [[
        local http = require("http")

        local res = http.get("https://raw.githubusercontent.com/sevir/essh/refs/heads/master/test/test.lua")
        local code = loadstring(res.body)
        code()
    ]]
}

host "host-example" {
    Hostname = "example.com",
    description = "An example host",
}