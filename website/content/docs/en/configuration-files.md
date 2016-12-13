+++
title = "configuration Files | Documentation"
type = "docs"
category = "docs"
lang = "en"
basename = "configuration-files.html"
+++

# Configuration Files

Essh configuration is written in [Lua](https://www.lua.org/). In the configuration files, you can use DSL syntax that is more human-readable.

## Example

~~~lua
host "web01.localhost" {
    HostName = "192.168.0.11",
    Port = "22",
    User = "kohkimakimoto",
    description = "web01 development server",
    tags = {
        "web",
    },
}

host "web02.localhost" {
    HostName = "192.168.0.12",
    Port = "22",
    User = "kohkimakimoto",
    description = "web02 development server",
    tags = {
        "web",
    },
}

task "uptime" {
    backend = "remote",
    targets = "web",
    script = "uptime",
}
~~~

## Another Syntax

The above example of configuration is written in Lua DSL style. You can also use plain Lua functions styles. The following examples are valid config code.

~~~lua
host("web01.localhost", {
    HostName = "192.168.0.11",
    Port = "22",
    User = "kohkimakimoto",
    description = "web01 development server",
    tags = {
        "web",
    },
})
~~~

or

~~~lua
local web01 = host "web01.localhost"
web01.HostName = "192.168.0.11"
web01.Port = "22"
web01.User = "kohkimakimoto"
web01.description = "web01 development server"
web01.tags = {
    "web",
}
~~~

## Evaluating Orders

Essh loads configuration files from several different places. All configuration files are not required. Essh loads these if they exist. Configuration are applied in the following order:

1. Loads `~/.essh/config.lua`.
1. Loads `esshconfig.lua` in the current directory.
1. Loads `esshconfig_override.lua` in the current directory.
1. Loads `~/.essh/config_override.lua`.

If you use `--config` command line option, Essh loads the file and the per-user configuration files (`~/.essh/config.lua` and `~/.essh/config_override.lua`) will be ignored. In this case, configuration are applied in the following order:

1. Loads a file specified by `--config` command line option.
1. Loads a file the name of which end in `_override`, that specified by `--config` command line option. ex) If you use `--config=myconfig.lua`, valid file name is `myconfig_override.lua`.

## Lua

Essh provides built-in Lua libraries that can be used in the configuration files.

Please see [Lua VM](lua-vm.html).
