# vagrant

## Usage

```lua
local vagrant = essh.require "github.com/kohkimakimoto/essh/modules/vagrant"

for name, config in pairs(vagrant.hosts()) do
    host(name, config)
end
```
