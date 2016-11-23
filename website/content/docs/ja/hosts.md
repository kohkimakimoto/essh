+++
title = "ホスト | ドキュメント"
type = "docs"
category = "docs"
lang = "ja"
basename = "hosts.html"
+++

# ホスト

Esshのホストとは、あなたが管理するSSHサーバです。ホスト設定を使用して、Esshは動的にSSHコンフィグを生成し、フック機能を提供し、タグでホストを分類します。

例:

~~~lua
host "web01.localhost" {
    -- SSH config properties.
    HostName = "192.168.0.11",
    Port = "22",
    User = "kohkimakimoto",

    -- Essh config properties.
    description = "web01 development server",
    hidden = false,
    private = false,
    props = {},
    tags = {},
    hooks_before_connect = {},
    hooks_after_connect = {},
    hooks_after_disconnect = {},
}
~~~

ホストは2種類のプロパティで構成されています。 **SSHコンフィグプロパティ** と **Esshコンフィグプロパティ** です。

## SSHコンフィグプロパティ

SSHコンフィグプロパティは、最初の文字を大文字にする必要があります。
例えば​​`HostName`や`Port`です。 このタイプのプロパティは**ssh_config**を生成するために使用されます。このプロパティはすべてのsshオプションを使用できます。ssh_config(5)を参照してください。

## Esshコンフィグプロパティ

Esshコンフィグプロパティは、最初の文字を小文字にする必要があります。
これらはssh_configではなくEsshの機能の特殊な目的に使用されます。

このタイプのプロパティのすべてを以下に記載します。

* `description` (string): ホストの説明。これは、ホストのリストとzsh補完に表示するために使用されます。

* `hidden` (boolean): trueに設定すると、zsh補完はこのホストを表示しません。

* `private` (boolean): trueに設定すると、このホストはタスクに対してのみ使用できます。[スコープ](#scope)も参照してください。

* `hooks_before_connect` (table): 接続する前に発火するフック。これはローカルで実行されます。フックはLuaテーブルとして定義されています。このテーブルは、複数の関数または文字列を持つことができます。例を参照してください:

    ~~~lua
    hooks_before_connect = {
        -- function
        function()
            print("foo")
        end,

        -- string (commands)
        "echo bar",

        -- If the function returns a string, Essh run the string as a command.
        function()
            return "echo foobar"
        end,
    }
    ~~~

    Lua関数で実装されたすべてのフック(`hooks_after_connect`, `hooks_after_disconnect`も含む)はローカルで実行されます。

    すべてのフック(`hooks_after_connect`, `hooks_after_disconnect`も含む)は、単にsshでログインしたときにのみ発火します。フックはタスクや`--exec`オプションで発火しません。

* `hooks_after_connect` (table): 接続後に発火するフック。これはリモートサーバで実行されます。

* `hooks_after_disconnect` (table): 切断後に発火するフック。これはローカルで実行されます。

* `tags` (array table): タグはホストを分類します。

    ~~~lua
    tags = {
        "web",
        "production",
    }
    ~~~

    タグをホスト名と重複させることはできません。

* `props` (table): Propsはホストがタスクで使用されるときの環境変数を`ESSH_HOST_PROPS_{KEY}`で設定します。テーブルキーは大文字に変更されます。

    ~~~lua
    props = {
        foo = "bar",
    }

    -- ESSH_HOST_PROPS_FOO=bar
    ~~~

## スコープ {#scope}

Esshのホストには、**プライベート** または **パブリック** のスコープがあります。

プライベートホストは同じレジストリ内で定義されたタスクにしか使用できません。このタイプのホストはsshログインには使用できません。また、zsh補完ではサジェストされません。`--exec`オプションでも使用できません。

`private_host`関数をエイリアスとして使用して、プライベートホストを定義することができます。以下の例を参照してください。

~~~lua
private_host "example" {
    HostName = "192.168.0.11",
    Port = "22",
    User = "kohkimakimoto",
}
~~~

これは以下と同じです。

~~~lua
host "example" {
    HostName = "192.168.0.11",
    Port = "22",
    User = "kohkimakimoto",
    private = true,
}
~~~

## 制約

[スコープ](#scope)と[レジストリ](configuration-files.html#registries)に関する制約があります。これは主にプロジェクト固有のタスクとホストをグローバルにSSH接続に使用したいホストと混在させないために、役に立ちます。

* **各パブリックホストは一意でなければなりません。**(パブリックホストはローカルレジストリとグローバルレジストリで同じ名前で定義することはできません)
* **すべてのホストは、同じレジストリ内で一意でなければなりません。**(同じレジストリ内で同じ名前でホストを定義することはできません)
* **タスクによって使用されるホストは、同じレジストリに定義する必要があります。**(タスクは、同じレジストリで定義されたホストのみを参照できます)
* **プライベートホストはタスクのみで使用できます。**
