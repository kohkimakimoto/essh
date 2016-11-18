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

* `host`: Defines a host. See [Hosts](/docs/ja/hosts.html).

* `private_host`: Defines a private host. See [Hosts](/docs/ja/hosts.html).

* `task`: Defines a task. See [Tasks](/docs/ja/tasks.html).

* `driver`: Defines a driver. See [Drivers](/docs/ja/drivers.html).

* `import`: Imports a module. See [Modules](/docs/ja/modules.html).

* `find_hosts`: Gets defined hosts. It is useful for overriding host config and set default values. For example, if you want to set a default ssh config: `ForwardAgent = yes`, you can achieve it the below code:

    ~~~lua
    -- ~/.essh/config_override.lua
    for _, h in pairs(find_hosts():get()) do
        if h.ForwardAgent == nil then
            h.ForwardAgent = "yes"
        end
    end
    ~~~

    Above example sets the default value to the all hosts. If you want to set the value to some hosts you selected, You can also use the below code:

    ~~~lua
    -- ~/.essh/config_override.lua
    -- Getting only the hosts that has `web` tag or name of the hosts is `web`.
    for _, h in pairs(find_hosts("web"):get()) do
        if h.ForwardAgent == nil then
            h.ForwardAgent = "yes"
        end
    end

    -- You can set filter multiple times.
    -- Getting only the hosts filtered by `web` and `production`.
    for _, h in pairs(find_hosts("web"):filter("production"):get()) do
        if h.ForwardAgent == nil then
            h.ForwardAgent = "yes"
        end
    end

    -- Getting only the first one host using `first` method.
    local h = find_hosts("web"):first()
    if h.ForwardAgent == nil then
        h.ForwardAgent = "yes"
    end
    ~~~

* `registry`: Gets a current registry object.

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
* `question`: [kohkimakimoto/question](https://github.com/kohkimakimoto/gluaquestion).
* `template`: [kohkimakimoto/gluatemplate](https://github.com/kohkimakimoto/gluatemplate).
* `env`: [kohkimakimoto/gluaenv](https://github.com/kohkimakimoto/gluaenv).
* `http`: [cjoudrey/gluahttp](https://github.com/cjoudrey/gluahttp).
* `re`: [yuin/gluare](https://github.com/yuin/gluare)
* `sh`:[otm/gluash](https://github.com/otm/gluash)

## Predefined Variables

Essh provides predefined variables. In the recent version of Essh, there is one predefined varilable: `essh`.

`essh` is a table that has some functions and variables. see below

* `ssh_config` (string): ssh_config is ssh_config file path. At default, it is a temporary file that is generated automatically when you run Essh. You can overwrite this value for generating ssh_config to a static destination. If you use a gateway host that is a server between your client computer and a target server, you may use this variable to specify `ProxyCommand`. See below example:

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
