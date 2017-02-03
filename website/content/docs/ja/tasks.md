+++
title = "タスク | ドキュメント"
type = "docs"
category = "docs"
lang = "ja"
basename = "tasks.html"
+++

# タスク

タスクは、リモートサーバーまたはローカルサーバーで実行されるスクリプトです。これを使用してシステム管理タスクを自動化することができます。

## 例

~~~lua
task "example" {
    description = "example task",
    
    targets = {
        "web",
    },
    
    filters = {
        "production",
    },
    
    backend = "local",
    
    parallel = true,
    
    prefix = true,
    
    script = {
        "echo foo",
        "echo bar"
    },
}
~~~

注意：タスク名はホスト名と重複して定義することはできません。

以下のコマンドでタスクを実行できます。

~~~
$ essh example
~~~

タスクに引数を渡すことができます。

~~~
$ essh example foo bar
~~~



## プロパティ

* `description` (string): タスクの説明。

* `pty` (boolean): trueに設定すると、SSH接続は`ssh -t -t`のように複数の-tオプションを使用してsshコマンドに擬似端末を割り当てます。

* `driver` (string): このタスクで使用するドライバ。[Drivers](drivers.html)を参照してください。

* `parallel` (boolean): trueに設定すると、タスクのスクリプトを並列に実行します。

* `privileged` (boolean): trueに設定すると、特権ユーザーがタスクのスクリプトを実行します。これを使用する場合は、パスワードなしでsudoを使用できるようにマシンを設定する必要があります。

* `hidden` (boolean): trueの場合、このタスクはタスクリストに表示されません。

* `targets` (string|table): タスクのスクリプトが実行されるホスト名またはタグ。

* `filters` (string|table): ターゲットホストをフィルタリングするためのホスト名またはタグ。このプロパティは`targets`と一緒に使わなければなりません。

* `backend` (string): タスクのスクリプトが実行される場所。`remote`か`local`を指定できます。

* `prefix` (boolean|string): trueの場合、Esshはタスクの出力にホスト名のプレフィックスをつけて表示します。文字列の場合、Esshはタスクの出力にカスタムのプレフィックスをつけて表示します。この文字列は `{{.Host.Name}}`のようなテキスト/テン​​プレート形式で使用できます。

* `prepare` (function): Prepareは、タスクの開始時に実行される関数です。例を参照してください:

    ~~~lua
    prepare = function (t)
        -- override task config
        t.targets = "web"
        -- get command line arguments
        print(t.args[1])
        -- cancel the task execution by returns false.
        return false
    end,
    ~~~

    prepare関数によってfalseが返されると、タスクのスクリプトの実行を取り消すことができます。

* `props` (table): Propsはタスクが実行されるとき、環境変数を`ESSH_TASK_PROPS_${KEY}=VALUE`で設定します。テーブルキーは大文字に変更されます。

    ~~~lua
    props = {
        foo = "bar",
    }

    -- export ESSH_TASK_PROPS_FOO="bar"
    ~~~
    
* `script` (string|table): 実行されるコードです。例:

    ~~~lua
    script = [=[
        echo aaa
        echo bbb
        echo ccc
    ]=]
    ~~~

    or

    ~~~lua
    script = {
        "echo aaa",
        "echo bbb",
        "echo ccc",
    }
    ~~~

    テーブルで設定すると、Esshはテーブルの文字列を改行コードで連結します。 Esshはスクリプトをbashスクリプトとして実行します。しかし、これはただのデフォルトの動作です。[ドライバ](drivers.html)で変更できます。

    スクリプトでは、定義済みの環境変数を使用できます。以下を参照してください。

  * `ESSH_TASK_NAME`: タスク名。

  * `ESSH_SSH_CONFIG`: 生成されたssh_configファイルのパス。

  * `ESSH_DEBUG`: CLIで `--debug`オプションを設定した場合この変数は"1"に設定されます。

  * `ESSH_TASK_PROPS_${KEY}`: タスクの`props`によって設定される値。
  
  * `ESSH_TASK_ARGS_${INDEX}`: コマンドライン引数によって渡される引数の値。インデックスは '1'から始まります。

  * `ESSH_HOSTNAME`: ホスト名。

  * `ESSH_HOST_HOSTNAME`: ホスト名。

  * `ESSH_HOST_SSH_{SSH_CONFIG_KEY}`: ssh_configのキー/バリュー ペア.

  * `ESSH_HOST_TAGS_{TAG}`: タグ。タグを設定すると、この変数の値は"1"になります。

  * `ESSH_HOST_PROPS_{KEY}`: ホストの`props`によって設定される値。[ホスト](hosts.html)を参照してください。

  * `ESSH_JOB_NAME`: ネームスペース名。[ネームスペース](namespaces.html)を参照してください。
  
* `script_file` (string): ファイルパスまたはhttpまたはhttpsでアクセスできるURL。ファイルの内容が実行されます。 `script_file`と` script`を同時に使うことはできません。
