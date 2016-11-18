+++
title = "他のツールとの統合 | ドキュメント"
type = "docs"
category = "docs"
lang = "ja"
basename = "integrating-other-tools.html"
+++

# 他のツールとの統合

Esshは `scp`、`rsync`、`git`と一緒に使うことができます。

## git

gitコマンドの中でEsshを使うには`~/.zshrc`次の行を書いてください

~~~
export GIT_SSH=essh
~~~

## scp

Esshはscpで使用することをサポートしています。

~~~
$ essh --exec 'scp -F $ESSH_SSH_CONFIG <scp command args...>'
~~~

もっと使いやすくするため`〜/.zshrc`で`eval "$(essh --aliases)"`を実行すると、上記のコードは以下のように書くことができます。

~~~
$ escp <scp command args...>
~~~

## rsync

Esshはrsyncで使用することをサポートしています。

~~~
$ essh --exec 'rsync -e "ssh -F $ESSH_SSH_CONFIG" <rsync command args...>'
~~~

もっと使いやすくするため`〜/.zshrc`で`eval "$(essh --aliases)"`を実行すると、上記のコードは以下のように書くことができます。

~~~
$ ersync <rsync command args...>
~~~
