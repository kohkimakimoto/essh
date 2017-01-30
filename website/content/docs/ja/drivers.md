+++
title = "ドライバ | ドキュメント"
type = "docs"
category = "docs"
lang = "ja"
basename = "drivers.html"
+++

# ドライバ

Esshのドライバとは、タスク実行時にシェルスクリプトを構築するためのテンプレートシステムです。ドライバを使用してタスクの動作を変更することができます。

## Example

~~~lua
-- defining a driver
driver "custom_driver" { 
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

-- using defined driver in a task
task "example" {
    driver = "custom_driver",
    script = {
        "echo aaa",
        "echo bbb",
    }
}
~~~

## Environment template

Essh provides environment template to generate bash code to set environment variables.
You can used it as `{{template "environment" .}}`.

## Predefined variables

You can use predefined variables in the driver engine text template.

* `.Scripts`: This is a task's `script` value.

## Default driver 

If you define `default` driver like the following. This driver is used at default in the task instead of built-in default driver.

~~~lua
driver "default" { 
    engine = [=[
    -- your driver code...
    ]=],
}

-- This task uses above default driver automatically.
task "example" {
    script = {
        "echo aaa",
        "echo bbb",
    }
}
~~~
