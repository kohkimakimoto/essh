local sh = require "glua.sh"
local fs = require "glua.fs"

local vagrant = {}
local cache_filename = "vagrant-ssh-config.cache"

vagrant.hosts = function()
    local vagrant_hosts = {}
    local registry = essh.registry()
    local tmp_dir = registry:tmp_dir()
    local cache_file = tmp_dir .. "/" .. cache_filename

    if fs.exists(cache_file) == false then
        sh.vagrant("ssh-config"):combinedOutput(cache_file)
    end

    local current_hostname = nil

    for line in sh.cat(cache_file):lines() do
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
