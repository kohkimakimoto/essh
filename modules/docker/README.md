# Essh Docker Module


## Usage

```lua
local docker = essh.require "github.com/kohkimakimoto/essh//modules/docker"

driver "docker" {
    engine = docker.driver,
    image = "centos:centos6",
}

task "example" {
    driver = "docker",
    description = "example",
    script = {
        [=[
            cat /etc/redhat-release
        ]=],
    }
}
```
