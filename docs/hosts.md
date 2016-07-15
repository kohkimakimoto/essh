# Hosts

Hosts in Essh are managed SSH servers. Using hosts configuration, Essh dynamically generates SSH config, provides hook functions, classify the hosts by tags.

## Example

```lua
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
    hooks = {
        before_connect = "",
        after_connect = "",
        after_disconnect = "",
    },
}
```

Host is composed of two different kinds of properties. **SSH Config Properties** and **Essh Config Properties**.

### SSH Config Properties

SSH config properties require that the first character is upper case.
For instance `HostName` and `Port`. They are used to generate **ssh_config**. You can use all ssh options to these properties. see ssh_config(5).

### Essh Config Properties

Essh config properties require that the first character is lower case.
They are used for special purpose of Essh functions, not ssh_config.

All the properties of this type are listed below.

* `description` (string): Description is a description of the host.

* `hidden` (boolean): If you set it true, zsh completion doesn't show the host.

* `private` (boolean): If you set it true, This host can be only used to the tasks. See also [private hosts](#private-hosts)

* `hooks` (table): Hooks is a table that defines `before_connect`(string, table or function), `after_connect`(string, table or function) and `after_disconnect`(string, table or function).

    ```lua
    hooks = {
        -- WIP...
    }
    ```

* `tags` (array table): Tags classifys hosts.

    ```lua
    tags = {
        "web",
        "production",
    }
    ```
* `props` (table): Props sets environment variables `ESSH_HOST_PROPS_{KEY}` when the host is used in tasks. The table key is modified to upper cased.

    ```lua
    props = {
        foo = "bar",
    }

    -- ESSH_HOST_PROPS_FOO=bar
    ```

## Private Hosts

Private hosts can be only used to the tasks. The hosts of this type can't be used by ssh login, and does not suggest by zsh-completion. also they can't be used with `--exec` option.

You can use `private_host` function as an alias to define a private host. See the below example:

```lua
private_host "example" {
    HostName = "192.168.0.11",
    Port = "22",
    User = "kohkimakimoto",
}
```

This is same the following:

```lua
host "example" {
    HostName = "192.168.0.11",
    Port = "22",
    User = "kohkimakimoto",
    private = true,
}
```

## Constraints

There are constraints about scopes and registries.

* Each public hosts must be unique. (You can NOT define public hosts by the same name in the local and global registry.)
* Any hosts must be unique in a same registry. (You can NOT define hosts by the same name in the same registry.)
* Hosts used by task must be defined in a same registry. (Tasks can refer to only hosts defined in the same registry.)
* Private hosts is only used by tasks.
* There can be duplicated hosts in the entire registries. (You can define private hosts even if you define same name public hosts.)
