+++
type = "docs"
lang = "en"
+++

# Zsh Completion

Essh supports zsh completion that lists SSH hosts. If you want to use it, add the following code in your `~/.zshrc`

~~~
eval "$(essh --zsh-completion)"
~~~

And then, edit your `esshconfig.lua`. Try to add the `description` property as the following.

~~~lua
host "web01.localhost" {
    ForwardAgent = "yes",
    HostName = "192.168.0.11",
    Port = "22",
    User = "kohkimakimoto",
    -- add description
    description = "web01 development server",
}

host "web02.localhost" {
    ForwardAgent = "yes",
    HostName = "192.168.0.12",
    Port = "22",
    User = "kohkimakimoto",
    -- add description
    description = "web02 development server",
}
~~~

You will get completion about hosts.

~~~
$ essh [TAB]
web01.localhost  -- web01 development server
web02.localhost  -- web02 development server
~~~

You can hide a host using `hidden` property. If you set it true, zsh completion doesn't show the host.

~~~lua
host "web01.localhost" {
    ForwardAgent = "yes",
    HostName = "192.168.0.11",
    Port = "22",
    User = "kohkimakimoto",
    description = "web01 development server",
    hidden = true,
}
~~~

You notice that the first characters of the `description` and `hidden` are lower case. Others are upper case. It is important point. Essh uses properties whose first character is upper case, as **ssh_config** generated to the temporary file. And the properties whose first character is lower case are used for special purpose of Essh functions, not ssh config.

For more information on hosts, see the [Hosts](/docs/en/hosts.html) section.

Let's read next section: [Using Hooks](using-hooks.html)
