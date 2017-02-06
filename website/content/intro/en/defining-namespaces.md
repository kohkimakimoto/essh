+++
title = "Defining Namespaces | Introduction"
type = "docs"
category = "intro"
lang = "en"
basename = "defining-namespaces.html"
+++

# Defining Namespaces

Namespaces in Essh encapsulate tasks, hosts and drivers. Hosts and drivers that are defined in a namespace can be used only by the tasks in the same namespace. It prevents to conflict public hosts with task's hosts.

Edit your `esshconfig.lua`:

~~~lua
namespace "mynamespace" {
    host "web01.localhost" {
        ForwardAgent = "yes",
        HostName = "192.168.0.11",
        Port = "22",
        User = "kohkimakimoto",
        tags = {
            "web",
        },
    },

    host "web02.localhost" {
        ForwardAgent = "yes",
        HostName = "192.168.0.12",
        Port = "22",
        User = "kohkimakimoto",
        tags = {
            "web",
        },
    },

    task "hello" {
        description = "say hello",
        prefix = true,
        backend = "remote",
        targets = "web",
        script = [=[
            echo "hello on $(hostname)"
        ]=],
    },
}
~~~

A Namespace's tasks have a prefix that is their namespace's name, So you can run the task like the following

~~~
$ essh mynamespace:hello
~~~

For more information on namespaces, see the [Namespaces](/docs/en/namespaces.html) section.


## Next Steps

In the [Introduction](/intro/en/index.html) guide, I have explained the basic features of Essh. If you want to get in-depth information about Essh, read the [documentation](/docs/en/index.html).

Enjoy!
