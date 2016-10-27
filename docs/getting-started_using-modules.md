# Using Modules

Essh has a modular system that allows you to use reusable code easily for Essh configuration.
For instance, I provide [bash module](https://github.com/kohkimakimoto/essh/tree/master/modules/bash) that is a collection of bash scripts for using in your Essh tasks.
You can use `import` function to load a module.

Example:

```lua
local bash = import "github.com/kohkimakimoto/essh/modules/bash"

task "example" {
    script = {
        bash.version,
        "echo foo",
    },
}
```

`bash.version` is a variable that actually is a simple string `bash --version`. So this task prints bash version and then runs `echo foo`.

The modules are installed automatically, when you run Essh.
You run the task, you will get as below.

```
$ essh example
Installing module: 'github.com/kohkimakimoto/essh/modules/bash' (into /path/to/directory/.essh)
GNU bash, version 4.1.2(1)-release (x86_64-redhat-linux-gnu)
Copyright (C) 2009 Free Software Foundation, Inc.
License GPLv3+: GNU GPL version 3 or later <http://gnu.org/licenses/gpl.html>

This is free software; you are free to change and redistribute it.
There is NO WARRANTY, to the extent permitted by law.
foo
```

For more information on Modules, see the [Modules](modules.md) section.

Let's read next section: [Using Drivers](getting-started_using-drivers.md)
