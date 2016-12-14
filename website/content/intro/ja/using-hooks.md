+++
title = "フックを使う | イントロダクション"
type = "docs"
category = "intro"
lang = "ja"
basename = "using-hooks.html"
+++

# フックを使う

Esshのフックは、リモートサーバーを接続する前後に実行されるスクリプトです。

以下のコードを`esshconfig.lua`に書いてください。

~~~lua
host "web01.localhost" {
    HostName = "192.168.0.11",
    Port = "22",
    User = "kohkimakimoto",

    hooks_before_connect = {
        "echo before_connect: $HOSTNAME",
    },
    hooks_after_connect = {
        "echo after_connect: $HOSTNAME",
    },
    hooks_after_disconnect = {
        "echo after_disconnect: $HOSTNAME",
    },
}
~~~

サーバに接続します。

~~~
$ essh web01.localhost 
before_connect: your-local-machine
after_connect: web01.localhost
[kohkimakimoto@web01.localhost ~]$ 
~~~

`hooks_before_connect` と `hooks_after_connect` が実行されました。サーバから切断してみましょう。

~~~
[kohkimakimoto@web01.localhost ~]$ exit
exit
Connection to 192.168.0.11 closed.
after_disconnect: your-local-machine
~~~

`hooks_after_disconnect` が実行されました。

Esshは以下のタイプのフックをサポートしています:

* `hooks_before_connect`
* `hooks_after_connect`
* `hooks_after_disconnect`

次のセクションに進みましょう: [ホストの管理](managing-hosts.html)
