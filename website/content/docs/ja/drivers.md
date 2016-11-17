+++
title = "ドライバ | ドキュメント"
type = "docs"
category = "docs"
lang = "ja"
basename = "drivers.html"
+++

# ドライバ

ドライバは、タスクの実行時にスクリプトを構築するためのテンプレートです。
タスク設定でドライバを指定しない場合、Esshはデフォルトのビルトインドライバを使用します。

ドライバが何をするのかを理解するために、次の短い例を参照してください

~~~lua
task "example" {
    script = {
        "echo aaa",
        "echo bbb",
    }
}
~~~

`--debug`オプションをつけてこのタスクを実行して、実際のスクリプトを表示してみてください。

~~~
$ essh example --debug
[essh debug] ...
[essh debug] real local command: [bash -c
export ESSH_TASK_NAME='example'
export ESSH_SSH_CONFIG=/var/folders/bt/xwh9qmcj00dctz53_rxclgtr0000gn/T/essh.ssh_config.544434412

echo aaa
echo bbb
]
~~~

デバッグメッセージによれば、タスクは次のようなbashスクリプトを実行しました:

~~~
export ESSH_TASK_NAME='example'
export ESSH_SSH_CONFIG=/var/folders/bt/xwh9qmcj00dctz53_rxclgtr0000gn/T/essh.ssh_config.544434412

echo aaa
echo bbb
~~~

この内容は**ビルトインドライバ**によって生成されたものです。ビルトインドライバとは次のテキストテンプレートです。

~~~go
{{template "environment" .}}
{{range $i, $script := .Scripts}}{{$script.code}}
{{end}}
~~~

`{{template environment"}}`は環境変数を生成します。上の例でこのセクションは以下の部分を出力します。

~~~
export ESSH_TASK_NAME='example'
export ESSH_SSH_CONFIG=/var/folders/bt/xwh9qmcj00dctz53_rxclgtr0000gn/T/essh.ssh_config.544434412
~~~

その後、Esshは `script`テキストを改行コードで連結します。上記の例では、以下を出力します。

~~~
echo aaa
echo bbb
~~~

結論：ドライバはbashスクリプトを出力するためのテンプレートです。

## カスタムドライバ

`driver`関数を使ってカスタムドライバを定義して使うことができます。

例:

~~~lua
driver "my_driver" {
    engine = [=[
        {{template "environment" .}}
        {{range $i, $script := .Scripts}}{{$script.code}}
        {{end}}
    ]=],
}

task "example" {
    driver = "my_driver",
    script = {
        "echo aaa",
        "echo bbb",
    }
}
~~~

`driver`の設定には、必須パラメータ`engine`が必要です。これがテンプレートテキストです。
カスタムドライバを使用するには、タスクの`driver`プロパティを設定する必要があります。

サンプル実装の詳細は以下を参照してください:

[bash.driver](https://github.com/kohkimakimoto/essh/blob/master/modules/bash/index.lua)

or

[docker.driver](https://github.com/kohkimakimoto/essh/blob/master/modules/docker/index.lua)
