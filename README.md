# ZSSH

Extended ssh command wrapper.

features

* Automatically generates `~/.ssh/config` from `~/.ssh/zssh.lua`
* Support zsh completion.

## Configuration

Create and edit `~/.ssh/zssh.lua`.

```lua
Host "web01.localhost" {
    HostName = "192.168.0.11",
    Port = "22",
    User = "kohkimakimoto",
    ForwardAgent = "yes",
    description = "my web01 server",
}
```

This configuration genarates the below ssh config when you run `zssh`.

```
Host web01.localhost
    ForwardAgent yes
    HostName 192.168.0.11
    Port 22
    User kohkimakimoto
```
