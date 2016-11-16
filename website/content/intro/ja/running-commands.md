+++
title = "Running Commands"
type = "docs"
category = "intro"
lang = "ja"
basename = "running-commands.html"
+++

# コマンドを実行する

Esshでは`--exec`、`--backend`、`--target`オプションを使って、選択したリモートホスト上でコマンドを実行することができます。

~~~sh
$ essh --exec --backend=remote --target=web uptime
 22:48:31 up  7:58,  0 users,  load average: 0.00, 0.01, 0.03
 22:48:31 up  7:58,  0 users,  load average: 0.00, 0.02, 0.04
~~~

`--prefix`オプションを使うとEsshはホスト名をプリフィクスにつけてコマンドの結果を出力します。

~~~sh
$ essh --exec --backend=remote --target=web --prefix uptime
[remote:web01.localhost]  22:48:31 up  7:58,  0 users,  load average: 0.00, 0.01, 0.03
[remote:web02.localhost]  22:48:31 up  7:58,  0 users,  load average: 0.00, 0.02, 0.04
~~~

`--parallel`オプションを使うと、並列にコマンドを実行します。

~~~sh
$ essh --exec --backend=remote --target=web --prefix --parallel uptime
[remote:web01.localhost]  22:48:31 up  7:58,  0 users,  load average: 0.00, 0.01, 0.03
[remote:web02.localhost]  22:48:31 up  7:58,  0 users,  load average: 0.00, 0.02, 0.04
~~~

次のセクション: [タスクを実行する](running-tasks.html)
