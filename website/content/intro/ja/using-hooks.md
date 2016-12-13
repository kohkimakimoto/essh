+++
title = "フックを使う | イントロダクション"
type = "docs"
category = "intro"
lang = "ja"
basename = "using-hooks.html"
+++

# フックを使う

Esshのフックは、リモートサーバーを接続する前後に実行されるスクリプトです。

例:

~~~lua
host "web01.localhost" {
    HostName = "192.168.0.11",
    Port = "22",
    User = "kohkimakimoto",

    hooks_before_connect = {
        "echo before_connect",
    },
    hooks_after_connect = {
        "echo after_connect",
    },
    hooks_after_disconnect = {
        "echo after_disconnect",
    },
}
~~~

Esshは以下のタイプのフックをサポートしています:

* `hooks_before_connect`
* `hooks_after_connect`
* `hooks_after_disconnect`

次のセクション: [ホストの管理](managing-hosts.html)
