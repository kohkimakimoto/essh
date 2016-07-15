# Essh vagrant module

Loading vagrant hosts as Essh hosts.

## Usage

```lua
local vagrant = essh.require "github.com/kohkimakimoto/essh/modules/vagrant"

vagrant.load_hosts()
```

override config:

```
local vagrant = essh.require "github.com/kohkimakimoto/essh/modules/vagrant"

vagrant.load_hosts({
    tags = {"vagrant", "local_dev"},
})
```
