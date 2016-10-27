# Essh Docker Module

This module provides docker driver engine for Essh driver.
It allows you to run a Essh task in a docker container.

## Usage Example

```lua
local docker = import "github.com/kohkimakimoto/essh/modules/docker"

driver "docker-centos6" {
    engine = docker.driver,
    image = "centos:centos6",
    privileged = true,
    remove_terminated_containers = true,
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
