+++
title = "Managing Hosts"
type = "docs"
category = "intro"
lang = "en"
+++

# Managing Hosts

Hosts in Essh can have tags. The tags allow you to classify the hosts.

For instance, edit `esshconfig.lua` to add some hosts and set tags.

```lua
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
```

Run `essh` with `--hosts` option.

```sh
$ essh --hosts
NAME                     DESCRIPTION                     TAGS                 REGISTRY
cache01.localhost        cache01 server                  cache,backend        local
db01.localhost           db01 server                     db,backend           local
web01.localhost          web01 development server        web                  local
web02.localhost          web02 development server        web                  local
```

You will see the all hosts. Next, try to run it with `--select` option.

```sh
$ essh --hosts --select=web
NAME                   DESCRIPTION                     TAGS        REGISTRY
web01.localhost        web01 development server        web         local
web02.localhost        web02 development server        web         local
```

You will get filtered hosts by `web` tag. `--select` can be specified multiple times. Each filters are used in OR condition.

```sh
$ essh --hosts --select=web --select=db
NAME                   DESCRIPTION                     TAGS              REGISTRY
db01.localhost         db01 server                     db,backend        local
web01.localhost        web01 development server        web               local
web02.localhost        web02 development server        web               local
```

For more information on hosts, see the [Hosts](/docs/en/hosts.html) section.

Let's read next section: [Running Commands](running-commands.html)
