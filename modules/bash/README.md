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
