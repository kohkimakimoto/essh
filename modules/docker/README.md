# Essh Docker Module

This module provides docker driver engine for Essh driver.
It allows you to run a Essh task in a docker container.

## Usage Example

```lua
local docker = essh.require "github.com/kohkimakimoto/essh/modules/docker"

driver "docker-centos6" {
    engine = docker.driver,
    image = "centos:centos6",
    privileged = true,
}

task "example" {
    backend = "remote",
    targets = "default",
    driver = "docker-centos6",
    script = {
        "echo hello",
    },
}
```

Experimental:

Building docker image before running if it doesn't exist (only local task).


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
