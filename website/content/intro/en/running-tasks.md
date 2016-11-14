+++
title = "Running Tasks"
type = "docs"
category = "intro"
lang = "en"
basename = "running-tasks.html"
+++

# Running Tasks

You can define tasks that are executed on remote and local servers.
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
[web01.localhost] hello on web01.localhost
[web02.localhost] hello on web02.localhost
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
[Local] hello on your-hostname
~~~

For more information on tasks, see the [Tasks](/docs/en/tasks.html) section.

Let's read next section: [Using Lua Libraries](using-lua-libraries.html)
