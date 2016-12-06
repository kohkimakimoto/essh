+++
title = "Hosts | Documentation"
type = "docs"
category = "docs"
lang = "en"
basename = "hosts.html"
+++

# Hosts

Hosts in Essh are SSH servers that you manage. Using hosts configuration, Essh dynamically generates SSH config, provides hook functions, classify the hosts by tags.

Example:
~~~lua
host "web01.localhost" {
    -- SSH config properties.
    HostName = "192.168.0.11",
    Port = "22",
    User = "kohkimakimoto",

    -- Essh config properties.
    description = "web01 development server",
    hidden = false,
    private = false,
    props = {},
    tags = {},
    hooks_before_connect = {},
    hooks_after_connect = {},
    hooks_after_disconnect = {},
}
~~~

Host is composed of two different kinds of properties. **SSH Config Properties** and **Essh Config Properties**.

## SSH Config Properties

SSH config properties require that the first character is upper case.
For instance `HostName` and `Port`. They are used to generate **ssh_config**. You can use all ssh options to these properties. see ssh_config(5).

## Essh Config Properties

Essh config properties require that the first character is lower case.
They are used for special purpose of Essh functions, not ssh_config.

All the properties of this type are listed below.

* `description` (string): Description of the host. This is used for displaying hosts list and zsh completion.

* `hidden` (boolean): If you set it true, zsh completion doesn't show the host.

* `hooks_before_connect` (table): Hooks that fire before connect. This hook runs on local. The hook is defined as a Lua table. This table can have mulitple functions or strings. See the example:

    ~~~lua
    hooks_before_connect = {
        -- function
        function()
            print("foo")
        end,

        -- string (commands)
        "echo bar",

        -- If the function returns a string, Essh run the string as a command.
        function()
            return "echo foobar"
        end,
    }
    ~~~

    All hooks (includes `hooks_after_connect`, `hooks_after_disconnect`) implemented in Lua function runs on local.

    All hooks (includes `hooks_after_connect`, `hooks_after_disconnect`) only fire when your simply login with ssh. Hooks don't fire in tasks and with `--exec` option.

* `hooks_after_connect` (table): Hooks that fire after connect. This hook runs on remote.

* `hooks_after_disconnect` (table): Hooks that fire after disconnect. This hook runs on local.

* `tags` (array table): Tags classifies hosts.

    ~~~lua
    tags = {
        "web",
        "production",
    }
    ~~~

    Tags mustn't be duplicated with any host names.

* `props` (table): Props sets environment variables `ESSH_HOST_PROPS_{KEY}` when the host is used in tasks. The table key is modified to upper cased.

    ~~~lua
    props = {
        foo = "bar",
    }

    -- ESSH_HOST_PROPS_FOO=bar
    ~~~
