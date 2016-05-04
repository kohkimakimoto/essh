# Essh Docker Module

This module provides docker driver engine for Essh driver.
It allows you to run a Essh task in a docker container.

## Usage

```lua
local docker = essh.require "github.com/kohkimakimoto/essh/modules/docker"

task "example" {
    description = "example",
    configure = function()
        driver "docker" {
            engine = docker.driver,
            image = "centos:centos6",
        }
    end,
    driver = "docker",
    script = {
        [=[
            cat /etc/redhat-release
        ]=],
    }
}
```

Building docker image before running if it doesn't exist.

```lua
driver "docker" {
    engine = docker.driver,
    image = "my-custom-image",
    build = {
        -- using current directory Dockerfile.
        url = ".",
    }
}
```
