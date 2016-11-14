+++
title = "Lua VM | Documentation"
type = "docs"
category = "docs"
lang = "en"
basename = "lua-vm.html"
+++

# Lua VM

Essh uses [GopherLua](https://github.com/yuin/gopher-lua) as a Lua VM to load configuration files written in Lua.

## Libraries

Essh provides built-in Lua libraries that you can use in your configuration files.
For instance, if you want to use `json` library, you should use Lua's `require` function like below.

~~~lua
local json = require("json")

local jsontext = json.encode({aaa = "bbb", ccc = "ddd"})
print(jsontext)
~~~

The following are the built-in libraries that are included in Essh.

* `json`: [layeh/gopher-json](https://github.com/layeh/gopher-json).
* `fs`: [kohkimakimoto/gluafs](https://github.com/kohkimakimoto/gluafs).
* `yaml`: [kohkimakimoto/gluayaml](https://github.com/kohkimakimoto/gluayaml).
* `question`: [kohkimakimoto/question](https://github.com/kohkimakimoto/gluaquestion).
* `template`: [kohkimakimoto/gluatemplate](https://github.com/kohkimakimoto/gluatemplate).
* `env`: [kohkimakimoto/gluaenv](https://github.com/kohkimakimoto/gluaenv).
* `http`: [cjoudrey/gluahttp](https://github.com/cjoudrey/gluahttp).
* `re`: [yuin/gluare](https://github.com/yuin/gluare)
* `sh`:[otm/gluash](https://github.com/otm/gluash)

## Predefined Variables

Essh provides predefined variables. In the recent version of Essh, there is one predefined varilable: `essh`.

`essh` is a table that has some functions and variables. see below

* `ssh_config` (string): ssh_config is ssh_config file path. At default, it is a temporary file that is generated automatically when you run Essh. You can overwrite this value for generating ssh_config to a static destination. If you use a gateway host that is a server between your client computer and a target server, you may use this variable to specify `ProxyCommand`. See below example:

    ~~~lua
    --
    -- network environment.
    -- [your-computer] -- [getway-server1] -- [web-server]
    --

    host "web-server" {
        HostName = "192.168.0.1",
        ProxyCommand = "ssh -q -F " .. essh.ssh_config .. " -W %h:%p getway-server1",
    }
    ~~~
