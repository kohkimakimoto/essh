+++
title = "ネームスペース | ドキュメント"
type = "docs"
category = "docs"
lang = "ja"
basename = "namespaces.html"
+++

# ネームスペース

Esshのネームスペースは、タスク、ホスト、ドライバをカプセル化します。名前空間に定義されているホストとドライバは、同じ名前空間内のタスクでのみ使用できます。

## Example

~~~lua
namespace "mynamespace" {
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


