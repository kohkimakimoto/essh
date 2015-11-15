# ZSSH

Extended ssh command wrapper.

* This is a single CLI app.
* Simply wraps `ssh` command. You can use it in the same way as `ssh`.
* Automatically generates `~/.ssh/config` from `~/.ssh/zssh.lua`. You can write SSH configuration in Lua programming language.
* Supports zsh completion.
* Provides some hook functions.

## Installation

Run `go get` command.

```
go get github.com/kohkimakimoto/zssh/cmd/zssh
```

## Usage

At first, you should copy your `~/.ssh/config` to `~/.ssh/config.backup` to keep a backup.
ZSSH override `~/.ssh/config` automatically when it runs.

Create and edit `~/.ssh/zssh.lua`.

```lua
Host "web01.localhost" {
    ForwardAgent = "yes",
    HostName = "192.168.0.11",
    Port = "22",
    User = "kohkimakimoto",
    description = "my web01 server",
}

Host "web02.localhost" {
    ForwardAgent = "yes",
    HostName = "192.168.0.12",
    Port = "22",
    User = "kohkimakimoto",
    description = "my web02 server",
}
```

This configuration generates the below ssh config when you run `zssh`.

```
Host web01.localhost
    ForwardAgent yes
    HostName 192.168.0.11
    Port 22
    User kohkimakimoto

Host web02.localhost
    ForwardAgent yes
    HostName 192.168.0.12
    Port 22
    User kohkimakimoto
```

### Zsh Completion

If you want to use zsh completion, add the following code in your `~/.zshrc`

```
eval "$(zssh --zsh-completion)"
```

You will get completion about hosts.

```
$ zssh [TAB]
web01.localhost          -- my web01 server
web02.localhost          -- my web02 server
```

You can hide a host using `hidden` property. If you set it true, zsh completion doesn't show the host.

```lua
Host "web01.localhost" {
    ForwardAgent = "yes",
    HostName = "192.168.0.11",
    Port = "22",
    User = "kohkimakimoto",
    description = "my web01 server",

    hidden = true,
}
```

### Hooks

You can add hook functions `before` and `after` in a host configuration.

```lua
Host "web01.localhost" {
    HostName = "192.168.0.11",
    Port = "22",
    User = "kohkimakimoto",
    ForwardAgent = "yes",
    description = "my web01 server",
    hooks = {
        before = function()
            -- your code ...
        end,
        after = function()
            -- your code ...
        end,
    }
}
```

`before` hook fires before you connect server via SSH. `after` hook fires after you disconnect SSH connection.

## Author

Kohki Makimoto <kohki.makimoto@gmail.com>

## License

MIT license.
