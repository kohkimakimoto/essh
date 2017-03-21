# Essh Bash Module

Collection of the bash code for Essh script. This module is just example for developers who want to create own modules.

## Usage 

```lua
local bash = import "github.com/kohkimakimoto/essh/modules/bash"

task "example" {
    script = {
        bash.indent,
        "echo hello | indent",
    },
}
```
