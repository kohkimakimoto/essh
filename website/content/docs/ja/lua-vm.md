+++
title = "Lua VM | ドキュメント"
type = "docs"
category = "docs"
lang = "ja"
basename = "lua-vm.html"
+++

# Lua VM

EsshはLuaで書かれた設定ファイルを読み込むために[GopherLua](https://github.com/yuin/gopher-lua)をLua VMとして使っています。

## ビルトイン関数

すでに`host`と`task`関数を見てきたように、Esshのコア機能はビルトイン関数で構成されています。Esshが提供しているすべての関数は以下の通りです。

* `host`: ホストを定義します。[ホスト](/docs/ja/hosts.html)を参照してください。

* `task`: タスクを定義します。[タスク](/docs/ja/tasks.html)を参照してください。

* `driver`: ドライバを定義します。[ドライバ](/docs/ja/drivers.html)を参照してください。

* `job`: ジョブを定義します。[ジョブ](/docs/ja/jobs.html)を参照してください。

* `import`: モジュールをインポートします。[モジュール](/docs/ja/modules.html)を参照してください。

## ビルトインライブラリ

Esshには、設定ファイルで使用できるビルトインLuaライブラリが用意されています。
たとえば、`json`ライブラリを使いたい場合は、以下のようにLuaの`require`関数を使います。

~~~lua
local json = require("json")

local jsontext = json.encode({aaa = "bbb", ccc = "ddd"})
print(jsontext)
~~~

以下は、Esshに組み込まれているビルトインライブラリです。

* `json`: [layeh/gopher-json](https://github.com/layeh/gopher-json).
* `fs`: [kohkimakimoto/gluafs](https://github.com/kohkimakimoto/gluafs).
* `yaml`: [kohkimakimoto/gluayaml](https://github.com/kohkimakimoto/gluayaml).
* `question`: [kohkimakimoto/gluaquestion](https://github.com/kohkimakimoto/gluaquestion).
* `template`: [kohkimakimoto/gluatemplate](https://github.com/kohkimakimoto/gluatemplate).
* `env`: [kohkimakimoto/gluaenv](https://github.com/kohkimakimoto/gluaenv).
* `http`: [cjoudrey/gluahttp](https://github.com/cjoudrey/gluahttp).
* `re`: [yuin/gluare](https://github.com/yuin/gluare)
* `sh`: [otm/gluash](https://github.com/otm/gluash)

## 定義済みの変数

Esshは事前定義された変数を提供します。 最新のEsshのバージョンには、定義済みの変数が1つあります。`essh`です。

`essh`はいくつかの関数と変数を持つテーブルです。下記を参照してください

* `ssh_config` (string): ssh_configはssh_configファイルのパスです。デフォルトでは、Esshを実行すると自動的に生成される一時ファイルです。 ssh_configを静的な宛先に生成するために、この値を上書きすることができます。クライアントコンピュータとターゲットサーバの間のサーバであるゲートウェイホストを使用する場合は、この変数を使用して `ProxyCommand`を指定することができます。以下の例を参照してください：

    ~~~lua
    --
    -- network environment.
    -- [your-computer] -- [getway-server1] -- [web-server]
    --

    host "web-server" {
        HostName = "192.168.0.1",
        ProxyCommand = "ssh -q -F " .. essh.ssh_config .. " -W %h:%p getway-server1",
    }
    ~~~


* `select_hosts` (function): 定義されたホストを取得します。これは、ホスト設定のオーバライドやデフォルト値の設定に役立ちます。たとえば、デフォルトのssh_config:`ForwardAgent = yes`を設定する場合は、以下のコードで実施できます。

    ~~~lua
    -- ~/.essh/config_override.lua
    for _, h in pairs(essh.select_hosts():get()) do
        if h.ForwardAgent == nil then
            h.ForwardAgent = "yes"
        end
    end
    ~~~

    上記の例では、すべてのホストにデフォルト値が設定されます。選択したホストに値を設定したい場合は、次のコードを使います:

    ~~~lua
    -- ~/.essh/config_override.lua
    -- Getting only the hosts that has `web` tag or name of the hosts is `web`.
    for _, h in pairs(essh.select_hosts("web"):get()) do
        if h.ForwardAgent == nil then
            h.ForwardAgent = "yes"
        end
    end

    -- Using a table, Getting the hosts both `web` or `db`
    for _, h in pairs(essh.select_hosts({"web", "db"}):get()) do
        if h.ForwardAgent == nil then
            h.ForwardAgent = "yes"
        end
    end

    -- You can set a filter.
    -- Getting only the `web` hosts filtered by `production`.
    for _, h in pairs(essh.select_hosts("web"):filter("production"):get()) do
        if h.ForwardAgent == nil then
            h.ForwardAgent = "yes"
        end
    end

    -- Getting only the first one host using `first` method.
    local h = essh.select_hosts("web"):first()
    if h.ForwardAgent == nil then
        h.ForwardAgent = "yes"
    end
    ~~~

* `get_job` (function): ジョブを取得します。以下の例を参照してください。

    ~~~lua
    for _, h in pairs(essh.get_job("myjob"):select_hosts():get()) do
        if h.ForwardAgent == nil then
            h.ForwardAgent = "yes"
        end
    end
    ~~~

* `module` (table): インポートされたモジュール内にのみで有効なテーブルです。モジュールスコープの変数として使用できます。いくつか定義済みの値を持ちます。

    * `path` (string): モジュールのパス
    
    * `import_path` (string): import関数の引数に使用されたパス

* `host` (function): `host` 関数のエイリアス。

* `task` (function): `task` 関数のエイリアス。

* `driver` (function): `driver` 関数のエイリアス。

* `job` (function): `job` 関数のエイリアス。

* `import` (function): `import` 関数のエイリアス。

* `debug` (function): デバッグメッセージを出力します。デバッグメッセージは`--debug`オプションつきでEsshを実行したときに出力されます。

    ~~~~lua
    essh.debug("foo")
    ~~~~
