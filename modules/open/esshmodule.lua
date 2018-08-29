local name = essh.module.params.name or error("require 'name'")
local url = essh.module.params.url or error("require 'url'")
local description = essh.module.params.description

task(name, {
    backend = "local",
    prefix = false,
    privileged = false,
    script = "open " .. url,
    description = description,
})
