# Using Hooks

Hooks in Essh are scripts executed before and after connecting remote servers.

Example:

```lua
host "web01.localhost" {
    HostName = "192.168.0.11",
    Port = "22",
    User = "kohkimakimoto",

    hooks = {
        before_connect = "echo before_connect",
        after_connect = "echo after_connect",
        after_disconnect = "echo after_disconnect",
    },
}
```

Essh supports below type of hooks:

* `before_connect` (string or table or function): fires on the localhost before you connect a server via SSH.

* `after_connect` (string or table or function): fires on the remote host after you connect a server via SSH.

* `after_disconnect` (string or table or function): fires on the local host after you disconnect from a SSH server.

> Note: I am using this functionality to change OSX terminal profile(color). See the below example.

```lua
host "web01.localhost" {
    -- ...
    hooks = {
        before_connect = "osascript -e 'tell application \"Terminal\" to set current settings of first window to settings set \"Blue Profile\"'",
        after_disconnect = "osascript -e 'tell application \"Terminal\" to set current settings of first window to settings set \"Normal Profile\"'",
    },
}
```

Let's read next section: [Managing Hosts](getting-started_managing-hosts.md).
