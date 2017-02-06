+++
title = "Groups | Documentation"
type = "docs"
category = "docs"
lang = "en"
basename = "groups.html"
+++

# Groups

## Example

~~~lua
group {
    hidden = true,
    privileged = true,
    backend = "remote",
    targets = {"web"},
    
    task "foo" {
        script = "echo foo"
    },

    task "foo" {
        script = "echo foo"
    },
}
~~~