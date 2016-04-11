# ESSH

Extended ssh command.

* Single binary CLI app.
* Simply wraps `ssh` command. You can use it in the same way as `ssh`.
* Supports to write SSH configuration in Lua programming language.
* Supports zsh completion.
* Provides some hook functions.
* Provides utility for managing remote hosts.
* Composes tasks for the remote hosts.

**Now it is on unstable stage**

Table of contents

* [Getting Started](#getting-started)
  * [Installation](#installation)
  * [Usage](#usage)
  * [Zsh Completion](#zsh-completion)
* [Configuration](#configuration)
  * [Hooks](#hooks)
  * [Variables](#variables)
* [Using with git](#using-with-git)
* [Author](#author)
* [License](#license)

## Getting Started

### Installation

ESSH is provided as a single binary. You can download it and drop it in your $PATH.

[Download latest version](https://github.com/kohkimakimoto/essh/releases/latest)

### Usage

Create and edit `~/.essh/config.lua`. This is a main configuration file for ESSH.
The configuration is written in Lua programming language.

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

This configuration generates the below ssh config to the temporary file like the `/tmp/essh.ssh_config.260398422` when you run `essh`.

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

ESSH uses the generated config file by default.
It internally runs the command like the `ssh -F /tmp/essh.ssh_config.260398422 <hostname>`.
And automatically removes the temporary file when `essh` process finishes.
So you can connect a server using below simple command.

```
$ essh web01.localhost
```

If you set a first character of keys as lower case like `description` in the config file, it is not SSH config.
It is used for specific functionality. Read the next section **Zsh Completion**.

### Zsh Completion

If you want to use zsh completion, add the following code in your `~/.zshrc`

```
eval "$(essh --zsh-completion)"
```

You will get completion about hosts.

```
$ essh [TAB]
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


## Configuration

### Hooks

You can add hook `before_connect`, `after_connect` and `after_disconnect` in a host configuration.

```lua
Host "web01.localhost" {
    HostName = "192.168.0.11",
    Port = "22",
    User = "kohkimakimoto",
    ForwardAgent = "yes",
    description = "my web01 server",
    hooks = {
        -- Runs the script on the local before connecting. This is an example to change screen color to red.
        before_connect = "osascript -e 'tell application \"Terminal\" to set current settings of first window to settings set \"Red Sands\"'",

        -- Runs the script on the remote after connecting.
        after_connect = [=[
        echo "Connected to $(hostname)"
        ]=],

        -- Runs the script on the local after disconnecting. This is an example to change screen color to black.
        after_disconnect = "osascript -e 'tell application \"Terminal\" to set current settings of first window to settings set \"Pro\"'",
    }
}
```

`before_connect` and `after_disconnect` also can be written as Lua function instead of shell script.

### Variables

ESSH provides `essh` object to the Lua context. And you can set and get below variable.

#### ssh_config

`ssh_config` is generated config file path. At default, a temporary file path when you run `essh`.

You can set static file path. For instance, you set `essh.ssh_config = os.getenv("HOME") .. "/.ssh/config"`, ESSH overrides `~/.ssh/config` that is standard ssh config file per user.

Example:

```lua
essh.ssh_config = os.getenv("HOME") .. "/.ssh/config"

Host "web01.localhost" {
    ForwardAgent = "yes",
    HostName = "192.168.0.11",
    Port = "22",
    User = "kohkimakimoto",
    description = "my web01 server",
    hidden = true,
}
```

## Using with git

Write the following line in your `~/.zshrc`.

```
export GIT_SSH=essh
```

## Author

Kohki Makimoto <kohki.makimoto@gmail.com>

## License

The MIT License (MIT)
