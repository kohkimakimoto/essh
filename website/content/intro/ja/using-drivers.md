+++
title = "ドライバを使う | イントロダクション"
type = "docs"
category = "intro"
lang = "ja"
basename = "using-drivers.html"
+++

# ドライバを使う

Esshのドライバとは、タスク実行時にシェルスクリプトを構築するためのテンプレートシステムです。
このチュートリアルで、あなたはすでにEsshバイナリに含まれているデフォルトの組み込みドライバを使用しています。カスタムドライバを使用してタスクの動作を変更することができます。

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
...
[essh debug] run task: example
[essh debug] driver: default 
[essh debug] real local command: [bash -c 
export ESSH_TASK_NAME='example'
export ESSH_SSH_CONFIG=/var/folders/bt/xwh9qmcj00dctz53_rxclgtr0000gn/T/essh.ssh_config.767200705
export ESSH_DEBUG="1"

echo aaa
echo bbb
]
~~~

デバッグメッセージによると、タスクは以下のbashスクリプトを実行しました：

~~~
export ESSH_TASK_NAME='example'
export ESSH_SSH_CONFIG=/var/folders/bt/xwh9qmcj00dctz53_rxclgtr0000gn/T/essh.ssh_config.767200705
export ESSH_DEBUG="1"

echo aaa
echo bbb
~~~

この内容は**組み込みドライバ**によって生成されたものです。組み込みドライバは、Esshバイナリに含まれている次のテキストテンプレートです。

~~~
{{template "environment" .}}
{{range $i, $script := .Scripts}}{{$script.code}}
{{end}}
~~~

`{{template "environment" .}}` は環境変数を生成します。上の例でこの部分は以下のコードになります。

~~~
export ESSH_TASK_NAME='example'
export ESSH_SSH_CONFIG=/var/folders/bt/xwh9qmcj00dctz53_rxclgtr0000gn/T/essh.ssh_config.767200705
export ESSH_DEBUG="1"
~~~

その後、Esshはスクリプトテキストを改行コードで連結します。

~~~
{{range $i, $script := .Scripts}}{{$script.code}}
{{end}}
~~~


上記のコードは次のようになります:

~~~
echo aaa
echo bbb
~~~

それでは `driver` 関数を使って最初のカスタムドライバを定義してみましょう。

~~~lua
driver "my_driver" {
    engine = [=[
        {{template "environment" .}}
        
        set -e
        indent() {
            local n="${1:-4}"
            local p=""
            for i in `seq 1 $n`; do
                p="$p "
            done;

            local c="s/^/$p/"
            case $(uname) in
              Darwin) sed -l "$c";;
              *)      sed -u "$c";;
            esac
        }
        
        {{range $i, $script := .Scripts -}}
        echo '==> step {{$i}}:{{if $script.description}} {{$script.description}}{{end}}'
        { 
            {{$script.code}} 
        } | indent; __essh_exit_status=${PIPESTATUS[0]}
        if [ $__essh_exit_status -ne 0 ]; then
            exit $__essh_exit_status
        fi
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

`driver` 関数には、必須パラメータ `engine` が必要です。これがテンプレートテキストです。 カスタムドライバを使用するには、タスクの`driver`プロパティを設定する必要があります。

この例のドライバはステップ番号と説明、インデントされたスクリプトの標準出力を表示します。上記のタスクを実行すると、次の出力が得られます。


~~~
==> step 0:
    aaa
==> step 1:
    bbb
~~~

説明はまだ表示されていません。各スクリプトのコードに`description`プロパティを設定してください。


~~~
task "example" {
    driver = "my_driver",
    script = {
        {
            description = "output aaa",
            code = "echo aaa",
        },
        {
            description = "output bbb",
            code = "echo bbb",
        },
    }
}
~~~

以下のような出力なります。

~~~
==> step 0: output aaa
    aaa
==> step 1: output bbb
    bbb
~~~

ドライバの詳細については、[ドライバ](/docs/ja/drivers.html)のセクションを参照してください。

次のセクションに進みましょう: [ジョブを定義する](defining-jobs.html)