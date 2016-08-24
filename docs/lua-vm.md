# Lua VM

Essh uses [GopherLua](https://github.com/yuin/gopher-lua) as a Lua VM to load configuration files written in Lua.

## Libraries

Essh provides built-in Lua libraries that you can use in your configuration files.
For instance, if you want to use `glua.json` library, you should use Lua's `require` function like below.

```lua
local json = require("glua.json")

local jsontext = json.encode({aaa = "bbb", ccc = "ddd"})
print(jsontext)
```

The following are the built-in libraries that are included in Essh.

* `json`: Json encoder/decoder. See [layeh/gopher-json](https://github.com/layeh/gopher-json).
* `fs`: Filesystem utility. See [kohkimakimoto/gluafs](https://github.com/kohkimakimoto/gluafs).
* `yaml`: Yaml parser. See [kohkimakimoto/gluayaml](https://github.com/kohkimakimoto/gluayaml).
* `template`: Text template. See [kohkimakimoto/gluatemplate](https://github.com/kohkimakimoto/gluatemplate).
* `env`: Utility package for manipulating environment variables. See [kohkimakimoto/gluaenv](https://github.com/kohkimakimoto/gluaenv).
* `http`: Http module. See [cjoudrey/gluahttp](https://github.com/cjoudrey/gluahttp).
* `re`: Regular expressions for the GopherLua. See [yuin/gluare](https://github.com/yuin/gluare)
* `sh`:A shell module for gopher-lua. See [otm/gluash](https://github.com/otm/gluash)

Of course, You can also use another Lua libraries by using `require`. See the Lua's [manual](http://www.lua.org/manual/5.1/manual.html#pdf-require).

## Predefined Variables

Essh provides predefined variables. In the recent version of Essh, there is one predefined varilable: `essh`.

`essh` is a table that has some functions and variables. see below

* `ssh_config` (string): ssh_config is ssh_config file path. At default, it is a temporary file that is generated automatically when you run Essh. You can overwrite this value for generating ssh_config to a static destination. If you use a gateway host that is a server between your client computer and a target server, you may use this variable to specify `ProxyCommand`. See below example:

    ```lua
    --
    -- network environment.
    -- [your-computer] -- [getway-server1] -- [web-server]
    --

    host "web-server" {
        HostName = "192.168.0.1",
        ProxyCommand = "ssh -q -F " .. essh.ssh_config .. " -W %h:%p getway-server1",
    }
    ```

* `debug` (function): debug is a function to output the debug message. The debug message outputs only when you run Essh with `--debug` option.

  ```lua
  essh.debug("this is a debug message!")
  ```

* `require` (function): require is a function to load Essh module. see the [Module](#module) section.

  ```lua
  local bash = essh.require "github.com/kohkimakimoto/essh/modules/bash"
  ```

* `hosts` (function): hosts is a function to get defined hosts.

* `registry` (function):
