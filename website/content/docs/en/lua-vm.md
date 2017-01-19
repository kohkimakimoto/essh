+++
title = "Lua VM | Documentation"
type = "docs"
category = "docs"
lang = "en"
basename = "lua-vm.html"
+++

# Lua VM

Essh uses [GopherLua](https://github.com/yuin/gopher-lua) as a Lua VM to load configuration files written in Lua.

## Built-in Functions

As you already seen `host` and `task` functions, Essh core features consist of built-in functions. All the functions that Essh provides are listed below.

* `host`: Defines a host. See [Hosts](/docs/en/hosts.html).

* `task`: Defines a task. See [Tasks](/docs/en/tasks.html).

* `driver`: Defines a driver. See [Drivers](/docs/en/drivers.html).

* `job`: Defines a job. See [Jobs](/docs/en/jobs.html).

* `import`: Imports a module. See [Modules](/docs/en/modules.html).

## Built-in Libraries

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
* `question`: [kohkimakimoto/gluaquestion](https://github.com/kohkimakimoto/gluaquestion).
* `template`: [kohkimakimoto/gluatemplate](https://github.com/kohkimakimoto/gluatemplate).
* `env`: [kohkimakimoto/gluaenv](https://github.com/kohkimakimoto/gluaenv).
* `http`: [cjoudrey/gluahttp](https://github.com/cjoudrey/gluahttp).
* `re`: [yuin/gluare](https://github.com/yuin/gluare)
* `sh`: [otm/gluash](https://github.com/otm/gluash)

## Predefined Variables

Essh provides predefined variables. In the recent version of Essh, there is one predefined variable: `essh`.

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

* `select_hosts` (function): Gets defined hosts. It is useful for overriding host config or setting default values. For example, if you want to set a default ssh_config: `ForwardAgent = yes`, you can achieve it the below code:

    ~~~lua
    -- ~/.essh/config_override.lua
    for _, h in pairs(essh.select_hosts():get()) do
        if h.ForwardAgent == nil then
            h.ForwardAgent = "yes"
        end
    end
    ~~~

    Above example sets the default value to the all hosts. If you want to set the value to some hosts you selected, You can also use the below code:

    ~~~lua
    -- ~/.essh/config_override.lua
    -- Getting only the hosts that has `web` tag or name of the hosts is `web`.
    for _, h in pairs(essh.select_hosts("web"):get()) do
        if h.ForwardAgent == nil then
            h.ForwardAgent = "yes"
        end
    end

    -- Using a table, Getting the hosts both `web` or `db`
    for _, h in pairs(essh.select_hosts({"web", "db"}):get()) do
        if h.ForwardAgent == nil then
            h.ForwardAgent = "yes"
        end
    end

    -- You can set a filter.
    -- Getting only the `web` hosts filtered by `production`.
    for _, h in pairs(essh.select_hosts("web"):filter("production"):get()) do
        if h.ForwardAgent == nil then
            h.ForwardAgent = "yes"
        end
    end

    -- Getting only the first one host using `first` method.
    local h = essh.select_hosts("web"):first()
    if h.ForwardAgent == nil then
        h.ForwardAgent = "yes"
    end
    ~~~

* `get_job` (function): Get a job. See the following examples:

    ~~~lua
    for _, h in pairs(essh.get_job("myjob"):select_hosts():get()) do
        if h.ForwardAgent == nil then
            h.ForwardAgent = "yes"
        end
    end
    ~~~

* `module` (table): A table only lives in the imported module. You can use it as a module scope variable in a module. It has some predefined values.

    * `path` (string): A path to the module.
    
    * `import_path` (string): A path that was used by import function's argument.

* `host` (function): An alias of `host` function.

* `task` (function): An alias of `task` function.

* `driver` (function): An alias of `driver` function.

* `job` (function): An alias of `job` function.

* `import` (function): An alias of `import` function.

* `debug` (function): Output a debug message. The debug message is outputed when you run Essh with `--debug` option.

    ~~~~lua
    essh.debug("foo")
    ~~~~

