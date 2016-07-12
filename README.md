# Essh

Essh is an extended `ssh` command. The features are the following:

* **Simple**: A single binary CLI tool. Essh simply wraps `ssh` command. You can use it in the same way as `ssh`.
* **Configuration As Code**: You can write SSH client configuration (aka:`~/.ssh/config`) in Lua code. So your ssh_config can become more dynamic.
* **Hooks**: Essh supports hooks that execute commands when it connects a remote server.
* **Servers List Management**: Essh provides utilities for managing hosts, that list and classify servers by using tags.
* **Zsh Completion**: Essh provides built-in zsh completion.
* **Per-Project Configuration**: Essh supports per-project configuration. This allows you to change SSH hosts config by changing current working directory.
* **Task Runner**: Task is a script that runs on remote and local servers. You can use it to automate your system administration tasks.
* **Modules**: Essh provides modular system that allows you to use, create and share reusable Lua code easily.

> **Now Essh is on unstable stage. API and code may be broken in future. And document lacks. sorry!**

**Demo1: Completing a hostname and changing terminal color by using hook.**

![optimized](essh-demo01.gif)

## Installation

Essh is provided as a single binary. You can download it and drop it in your $PATH.

[Download latest version](https://github.com/kohkimakimoto/essh/releases/latest)

## Documentation

TODO...

## Developing

Requirements

* Go 1.6 or later
* [Gom](https://github.com/mattn/gom)

Installing dependences

```
$ make deps
```

Building dev binary.

```
$ make
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
