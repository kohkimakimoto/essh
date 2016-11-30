# Essh Docker Module

This module provides docker driver engine for Essh.
It allows you to run an Essh task in a docker container.

## Usage

```lua
local docker = import "github.com/kohkimakimoto/essh/modules/docker"

driver "docker-centos7" {
    engine = docker.driver,
    image = "centos:centos7",
    privileged = true,
    remove_terminated_container = true,
}

task "example" {
    backend = "remote",
    targets = "default",
    driver = "docker-centos7",
    script = {
        "echo hello",
    },
}
```
