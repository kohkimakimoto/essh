local example = {}

example.hello = function()
    print("hello world!")
end

example.tasks = {
    hello = {
        description = "example hello",
        script = "echo 'hell oworld!'",
    },
    uptime = {
        description = "example uptime",
        script = "uptime",
    },
}

return example
