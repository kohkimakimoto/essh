+++
title = "Managing Hosts"
type = "docs"
category = "intro"
lang = "ja"
basename = "managing-hosts.html"
+++

# ホストの管理

Esshのホストにはタグを付けることができます。タグを使用するとホストを分類できます。

例として`esshconfig.lua`を編集していくつかのホストとタグを追加します

~~~lua
host "web01.localhost" {
    -- ... your config
    description = "web01 development server",
    tags = {
        "web",
    }
}

host "web02.localhost" {
    -- ... your config
    description = "web02 development server",
    tags = {
        "web",
    }
}

host "db01.localhost" {
    -- ... your config
    description = "db01 server",
    tags = {
        "db",
        "backend",
    }
}

host "cache01.localhost" {
    -- ... your config
    description = "cache01 server",
    tags = {
        "cache",
        "backend",
    }
}
~~~

`essh`を`--hosts`オプションを付けて実行します。

~~~sh
$ essh --hosts
NAME                     DESCRIPTION                     TAGS                 REGISTRY
cache01.localhost        cache01 server                  cache,backend        local
db01.localhost           db01 server                     db,backend           local
web01.localhost          web01 development server        web                  local
web02.localhost          web02 development server        web                  local
~~~

すべてのホストが表示されます。次に`--select`オプションを付けて実行してください。

~~~sh
$ essh --hosts --select=web
NAME                   DESCRIPTION                     TAGS        REGISTRY
web01.localhost        web01 development server        web         local
web02.localhost        web02 development server        web         local
~~~

あなたは`web`タグでフィルタリングされたホストを取得します。`--select`は複数回指定できます。各フィルタはOR条件で適用されます。

~~~sh
$ essh --hosts --select=web --select=db
NAME                   DESCRIPTION                     TAGS              REGISTRY
db01.localhost         db01 server                     db,backend        local
web01.localhost        web01 development server        web               local
web02.localhost        web02 development server        web               local
~~~

ホストの詳細については、[ホスト](/docs/ja/hosts.html) セクションを参照してください。

次のセクション: [コマンドを実行する](running-commands.html)
