# ZSSH

Extended ssh command wrapper.

features

* Automatically generates `~/.ssh/config` from `~/.ssh/zssh.lua`
* Support zsh completion.

## Installation

### Installation

#### Compiled binary

zssh is provided as a single binary. You can download it below links and drop it in your $PATH.

[Download latest version](https://github.com/kohkimakimoto/zssh/releases/latest)

#### Using ***go get*** command

You can also install zssh by the `go get` command like the following.

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

You will get like the below completion.

```
$ zssh [TAB]
web01.localhost          -- my web01 server
web02.localhost          -- my web02 server
web03.localhost          -- my web03 server
```
