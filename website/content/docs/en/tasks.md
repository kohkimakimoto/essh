+++
title = "Tasks | Documentation"
type = "docs"
category = "docs"
lang = "en"
basename = "tasks.html"
+++

# Tasks

Task is a script that runs on remote or local servers. You can use it to automate your system administration tasks.

Example:

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

You can run a task below command.

~~~
$ essh example
~~~

Notice: Task name mustn't be duplicated with any host names.

## Properties

* `description` (string): Description of the task.

* `pty` (boolean): If it is true, SSH connection allocates pseudo-terminal by running ssh command with multiple -t options like `ssh -t -t`.

* `driver` (string): driver name is used in the task. see [Drivers](drivers.html).

* `parallel` (boolean): If it is true, runs task's script in parallel.

* `privileged` (boolean): If it is true, runs task's script by privileged user. If you use it, you have to configure your machine to be able to be used `sudo` without password.

* `disabled` (boolean): If it is true, this task does not run and is not displayed in tasks list.

* `hidden` (boolean): If it is true, this task is not displayed in tasks list.

* `targets` (string|table): Host names or tags that the task's scripts is executed for. You can use only hosts and tags which defined by same configuration registry of the task. For example, if you define a task in `/var/tmp/esshconfig.lua`, this task can not use hosts defined in `~/.essh/config.lua`. The first configuration file is **local** registry. But the second configuration file is **global** registry.

* `filters` (string|table): Host names or tags to filter target hosts. This property must be used with `targets`.

* `backend` (string): A place where the task's scripts will be executed on. You can set value only `remote` or `local`.

* `prefix` (boolean|string): If it is true, Essh displays task's output with hostname prefix. If it is string, Essh displays task's output with custom prefix. This string can be used with text/template format like `{{.Host.Name}}`.

* `prepare` (function): Prepare is a function to be executed when the task starts. See example:

    ~~~lua
    prepare = function ()
        -- cancel the task execution by returns false.
        return false
    end,
    ~~~

    By the prepare function returns false, you can cancel to execute the task.

* `script` (string|table): Script is code that will be executed. Example:

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

  * `ESSH_TASK_NAME`: task name.

  * `ESSH_SSH_CONFIG`: generated ssh_config file path.

  * `ESSH_DEBUG`: If you set `--debug` option by CLI. this variable is set "1".

  * `ESSH_HOSTNAME`: host name.

  * `ESSH_HOST_HOSTNAME`: host name.

  * `ESSH_HOST_SSH_{SSH_CONFIG_KEY}`: ssh_config key/value pare.

  * `ESSH_HOST_TAGS_{TAG}`: tag.

  * `ESSH_HOST_PROPS_{KEY}`: property that is set by host's props. See [Hosts](hosts.html).
