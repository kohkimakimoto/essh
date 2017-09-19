+++
title = "Zsh補完 | イントロダクション"
type = "docs"
category = "intro"
lang = "ja"
basename = "zsh-completion.html"
+++

# Zsh補完

EsshはSSHホストをリストするzsh補完をサポートしています。使用する場合は次のコードを`~/.zshrc`に追加してください。

~~~
eval "$(essh --zsh-completion)"
~~~

それから`.esshconfig.lua`を編集してください。`description`プロパティを次のように追加してみてください。

~~~lua
host "web01.localhost" {
    ForwardAgent = "yes",
    HostName = "192.168.0.11",
    Port = "22",
    User = "kohkimakimoto",
    -- add description
    description = "web01 development server",
}

host "web02.localhost" {
    ForwardAgent = "yes",
    HostName = "192.168.0.12",
    Port = "22",
    User = "kohkimakimoto",
    -- add description
    description = "web02 development server",
}
~~~

ホストが補完されるはずです。

~~~
$ essh [TAB]
web01.localhost  -- web01 development server
web02.localhost  -- web02 development server
~~~

`hidden`プロパティを使ってホストを隠すことができます。これをtrueに設定すると、zsh補完はホストを表示しません。

~~~lua
host "web01.localhost" {
    ForwardAgent = "yes",
    HostName = "192.168.0.11",
    Port = "22",
    User = "kohkimakimoto",
    description = "web01 development server",
    hidden = true,
}
~~~

`description`と`hidden`は最初の文字が小文字であることにお気づきですか。その他は大文字です。これは重要なポイントです。 Esshは一時ファイルに生成されるssh_configに最初の文字が大文字のプロパティを使用します。最初の文字が小文字のプロパティはssh_configではなく、Esshの機能の特殊な目的に使用されます。

> zshではなくbashを使用している場合は、`eval "$(essh --bash-completion)"`を使用できます。 だたしbash補完は説明の表示はサポートしていません。

ホストの詳細については、[ホスト](/essh/docs/ja/hosts.html) セクションを参照してください。

次のセクションに進みましょう: [フックを使う](using-hooks.html)
