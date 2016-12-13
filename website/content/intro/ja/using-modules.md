+++
title = "モジュールを使う | イントロダクション"
type = "docs"
category = "intro"
lang = "ja"
basename = "using-modules.html"
+++

# モジュールを使う


Esshには設定のための再利用可能なコードを簡単に使用できるモジュラーシステムがあります。
たとえば、Esshタスクで使用するbashスクリプトのコレクションである[bash module](https://github.com/kohkimakimoto/essh/tree/master/modules/bash)があります。
モジュールは`import`関数を使ってロードすることができます。

例:

~~~lua
local bash = import "github.com/kohkimakimoto/essh/modules/bash"

task "example" {
    script = {
        bash.version,
        "echo foo",
    },
}
~~~

`bash.version`は実際には単純な文字列`bash --version`です。このタスクはbash版を出力し`echo foo`を実行します。

モジュールはEssh実行時に自動的にインストールされます。
タスクを実行すると、以下のようになります。

~~~
$ essh example
Installing module: 'github.com/kohkimakimoto/essh/modules/bash' (into /path/to/directory/.essh)
GNU bash, version 4.1.2(1)-release (x86_64-redhat-linux-gnu)
Copyright (C) 2009 Free Software Foundation, Inc.
License GPLv3+: GNU GPL version 3 or later <http://gnu.org/licenses/gpl.html>

This is free software; you are free to change and redistribute it.
There is NO WARRANTY, to the extent permitted by law.
foo
~~~


モジュールの詳細については[モジュール](/docs/ja/modules.html)セクションを参照してください。

次のセクション: [ドライバを使う](using-drivers.html)
