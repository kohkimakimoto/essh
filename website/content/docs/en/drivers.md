+++
title = "Drivers | Documentation"
type = "docs"
category = "docs"
lang = "en"
basename = "drivers.html"
+++

# Drivers

Drivers in Essh are template system to construct shell scripts in tasks execution. You can use a driver to modify behavior of tasks.

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

## Running scripts in other programming languages

You can run scripts in other programming languages by using a driver.

Why would you want to do this? There are a few reasons:

* You may want to use a language that is more familiar to you with more features than shell scripts.
* Although you can write in lua language the hooks functions, scripts always run in a shell environment. If you want to run a script in a different environment, you can use a driver.
* You may want to use a language that is more suitable for the task at hand.
* If the language script interpreter is not installed on the target machine, you can use a driver to install it.


For example, you can run a script in Python as follows.

Define the new driver:
~~~lua
driver "python" {
    engine = [=[
        {{template "environment" .}}

        # Detect if python or python3 is installed, python3 is preferred
        if ! command -v python3 &> /dev/null
        then
            if ! command -v python &> /dev/null
            then
                echo "Python could not be found. Please install it."
                exit 1
            else
                CMD=python
            fi
        else
                CMD=python3
        fi
        
        
        set -e

$CMD <<EOF
{{range $i, $script := .Scripts -}}
{{$script.code}} 
{{end -}}
EOF

        set +e
    ]=],
}
~~~

Use the driver in a task:

~~~lua
task "example" {
    driver = "python",
    script = [=[
# This is a python script, be careful with indentation
# Detect hostname and print it
import socket
hostname = socket.gethostname()
print(f"Hello world from Python in {hostname}")
    ]=]
}
~~~

### Other useful drivers

**Deno driver**

~~~lua
driver "deno" { 
    engine = [=[
        {{template "environment" .}}

        # Detect if deno is installed
        if ! command -v deno &> /dev/null
        then
            echo "Deno could not be found. Please install it."
            exit 1
        fi
        
        set -e
deno run --unstable -q -A - <<EOF
{{range $i, $script := .Scripts -}}
    {{$script.code}} 
{{end -}}

EOF
        set +e
    ]=],
}
~~~

Example:

~~~lua
task "example-deno-task" {
    driver = "deno",
    script = [=[
        // This is a deno script
        // Detect hostname and print it
        console.log("Hello world from Deno in", Deno.hostname());
    ]=]
}
~~~

**ESSH Lua driver (since *3.5.0*)**

~~~lua
driver "essh-lua" {
    engine = [=[
        {{template "environment" .}}

        # Detect if essh is installed
        if ! command -v essh &> /dev/null
        then
            echo "essh could not be found. Please install it."
            exit 1
        fi
        
        set -e

essh --eval <<EOF
{{range $i, $script := .Scripts -}}
{{$script.code}} 
{{end -}}
EOF

        set +e
    ]=],
}
~~~

Example:

~~~lua
task "example-essh-lua-task" {
    driver = "essh-lua",
    script = [=[
        -- This is a essh-lua script
        print("Hello world from Essh-Lua")

        local sh = require("sh")
        sh.ls("-l","/"):print()
    ]=]
}
~~~