+++
title = "設定ファイル | ドキュメント"
type = "docs"
category = "docs"
lang = "ja"
basename = "configuration-files.html"
+++

# 設定ファイル

Esshの設定は[Lua](https://www.lua.org/)で書きます。設定ファイルではより人間に読みやすい形式のDSL構文を使用できます。

以下は例です:

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

Esshはいくつかの異なる場所から設定ファイルを読み込みます。必須である設定ファイルはありません。Esshはこれらが存在する場合のみロードします。設定は次の順序で適用されます。

1. `/etc/essh/config.lua` (`global`レジストリ)。
1. `~/.essh/config.lua` (`global` レジストリ)。
1. カレントディレクトリの`esshconfig.lua`またはコマンドラインオプションの`--config`で指定したファイル (`local` レジストリ).
1. カレントディレクトリの`esshconfig_override.lua`またはコマンドラインオプションの`--config`で指定したファイル名の最後に`_override`をつけたファイル。例)`--config=myconfig.lua`なら`myconfig_override.lua` (`local` レジストリ)。
1. `~/.essh/config_override.lua` (`global` レジストリ)。
1. `/etc/essh/config_override.lua` (`global` レジストリ)。

## レジストリ {#registries}

各設定ファイルには**レジストリ**という概念があります。レジストリは、設定がロードされることによって定義されるホストやタスクなどのリソースを保持する論理空間です。

Esshには**local**と**global**の2つのレジストリがあります。

### 例

`/etc/essh/config.lua`にホストを定義すると、このホストの設定は`global`レジストリに格納されます。

### 制約

レジストリはリソースに関するいくつかの制約を提供します。たとえば（最も重要な制約）は次のとおりです。

> タスクによって使用されるホストは、同じレジストリに定義する必要がある。

`global`レジストリでタスクを定義した場合、このタスクは`local`レジストリで定義されたホストを使用できません。

詳細については、[ホスト](hosts.html)を参照してください。

## Lua

Esshは、設定ファイルで使用できる組み込みのLuaライブラリを提供しています。

[Lua VM](lua-vm.html)を参照してください。
