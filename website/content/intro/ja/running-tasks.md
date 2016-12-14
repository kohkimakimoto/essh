+++
title = "タスクを実行する | イントロダクション"
type = "docs"
category = "intro"
lang = "ja"
basename = "running-tasks.html"
+++

# タスクを実行する

タスクは、リモートサーバまたはローカルで実行されるスクリプトです。
それでは例として、以下のように`esshconfig.lua`を編集してみましょう。

~~~lua
task "hello" {
    description = "say hello",
    prefix = true,
    backend = "remote",
    targets = "web",
    script = [=[
        echo "hello on $(hostname)"
    ]=],
}
~~~

タスクを実行します。

~~~
$ essh hello
[remote:web01.localhost] hello on web01.localhost
[remote:web02.localhost] hello on web02.localhost
~~~

`backend`プロパティに`local`を設定すると、Esshはローカルでタスクを実行します。

~~~lua
task "hello" {
    description = "say hello",
    prefix = true,
    backend = "local",
    script = [=[
        echo "hello on $(hostname)"
    ]=],
}
~~~

~~~
$ essh hello
[local] hello on your-hostname
~~~

タスクの詳細については,[タスク](/docs/ja/tasks.html)セクションを参照してください。

次のセクションに進みましょう: [Luaライブラリを使う](using-lua-libraries.html)
