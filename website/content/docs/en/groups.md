+++
title = "Groups | Documentation"
type = "docs"
category = "docs"
lang = "en"
basename = "groups.html"
+++

# Groups

Groups are used for defining default parameters for hosts, tasks and drivers. See below example. At one group can only include one type of resource. For example, If you define a group with some hosts definitions, you can not define tasks and drivers in this group.

## Example

~~~lua
group {
    -- Define default parameters.
    hidden = true,
    privileged = true,
    backend = "remote",
    targets = {"web"},
    
    task "foo" {
        script = "echo foo"
    },

    task "foo" {
        -- You can override parameters.
        hidden = false,
        script = "echo foo"
    },
}

group {
    User = "kohkimakimoto",
    
    -- You can define only one type resource in a group.
    host "web01" {
        -- ...
    },
    
    host "web02" {
        -- ...
    },
}
~~~

