+++
title = "グールプ | ドキュメント"
type = "docs"
category = "docs"
lang = "ja"
basename = "groups.html"
+++

# グループ

グループは、ホスト、タスク、およびドライバのデフォルトパラメータを定義するために使用されます。下記の例を参照してください。 1つのグループには、1種類のリソースしか含めることができません。たとえば、いくつかのホスト定義を持つグループを定義した場合、このグループにタスクとドライバを定義することはできません。

## 例

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