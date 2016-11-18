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

`--target`オプションは複数指定することができます。

~~~sh
$ essh --exec --backend=remote --target=web --target=db uptime
 16:47:02 up 270 days, 13:29,  0 users,  load average: 0.11, 0.18, 0.11
 16:47:02 up 270 days, 13:26,  0 users,  load average: 0.00, 0.01, 0.00
 16:47:02 up 10 days,  1:02,  0 users,  load average: 0.01, 0.03, 0.00
 16:47:03 up 2 days, 22:24,  1 user,  load average: 0.00, 0.01, 0.05
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

`--privileged`オプションを使うと、特権(root)ユーザでコマンドを実行します。

~~~sh
$ essh --exec --backend=remote --target=web --prefix --privileged whoami
[remote:web01.localhost] root
[remote:web01.localhost] root
~~~

`--backend=local`をセットすると, Esshはローカルでコマンドを実行します

~~~sh
$ essh --exec --backend=local --target=web --parallel --prefix 'echo $ESSH_HOSTNAME'
[local:web01.localhost] web01.localhost
[local:web02.localhost] web02.localhost
~~~

上記の例では`ESSH_HOSTNAME`環境変数を使用しています。Esshは内部的に一時的な[タスク](/docs/ja/tasks.html)を使用してコマンドを実行します。したがって、定義済みの変数を使用することができます。[タスク](/docs/en/tasks.html)を参照してください。

次のセクション: [タスクを実行する](running-tasks.html)
