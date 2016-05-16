# Essh Bash Module

Collection of the bash code for Essh script.

## Usage

```lua
local bash = essh.require "github.com/kohkimakimoto/essh/modules/bash"

task "example" {
    script = {
        bash.indent,
        "echo hello | indent",
    },
}
```

## Properties

* `indent` (string): Bash script code defines a `indent` function.

* `prefix` (string): Bash script code defines a `prefix` function.

* `upper` (string):

* `xterm` (string):

* `errexit_on` (string):

* `version` (string):

* `lock` (string):

* `driver` (string): Custom driver engine.

    ```lua
    driver "bash" {
        engine = bash.driver,
        show_version = true,
    }
    ```
