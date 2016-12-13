+++
title = "SSHクライアントとして使う | イントロダクション"
type = "docs"
category = "intro"
lang = "ja"
basename = "using-as-a-ssh-client.html"
+++

# SSHクライアントとして使う

Esshは`ssh`コマンドのラッパーとして実装されています。つまりEsshは`ssh`と同じように使うことができます。`ssh`コマンドの代わりにEsshを使ってリモートサーバに接続してみましょう。

カレントディレクトリに`esshconfig.lua`を作成します。これはEsshのデフォルト設定ファイルです。設定は[Lua](https://www.lua.org/)プログラミング言語で書かれています。このファイルを次のように編集します。

> `HostName`や`User`などのパラメータは自分の環境に合わせて置き換えてください。

~~~lua
host "web01.localhost" {
    ForwardAgent = "yes",
    HostName = "192.168.0.11",
    Port = "22",
    User = "kohkimakimoto",
}

host "web02.localhost" {
    ForwardAgent = "yes",
    HostName = "192.168.0.12",
    Port = "22",
    User = "kohkimakimoto",
}
~~~

この設定は`essh`を実行するたびに`/tmp/essh.ssh_config.260398422`のような一時ファイルに以下の***ssh_config***を自動的に生成します。

~~~
Host web01.localhost
    ForwardAgent yes
    HostName 192.168.0.11
    Port 22
    User kohkimakimoto

Host web02.localhost
    ForwardAgent yes
    HostName 192.168.0.12
    Port 22
    User kohkimakimoto
~~~

Esshはデフォルトでこの生成された設定ファイルを使用します。以下のコマンドを実行すると、

~~~
$ essh web01.localhost
~~~

Esshは内部的に次のような`ssh`コマンドを実行します。

~~~
$ ssh -F /tmp/essh.ssh_config.260398422 web01.localhost
~~~

このようにしてLua設定を使用してsshサーバに接続することができます。生成される***ssh_config***を確認したいのなら、`--print`オプションを使ってください。

~~~
$ essh --print
Host web01.localhost
    ForwardAgent yes
    HostName 192.168.0.11
    Port 22
    User kohkimakimoto

Host web02.localhost
    ForwardAgent yes
    HostName 192.168.0.12
    Port 22
    User kohkimakimoto
~~~


プロセスが終了すると、Esshは自動的に一時ファイルを削除します。したがって、通常の操作で実際のssh設定を意識する必要はありません。

次のセクション: [Zsh補完](zsh-completion.html)
