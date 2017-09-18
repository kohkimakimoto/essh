+++
title = "Running Tasks | Introduction"
type = "docs"
category = "intro"
lang = "en"
basename = "running-tasks.html"
+++

# Running Tasks

Task is a script that runs on remote servers or local.
For example, edit your `esshconfig.lua`.

~~~lua
task "hello" {
    description = "say hello",
    prefix = true,
    backend = "remote",
    targets = "web",
    script = [=[
        echo "hello on $(hostname)"
    ]=],
}
~~~

Run the task.

~~~
$ essh hello
[remote:web01.localhost] hello on web01.localhost
[remote:web02.localhost] hello on web02.localhost
~~~

If you set `local` at `backend` property, Essh runs a task locally.

~~~lua
task "hello" {
    description = "say hello",
    prefix = true,
    backend = "local",
    script = [=[
        echo "hello on $(hostname)"
    ]=],
}
~~~

~~~
$ essh hello
[local] hello on your-hostname
~~~

For more information on tasks, see the [Tasks](/essh/docs/en/tasks.html) section.

Let's read next section: [Using Lua Libraries](using-lua-libraries.html)
