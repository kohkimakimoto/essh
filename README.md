# ZSSH

Extended ssh command wrapper.

* This is a single CLI app.
* Simply wraps `ssh` command. You can use it in the same way as `ssh`.
* Automatically generates `~/.ssh/config` from `~/.ssh/zssh.lua`
* Support zsh completion.

## Installation

Run `go get` command.

```
go get github.com/kohkimakimoto/zssh/cmd/zssh
```

## Usage

Create and edit `~/.ssh/zssh.lua`.

```lua
Host "web01.localhost" {
    HostName = "192.168.0.11",
    Port = "22",
    User = "kohkimakimoto",
    ForwardAgent = "yes",
    description = "my web01 server",
}

Host "web02.localhost" {
    HostName = "192.168.0.12",
    Port = "22",
    User = "kohkimakimoto",
    ForwardAgent = "yes",
    description = "my web02 server",
}
```

This configuration genarates the below ssh config when you run `zssh`.

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
