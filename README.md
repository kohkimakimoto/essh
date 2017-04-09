# Essh [![Build Status](https://travis-ci.org/kohkimakimoto/essh.svg?branch=master)](https://travis-ci.org/kohkimakimoto/essh)

Extended ssh command. 

* [Website](https://essh.sitespread.net/)
* [Documentation](https://essh.sitespread.net/docs/en/index.html)
* [Gettting Started](https://essh.sitespread.net/intro/en/index.html)

## Overview

Essh is an extended `ssh` command. If you use `essh` command instead of `ssh`, Your SSH operation becomes more efficient and convenient. Essh is a single binary CLI tool and simply wraps ssh command. You can use it in the same way as ssh. And it has useful features over ssh.

## Features

* **Configuration As Code**: You can write SSH client configuration (aka:`~/.ssh/config`) in [Lua](https://www.lua.org/) code. So your ssh_config can become more dynamic.

* **Hooks**: Essh supports hooks that execute commands when it connects a remote server.

* **Servers List Management**: Essh provides utilities for managing hosts, that list and classify servers by using tags.

* **Per-Project Configuration**: Essh supports per-project configuration. This allows you to change SSH hosts config by changing current working directory.

* **Task Runner**: Task is a script that runs on remote and local servers. You can use it to automate your system administration tasks.

* **Modules**: Essh provides modular system that allows you to use, create and share reusable Lua code easily.

## Developing

Requirements

* Go 1.7 or later (my development env)
* [Glide](https://glide.sh/)
* [gox](https://github.com/mitchellh/gox)

Installing dependences

```
$ make deps
```

Building dev binary.

```
$ make dev
```

Building distributed binaries.


```
$ make dist
```

Building packages (now support only RPM)

```
$ make dist
$ make packaging
```

## Author

Kohki Makimoto <kohki.makimoto@gmail.com>

## License

The MIT License (MIT)
