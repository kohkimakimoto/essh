# Running Tasks

You can define tasks that are executed on remote and local servers.
For example, edit your `esshconfig.lua`.

```lua
task "hello" {
    description = "say hello",
    prefix = true,
    backend = "remote",
    targets = "web",
    script = [=[
        echo "hello on $(hostname)"
    ]=],
}
```

Run the task.

```
$ essh hello
[web01.localhost] hello on web01.localhost
[web02.localhost] hello on web02.localhost
```

If you set `local` at `backend` property, Essh runs a task locally.

```lua
task "hello" {
    description = "say hello",
    prefix = true,
    backend = "local",
    script = [=[
        echo "hello on $(hostname)"
    ]=],
}
```

```
$ essh hello
[Local] hello on your-hostname
```

For more information on tasks, see the [Tasks](tasks.md) section.

Let's read next section: [Using Lua Libraries](getting-started_using-lua-libraries.md)
