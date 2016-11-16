+++
title = "Using Lua Libraries"
type = "docs"
category = "intro"
lang = "en"
basename = "using-lua-libraries.html"
+++

# Using Lua Libraries

Essh uses Lua for configuration and also has several built-in Lua libraries. You can use `require` function to load the libraries.

Example:

~~~lua
local question = require "question"

task "example" {
    prepare = function ()
        local r = question.ask("Are you OK? [y/N]: ")
        if r ~= "y" then
            -- return false, the task does not run.
            return false
        end
    end,
    script = [=[
        echo "foo"
    ]=],
}
~~~

`question` is a built-in library of Essh, that is implemented by [gluaquestion](https://github.com/kohkimakimoto/gluaquestion). It provides functions to get user input from a terminal.
And task's property `prepare` is a configuration that defines a function executed when the task starts.

When you run the task Essh, displays a message and waits your input.

~~~
$ essh example
Are you OK? [y/N]: y
foo
~~~

For more information on Lua libraries, see the [Lua VM](/docs/en/lua-vm.html) section.

Let's read next section: [Using Modules](using-modules.html)
