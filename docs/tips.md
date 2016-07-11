
vagrant connections.

```lua
local sh = require "glua.sh"

local vagrant_hosts = {}
local current_hostname = nil

for line in sh.vagrant("ssh-config"):lines() do

    local _, hostname = string.match(line, "^(Host )(.-)$")

    if hostname ~= nil then
        current_hostname = hostname
        vagrant_hosts[hostname] = {
            description = hostname .. " (vagrant vm)",
            tags = {"vagrant"},
        }
    else
        local _, config_key, config_value = string.match(line, "^(  )(.-) (.-)$")
        if config_key ~= nil then
            vagrant_hosts[current_hostname][config_key] = config_value
        end
    end

end

for name, config in pairs(vagrant_hosts) do
    host(name, config)
end
```
