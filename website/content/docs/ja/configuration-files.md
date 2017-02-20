+++
title = "設定ファイル | ドキュメント"
type = "docs"
category = "docs"
lang = "ja"
basename = "configuration-files.html"
+++

# 設定ファイル

Esshの設定は[Lua](https://www.lua.org/)で書きます。設定ファイルではより人間に読みやすい形式のDSL構文を使用できます。

## 例

~~~lua
host "web01.localhost" {
    HostName = "192.168.0.11",
    Port = "22",
    User = "kohkimakimoto",
    description = "web01 development server",
    tags = {
        "web",
    },
}

host "web02.localhost" {
    HostName = "192.168.0.12",
    Port = "22",
    User = "kohkimakimoto",
    description = "web02 development server",
    tags = {
        "web",
    },
}

task "uptime" {
    backend = "remote",
    targets = "web",
    script = "uptime",
}
~~~

## 別の構文

上記の設定例は、LuaのDSLスタイルで書かれています。プレーンなLua関数のスタイルを使うこともできます。以下の例も有効な設定コードです。

~~~lua
host("web01.localhost", {
    HostName = "192.168.0.11",
    Port = "22",
    User = "kohkimakimoto",
    description = "web01 development server",
    tags = {
        "web",
    },
})
~~~

または

~~~lua
local web01 = host "web01.localhost"
web01.HostName = "192.168.0.11"
web01.Port = "22"
web01.User = "kohkimakimoto"
web01.description = "web01 development server"
web01.tags = {
    "web",
}
~~~

## 評価の順序

Esshはいくつかの異なる場所から設定ファイルを読み込みます。設定は次の順序で適用されます。

1. カレントディレクトリに`esshconfig.lua`が存在する場合、これを読み込みます。
1. カレントディレクトリに`esshconfig.lua`が存在しない場合、`〜/.essh/config.lua`を読み込みます。
1. カレントディレクトリの`esshconfig_override.lua`を読み込みます。
1. `~/.essh/config_override.lua`を読み込みます。

`--config`コマンドラインオプションや`ESSH_CONFIG`環境変数を使うと、現在のディレクトリにある読み込みファイルを変更することができます。

## Lua

Esshは、設定ファイルで使用できる組み込みのLuaライブラリを提供しています。

[Lua VM](lua-vm.html)を参照してください。
