+++
title = "Using Hooks"
type = "docs"
category = "intro"
lang = "en"
basename = "using-hooks.html"
+++

# Using Hooks

Hooks in Essh are scripts executed before and after connecting remote servers.

Example:

~~~lua
host "web01.localhost" {
    HostName = "192.168.0.11",
    Port = "22",
    User = "kohkimakimoto",

    hooks_before_connect = {
        "echo before_connect",
    },
    hooks_after_connect = {
        "echo after_connect",
    },
    hooks_after_disconnect = {
        "echo after_disconnect",
    },
}
~~~

Essh supports below type of hooks:

* `hooks_before_connect` (table): fires on the localhost before you connect a server via SSH.

* `hooks_after_connect` (table): fires on the remote host after you connect a server via SSH.

* `hooks_after_disconnect` (table): fires on the local host after you disconnect from a SSH server.

Let's read next section: [Managing Hosts](managing-hosts.html).
