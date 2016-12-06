+++
title = "Modules | Documentation"
type = "docs"
category = "docs"
lang = "en"
basename = "modules.html"
+++

# Modules

Module allows you to use, create and share reusable code easily for Essh configuration.

## Using Modules

You can use `import` function to load a Essh module.

~~~lua
local bash = import "github.com/kohkimakimoto/essh/modules/bash"
~~~

`import` returns Lua value. In the above case, `bash` is Lua table that has several variables and functions.

~~~lua
local bash = import "github.com/kohkimakimoto/essh//modules/bash"

task "example" {
    script = {
        bash.indent,
        "echo hello | indent",
    },
}
~~~

`bash.indent` is [this code snippet](https://github.com/kohkimakimoto/essh/blob/master/modules%2Fbash%2Findex.lua#L3-L17).
So the task displays indented output.

`import` is implemented by using [hashicorp/go-getter](https://github.com/hashicorp/go-getter). You can use git url and local filesystem path to specify a module path. For example:

* Getting a module from a github repository:

    ~~~
    local mod = import "github.com/username/repository"
    ~~~

* Getting a module from a github repository and checkout tag, commit or branch:

    ~~~
    local mod = import "github.com/username/repository?ref=master"
    ~~~

* Getting a module from a github repository's subdirectory:

    ~~~
    local mod = import "github.com/username/repository//path/to/module"
    ~~~

    The double-slash, `//` is the separator for a subdirectory, and not part of the repository itself.

* Getting a module from a generic git repository:
    
    ~~~~
    local mod = import "git::ssh://your-private-git-server/path/to/repo.git"
    ~~~~

* Getting a module from a local filesystem:

    ~~~
    local mod = import "/path/to/module"
    ~~~

For detail, see [hashicorp/go-getter](https://github.com/hashicorp/go-getter).

Modules are installed automatically when Essh runs. The installed modules are stored in `.essh/modules` directory if the config file that is written `import` is `esshconfig.lua` in the current directory.
If the config file is `~/.essh/config.lua`, the modules are stored in `~/.essh/modules` directory.

If you need to update installed modules, runs `essh --update`.

~~~
$ essh --update
~~~

At default, Essh updates only modules in `.essh/modules`. If you want to update `~/.essh/modules` , run the following command:

~~~
$ essh --update --with-global
~~~

### Creating Modules

Creating new modules is easy. A minimum module is a directory that includes only `index.lua`.
Try to create `my_module` directory and `index.lua` file in the directory.

~~~lua
-- my_module/index.lua
local m = {}

m.hello = "echo hello"

return m
~~~

`index.lua` is the entry-point that have to return Lua value. This example returns a table that has `hello` variable. That's it. To use this module, write below config.

~~~lua
local my_module = import "./my_module"

task "example" {
    script = {
        my_module.hello,
    },
}
~~~

Run it.

~~~
$ essh example
hello
~~~

If you want to share the module, create a git repository from the module directory and push it to a remote repository as github.com. To use the module of git repository, you update `import` path to the url.

~~~lua
local my_module = import "github.com/your_account/my_module"

task "example" {
    script = {
        my_module.hello,
    },
}
~~~
