+++
title = "Using Hooks | Introduction"
type = "docs"
category = "intro"
lang = "en"
basename = "using-hooks.html"
+++

# Using Hooks

Hooks in Essh are scripts executed before and after connecting remote servers.

Write the following code in your `esshconfig.lua`.

~~~lua
host "web01.localhost" {
    HostName = "192.168.0.11",
    Port = "22",
    User = "kohkimakimoto",

    hooks_before_connect = {
        "echo before_connect: $HOSTNAME",
    },
    hooks_after_connect = {
        "echo after_connect: $HOSTNAME",
    },
    hooks_after_disconnect = {
        "echo after_disconnect: $HOSTNAME",
    },
}
~~~

Connect with the server.

~~~
$ essh web01.localhost 
before_connect: your-local-machine
after_connect: web01.localhost
[kohkimakimoto@web01.localhost ~]$ 
~~~

The `hooks_before_connect` and `hooks_after_connect` were executed. Disconnect with the server.

~~~
[kohkimakimoto@web01.localhost ~]$ exit
exit
Connection to 192.168.0.11 closed.
after_disconnect: your-local-machine
~~~

The `hooks_after_disconnect` was executed.

Essh supports below type of hooks:

* `hooks_before_connect`
* `hooks_after_connect`
* `hooks_after_disconnect`


Let's read next section: [Managing Hosts](managing-hosts.html).
