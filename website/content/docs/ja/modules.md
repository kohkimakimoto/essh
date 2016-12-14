+++
title = "モジュール | ドキュメント"
type = "docs"
category = "docs"
lang = "ja"
basename = "modules.html"
+++

# モジュール

モジュールを使用すると、簡単にEssh設定のための再利用可能なコードを使用、作成、共有できます。

## モジュールを使う

`import`関数を使ってEsshモジュールをロードすることができます。

~~~lua
local bash = import "github.com/kohkimakimoto/essh/modules/bash"
~~~

`import`はLuaの値を返します。上記の場合`bash`はいくつかの変数と関数を持つLuaのテーブルです。

~~~lua
local bash = import "github.com/kohkimakimoto/essh//modules/bash"

task "example" {
    script = {
        bash.indent,
        "echo hello | indent",
    },
}
~~~


`bash.indent`は[このコードスニペット](https://github.com/kohkimakimoto/essh/blob/master/modules%2Fbash%2Findex.lua#L3-L17)です。
したがって、タスクはインデントされた出力を表示します。

`import`は[hashicorp/go-getter](https://github.com/hashicorp/go-getter)を使って実装されています。 git urlとローカルファイルシステムパスを使用して、モジュールパスを指定することができます。例えば、以下のようなものです。

* githubリポジトリからモジュールを取得します:

    ~~~
    local mod = import "github.com/username/repository"
    ~~~

* githubリポジトリからモジュールを取得し、タグ、コミットまたはブランチをチェックアウトします:

    ~~~
    local mod = import "github.com/username/repository?ref=master"
    ~~~

* githubリポジトリのサブディレクトリからモジュールを取得します:

    ~~~
    local mod = import "github.com/username/repository//path/to/module"
    ~~~

    ダブルスラッシュ`//`はサブディレクトリの区切りで、リポジトリ自体の一部ではありません。

* 汎用のgitリポジトリからモジュールを取得する:
    
    ~~~~
    local mod = import "git::ssh://your-private-git-server/path/to/repo.git"
    ~~~~

* ローカルファイルシステムからモジュールを取得する:

    ~~~
    local mod = import "/path/to/module"
    ~~~

詳細は[hashicorp/go-getter](https://github.com/hashicorp/go-getter)を参照してください。

モジュールはEsshの実行時に自動的にインストールされます。インストールされたモジュールは、`import`が書かれた設定ファイルが現在のディレクトリの`esshconfig.lua`である場合、`.essh/modules`ディレクトリに保存されます。設定ファイルが`~/.essh/config.lua`の場合、モジュールは`~/.essh/modules`ディレクトリに保存されます。

インストールされたモジュールを更新する必要がある場合は`essh --update`を実行してください。

~~~
$ essh --update
~~~

デフォルトでは、Esshは`.essh/modules`内のモジュールのみを更新します。`~/.essh/modules`内のモジュールを更新するには、次のコマンドを実行します：


~~~
$ essh --update --with-global
~~~

## モジュールの作成

新しいモジュールの作成は簡単です。最小のモジュールは`index.lua`だけを含むディレクトリです。
`my_module`ディレクトリと` index.lua`ファイルをディレクトリに作成してみてください。

~~~lua
-- my_module/index.lua
local m = {}

m.hello = "echo hello"

return m
~~~

`index.lua`はLuaの値を返さなければならないエントリポイントです。この例では`hello`変数を持つテーブルを返します。これだけです。このモジュールを使用するには、あなたの設定に次のコードを書いてください。

~~~lua
local my_module = import "./my_module"

task "example" {
    script = {
        my_module.hello,
    },
}
~~~

実行しましょう。

~~~
$ essh example
hello
~~~

モジュールを共有する場合は、モジュールディレクトリからgitリポジトリを作成し、github.comなどのリモートリポジトリにプッシュします。 gitリポジトリのモジュールを使用するには、URLへ`import`パスを更新します。

~~~lua
local my_module = import "github.com/your_account/my_module"

task "example" {
    script = {
        my_module.hello,
    },
}
~~~
