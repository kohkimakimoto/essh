+++
title = "ジョブを定義する | イントロダクション"
type = "docs"
category = "intro"
lang = "ja"
basename = "defining-jobs.html"
+++

# ジョブを定義する

Esshのジョブはタスク、ホスト、ドライバをカプセル化します。ジョブに定義されているホストとドライバは、同じジョブ内のタスクでのみ使用できます。

`esshconfig.lua`を編集してください。

~~~lua
job "myjob" {
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
$ essh myjob:hello
~~~

ジョブの詳細については、[ジョブ](/docs/ja/jobs.html)のセクションを参照してください。

## 次のステップ

この[イントロダクション](/intro/ja/index.html)ガイドでは、Esshの基本的な機能について説明しました。 Esshに関する詳細な情報を知りたい場合は、[ドキュメント](/docs/ja/index.html)を参照してください。

それでは。
