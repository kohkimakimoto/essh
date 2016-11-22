+++
title = "Managing Hosts"
type = "docs"
category = "intro"
lang = "en"
basename = "managing-hosts.html"
+++

# Managing Hosts

Hosts in Essh can have tags. The tags allow you to classify the hosts.

For instance, edit `esshconfig.lua` to add some hosts and set tags.

~~~lua
host "web01.localhost" {
    -- ... your config
    description = "web01 development server",
    tags = {
        "web",
    }
}

host "web02.localhost" {
    -- ... your config
    description = "web02 development server",
    tags = {
        "web",
    }
}

host "db01.localhost" {
    -- ... your config
    description = "db01 server",
    tags = {
        "db",
        "backend",
    }
}

host "cache01.localhost" {
    -- ... your config
    description = "cache01 server",
    tags = {
        "cache",
        "backend",
    }
}
~~~

Run `essh` with `--hosts` option.

~~~sh
$ essh --hosts
SCOPE         NAME                     DESCRIPTION                     TAGS                 REGISTRY        HIDDEN
public        cache01.localhost        cache01 server                  cache,backend        local           false
public        db01.localhost           db01 server                     db,backend           local           false
public        web01.localhost          web01 development server        web                  local           false
public        web02.localhost          web02 development server        web                  local           false
~~~

You will see the all hosts. Next, try to run it with `--select` option.

~~~sh
$ essh --hosts --select=web
SCOPE         NAME                   DESCRIPTION                     TAGS        REGISTRY        HIDDEN
public        web01.localhost        web01 development server        web         local           false
public        web02.localhost        web02 development server        web         local           false
~~~

You will get filtered hosts by `web` tag. `--select` can be specified multiple times. Each filters are used in OR condition.

~~~sh
$ essh --hosts --select=web --select=db
SCOPE         NAME                   DESCRIPTION                     TAGS              REGISTRY        HIDDEN
public        db01.localhost         db01 server                     db,backend        local           false
public        web01.localhost        web01 development server        web               local           false
public        web02.localhost        web02 development server        web               local           false
~~~

For more information on hosts, see the [Hosts](/docs/en/hosts.html) section.

Let's read next section: [Running Commands](running-commands.html)
