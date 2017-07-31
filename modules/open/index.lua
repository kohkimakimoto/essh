local name = essh.module.var.name or error("require 'name'")
local url = essh.module.var.url or error("require 'url'")
local description = essh.module.var.description

task(name, {
    backend = "local",
    prefix = false,
    privileged = false,
    script = "open " .. url,
    description = description,
})

