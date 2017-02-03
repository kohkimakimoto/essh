+++
title = "ジョブを定義する | イントロダクション"
type = "docs"
category = "intro"
lang = "ja"
basename = "defining-namespaces.html"
+++

# ジョブを定義する

Esshのジョブはタスク、ホスト、ドライバをカプセル化します。ジョブに定義されているホストとドライバは、同じジョブ内のタスクでのみ使用できます。これはパブリックなホストがタスクのホストと競合するのを防止します。

`esshconfig.lua`を編集してください。

~~~lua
namespace "mynamespace" {
    host "web01.localhost" {
        ForwardAgent = "yes",
        HostName = "192.168.0.11",
        Port = "22",
        User = "kohkimakimoto",
        tags = {
            "web",
        },
    },

    host "web02.localhost" {
        ForwardAgent = "yes",
        HostName = "192.168.0.12",
        Port = "22",
        User = "kohkimakimoto",
        tags = {
            "web",
        },
    },

    task "hello" {
        description = "say hello",
        prefix = true,
        backend = "remote",
        targets = "web",
        script = [=[
            echo "hello on $(hostname)"
        ]=],
    },
}
~~~

ジョブのタスクには、ジョブの名前でプレフィックスが付きます。そのため以下のようにしてタスクを実行します。

~~~
$ essh mynamespace:hello
~~~

ジョブの詳細については、[ジョブ](/docs/ja/namespaces.html)のセクションを参照してください。

## 次のステップ

この[イントロダクション](/intro/ja/index.html)ガイドでは、Esshの基本的な機能について説明しました。 Esshに関する詳細な情報を知りたい場合は、[ドキュメント](/docs/ja/index.html)を参照してください。

それでは。
