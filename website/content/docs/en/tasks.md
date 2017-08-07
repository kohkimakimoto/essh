+++
title = "Tasks | Documentation"
type = "docs"
category = "docs"
lang = "en"
basename = "tasks.html"
+++

# Tasks

Task is a script that runs on remote or local servers. You can use it to automate your system administration tasks.

## Example

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

Notice: Task name mustn't be duplicated with any host names.

You can run a task below command.

~~~
$ essh example
~~~

You can pass arguments to a task.

~~~
$ essh example foo bar
~~~


## Properties

* `description` (string): Description of the task.

* `pty` (boolean): If it is true, SSH connection allocates pseudo-terminal by running ssh command with multiple -t options like `ssh -t -t`.

* `driver` (string): driver name is used in the task. see [Drivers](drivers.html).

* `parallel` (boolean): If it is true, runs task's script in parallel.

* `privileged` (boolean): If it is true, runs task's script by privileged user. If you use it, you have to configure your machine to be able to be used `sudo` without password.

* `user` (string): Runs task's script by specific user. If you use it, you have to configure your machine to be able to be used `sudo` without password.

* `hidden` (boolean): If it is true, this task is not displayed in tasks list.

* `targets` (string|table): Host names or tags that the task's scripts is executed for.

* `filters` (string|table): Host names or tags to filter target hosts. This property must be used with `targets`.

* `backend` (string): A place where the task's scripts will be executed on. You can set value only `remote` or `local`.

* `prefix` (boolean|string): If it is true, Essh displays task's output with hostname prefix. If it is string, Essh displays task's output with custom prefix. This string can be used with text/template format like `{{.Host.Name}}`.

* `prepare` (function): Prepare is a function to be executed when the task starts. See example:

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

    By the prepare function returns false, you can cancel to execute the task's script.

* `props` (table): Props sets environment variables `ESSH_TASK_PROPS_${KEY}=VALUE` when the task is executed. The table key is modified to upper cased.

    ~~~lua
    props = {
        foo = "bar",
    }

    -- export ESSH_TASK_PROPS_FOO="bar"
    ~~~

* `script` (string|table): Code that will be executed. Example:

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

    If you set it as a table, Essh concatenates strings in the table with newline code. And Essh runs the script as a bash script. But this is just default behavior. You can change it by [Drivers](drivers.html).

    You can use predefined environment variables in your script, See below:

  * `ESSH_TASK_NAME`: Task name.

  * `ESSH_SSH_CONFIG`: Generated ssh_config file path.

  * `ESSH_DEBUG`: If you set `--debug` option by CLI. this variable is set "1".

  * `ESSH_TASK_PROPS_${KEY}`: The value that is set by task's `props`.
  
  * `ESSH_TASK_ARGS_${INDEX}`: The argument's value that is passed by a command line arguments. The index starts at '1'.

  * `ESSH_HOSTNAME`: Host name.

  * `ESSH_HOST_HOSTNAME`: Host name.

  * `ESSH_HOST_SSH_{SSH_CONFIG_KEY}`: ssh_config key/value pare.

  * `ESSH_HOST_TAGS_{TAG}`: Tag. If you set a tag, This variable has a value "1".

  * `ESSH_HOST_PROPS_{KEY}`: The value that is set by host's `props`. See [Hosts](hosts.html).

  * `ESSH_NAMESPACE_NAME`: Namespace name. See [Namespaces](namespaces.html).
  
* `script_file` (string): A file path or URL that can be accessed by http or https. The file's content will be executed. You can't use `script_file` and `script` at the same time.