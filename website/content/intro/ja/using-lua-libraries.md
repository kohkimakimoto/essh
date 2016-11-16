+++
title = "Using Lua Libraries"
type = "docs"
category = "intro"
lang = "ja"
basename = "using-lua-libraries.html"
+++

# Luaライブラリを使う

Esshは設定にLuaを使用し、いくつかの組み込みのLuaライブラリも持っています。ライブラリをロードするには`require`関数を使います。

例:

~~~lua
local question = require "question"

task "example" {
    prepare = function ()
        local r = question.ask("Are you OK? [y/N]: ")
        if r ~= "y" then
            -- return false, the task does not run.
            return false
        end
    end,
    script = [=[
        echo "foo"
    ]=],
}
~~~

`question`はEsshの組み込みライブラリで[gluaquestion]（https://github.com/kohkimakimoto/gluaquestion）によって実装されています。これはターミナルからのユーザ入力を取得する機能を提供します。
タスクのプロパティ`prepare`は、タスクの開始時に実行される関数を定義する設定です。

タスクを実行すると、Esshはメッセージを表示し、入力を待ちます。

~~~
$ essh example
Are you OK? [y/N]: y
foo
~~~

Luaライブラリの詳細については、[Lua VM](/docs/ja/lua-vm.html)セクションを参照してください。

次のセクション: [モジュールを使う](using-modules.html)
