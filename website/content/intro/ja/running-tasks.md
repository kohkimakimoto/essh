+++
title = "Running Tasks"
type = "docs"
category = "intro"
lang = "ja"
basename = "running-tasks.html"
+++

# タスクを実行する

リモートサーバーとローカルサーバーで実行されるタスクを定義できます。
例として、以下のように`esshconfig.lua`を編集します。

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

タスクの詳細については,[Tasks](/docs/ja/tasks.html)セクションを参照してください。

次のセクション: [Luaライブラリを使う](using-lua-libraries.html)
