#  Configuration Files

Essh configuration is written in [Lua](https://www.lua.org/). In the configuration files, you can also use DSL syntax that is more human-readable. Here is an example:

```lua
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
```

Essh loads configuration files from several different places. All configuration files are not required. Essh loads these if they exist.

* At first, loads `/etc/essh/config.lua`.
* At second, loads `~/.essh/config.lua`.
* At third, loads `esshconfig.lua` in the current directory or loads a file specified by `--config` command line option.
* At fourth, loads `esshconfig_override.lua` in the current directory or loads a file the name of which with end in `_override`, that specified by `--config` command line option. ex) If you use `--config=myconfig.lua`, valid file name is `myconfig_override.lua`.
* At fifth, loads `~/.essh/config_override.lua`.
* At seventh, load `/etc/essh/config_override.lua`

## Registry
