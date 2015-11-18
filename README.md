# ZSSH

Extended ssh command.

* This is a single binary CLI app.
* Simply wraps `ssh` command. You can use it in the same way as `ssh`.
* Automatically generates `~/.ssh/config` from `~/.ssh/zssh.lua`. You can write SSH configuration in Lua programming language.
* Supports zsh completion.
* Provides some hook functions.
* You can define macros to execute command via SSH.

![zssh.gif](zssh.gif)


## Installation

#### Compiled binary

ZSSH is provided as a single binary. You can download it and drop it in your $PATH.

[Download latest version](https://github.com/kohkimakimoto/zssh/releases/latest)

#### Using ***go get*** command

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

If you set a first character of keys as lower case like `description`, it is not SSH config.
It uses specific functionality. Read the next section **Zsh Completion**.

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
            -- This is an example to change screen color to red.
            os.execute("osascript -e 'tell application \"Terminal\" to set current settings of first window to settings set \"Red Sands\"'")
        end,
        after = function()
            -- This is an example to change screen color to black.
            os.execute("osascript -e 'tell application \"Terminal\" to set current settings of first window to settings set \"Pro\"'")
        end,
    }
}
```

`before` hook fires before you connect server via SSH. `after` hook fires after you disconnect SSH connection.


### Macros

You can define macros to run commands local or remote hosts.

```lua
Host "web01.localhost" {
    HostName = "192.168.0.11",
    Port = "22",
    User = "kohkimakimoto",
    ForwardAgent = "yes",
    description = "my web01 server",
    tags = {
        role = "web"
    },
}

Host "web02.localhost" {
    HostName = "192.168.0.12",
    Port = "22",
    User = "kohkimakimoto",
    ForwardAgent = "yes",
    description = "my web02 server",
    tags = {
        role = "web"
    },
}

Macro "example" {
    -- parallel execution: default false
    parallel = true,
    -- display confirm prompt: default false
    confirm = "Are you OK?",
    -- description that is showed on zsh completion.
    description = "example macro",
    -- specify remote servers to run a command by tags. if it isn't set, runs command locally.
    on = {role = "web"},
    -- allocate tty: default false
    tty = false,
    -- command.
    command = [[
        ls -la
    ]],
}
```

Run a macro.

```
$ zssh example
```

## Author

Kohki Makimoto <kohki.makimoto@gmail.com>

## License

MIT license.
