# essh bash module

Collection of the bash code for essh script.

## Usage

```lua
local bash = essh.require "github.com/kohkimakimoto/essh//modules/bash"

task "example" {
    script = {
        bash.indent,
        "echo hello | indent",
    },
}
```
