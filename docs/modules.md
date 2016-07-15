# Modules

Module allows you to use, create and share reusable code easily for Essh configuration.

## Using Modules

You can use `essh.require` function to load a Essh module.

```lua
local bash = essh.require "github.com/kohkimakimoto/essh/modules/bash"
```

`essh.require` returns Lua value. In the above case, `bash` is Lua table that has several variables and functions. You can use `bash` in your configuration.

```lua
local bash = essh.require "github.com/kohkimakimoto/essh//modules/bash"

task "example" {
    script = {
        bash.indent,
        "echo hello | indent",
    },
}
```

`bash.indent` is [this code snippet](https://github.com/kohkimakimoto/essh/blob/master/modules%2Fbash%2Findex.lua#L3-L17).
So the task displays indented output.

`essh.require` is implemented by using [hashicorp/go-getter](https://github.com/hashicorp/go-getter). You can use git url and local filesystem path to specify a module path.

Modules are installed automatically when Essh runs. The installed modules are stored in `.essh` directory. If you need to update installed modules, runs `essh --update`.

```
$ essh --update
```

### Creating Modules

Creating new modules is easy. A minimum module is a directory that includes only `index.lua`.
Try to create `my_module` directory and `index.lua` file in the directory.

```lua
-- my_module/index.lua
local m = {}

m.hello = "echo hello"

return m
```

`index.lua` is the entry-point that have to return Lua value. This example returns a table that has `hello` variable. That's it. To use this module, write below config.

```lua
local my_module = essh.require "./my_module"

task "example" {
    script = {
        my_module.hello,
    },
}
```

Run it.

```
$ essh example
hello
```

If you want to share the module, create a git repository from the module directory and push it to a remote repository as github.com. To use the module of git repository, you update `essh.require` path to the url.

```lua
local my_module = essh.require "github.com/your_account/my_module"

task "example" {
    script = {
        my_module.hello,
    },
}
```
