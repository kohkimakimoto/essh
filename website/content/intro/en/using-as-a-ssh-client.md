+++
title = "Using As A SSH Client | Introduction"
type = "docs"
category = "intro"
lang = "en"
basename = "using-as-a-ssh-client.html"
+++

# Using As A SSH Client

Essh is implemented as a wrapper of `ssh` command. That means you can use Essh in the same way as `ssh`. Try to connect a remote server by using Essh instead of `ssh` command.

Create `esshconfig.lua` in your current directory. This is a default configuration file for Essh. The configuration is written in [Lua](https://www.lua.org/) programming language. Now edit this file as the following.

> Replace the `HostName`, `User` and some parameters for your environment.

~~~lua
host "web01.localhost" {
    ForwardAgent = "yes",
    HostName = "192.168.0.11",
    Port = "22",
    User = "kohkimakimoto",
}

host "web02.localhost" {
    ForwardAgent = "yes",
    HostName = "192.168.0.12",
    Port = "22",
    User = "kohkimakimoto",
}
~~~

This configuration automatically generates the below ***ssh_config*** to the temporary file like the `/tmp/essh.ssh_config.260398422` whenever you run `essh`.

~~~
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
~~~

Essh uses this generated config file by default. If you run the below command

~~~
$ essh web01.localhost
~~~

Essh internally runs the `ssh` command like the following.

~~~
$ ssh -F /tmp/essh.ssh_config.260398422 web01.localhost
~~~

Therefore you can connect with a ssh server using Lua config. If you want to see the generated ***ssh_config***, use `--print` options.

~~~
$ essh --print
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
~~~

Essh also automatically removes the temporary file when the process finishes. So you don't have to be conscious of the real ssh configuration in the normal operations.

Essh configuration file can also be placed at `~/.essh/config.lua`. This is a per user's configuration file like a `~/.ssh/config`.

For more information on configuration files, see the [Configuration Files](/docs/en/configuration-files.html) section.

Let's read next section: [Zsh Completion](zsh-completion.html).
