local sh = require "glua.sh"
local fs = require "glua.fs"

local vagrant = {}
vagrant.cachne = ".vagrant-ssh-config.cache"

vagrant.hosts = function()
    local vagrant_hosts = {}

    if fs.exists(vagrant.cachne) == false then
        sh.vagrant("ssh-config"):combinedOutput(vagrant.cachne)
    end

    local current_hostname = nil

    for line in sh.cat(vagrant.cachne):lines() do
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

    return vagrant_hosts
end

return vagrant
