# Managing Hosts

Hosts in Essh can have tags. The tags allow you to classify hosts.

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

```
$ essh --hosts
NAME                     DESCRIPTION                     TAGS                 REGISTRY        SCOPE         HIDDEN
cache01.localhost        cache01 server                  cache,backend        local           public        false
db01.localhost           db01 server                     db,backend           local           public        false
web01.localhost          web01 development server        web                  local           public        false
web02.localhost          web02 development server        web                  local           public        false
```

You will see the all hosts. Next, try to run it with `--select` option.

```
$ essh --hosts --select=web
NAME                   DESCRIPTION                     TAGS        REGISTRY        SCOPE         HIDDEN
web01.localhost        web01 development server        web         local           public        false
web02.localhost        web02 development server        web         local           public        false
```

You will get filtered hosts by `web` tag. `--select` can be specified multiple times. Each filters are used in OR condition.

```
$ essh --hosts --select=web --select=db
NAME                   DESCRIPTION                     TAGS              REGISTRY        SCOPE         HIDDEN
db01.localhost         db01 server                     db,backend        local           public        false
web01.localhost        web01 development server        web               local           public        false
web02.localhost        web02 development server        web               local           public        false
```

For more information on hosts, see the [Hosts](#hosts) section.

Let's read next section: [Running Commands](getting-started_running-commands.md)
