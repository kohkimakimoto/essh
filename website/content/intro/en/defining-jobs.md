+++
title = "Defining Jobs | Introduction"
type = "docs"
category = "intro"
lang = "en"
basename = "defining-jobs.html"
+++

# Defining Jobs

Jobs in Essh encapsulate tasks, hosts and drivers. Hosts and drivers that are defined in a job can be used only by the tasks in the same job. 

Edit your `esshconfig.lua`:

~~~lua
job "myjob" {
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

A Job's tasks have a prefix that is their job's name, So you will run the task like the following

~~~
$ essh myjob:hello
~~~

For more information on jobs, see the [Jobs](/docs/en/jobs.html) section.


## Next Steps

In the [Introduction](/intro/en/index.html) guide, I have explained the basic features of Essh. If you want to get in-depth information about Essh, read the [documentation](/docs/en/index.html).

Enjoy!
