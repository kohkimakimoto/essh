+++
title = "グールプ | ドキュメント"
type = "docs"
category = "docs"
lang = "ja"
basename = "groups.html"
+++

# グループ

## 例

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