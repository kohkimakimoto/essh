+++
title = "ネームスペース | ドキュメント"
type = "docs"
category = "docs"
lang = "ja"
basename = "namespaces.html"
+++

# ネームスペース

Namespaces in Essh encapsulate tasks, hosts and drivers. Hosts and drivers that are defined in a namespace can be used only by the tasks in the same namespace.

## Example

~~~lua
namespace "mynamespace" {
    -- define description of the namespace
    description = "this is my namespace",
    
    -- If you set it true, defined tasks in this namespace are hidden at default.
    hidden = false,
    
    -- prepare function 
    prepare = function()
    
    end,
    
    -- defining hosts
    host "web01.localhost" {
        -- ...
    },

    host "web02.localhost" {
        -- ...
    },
    
    -- defining drivers
    driver "default" {
        -- ...
    },

    driver "custom" {
        -- ...
    },

    -- defining tasks
    task "foo" {
        -- ...
    },
    
    task "bar" {
        -- ...
    },
    
}
~~~

You can also define hosts, drivers and tasks by using tables. see below example:

~~~lua
namespace "mynamespace" {
    hosts = {
        ["web01.localhost"] = {
            -- ...
        },
        ["web02.localhost"] = {
            -- ...
        },
    },
    drivers = {
        ["default"] = {
            -- ...
        },
        ["custom"] = {
            -- ...
        },
    }
    tasks = {
        ["foo"] = {
            --- ...
        },
        ["bar"] = {
            --- ...
        }
    }
    
}
~~~

## Running namespace's task

A Namespace’s tasks have a prefix that is their namespace’s name, So you can run the task like the following

~~~
$ essh mynamespace:foo
~~~

## Default namespace 

If you define `default` namespace like the following. This namespace's tasks can be run without prefix.

~~~lua
namespace "default" {
    task "foo" {
        -- ...
    }
}
~~~

You can run the `foo` task without prefix.

~~~
$ essh foo
~~~


