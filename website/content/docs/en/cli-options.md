+++
title = "CLI Options | Documentation"
type = "docs"
category = "docs"
lang = "en"
basename = "cli-options.html"
+++

# CLI Options

In order to extend `ssh` command, Essh uses several command line options.
All these options are long option that names with the double-dash prefix to prevent conflict with `ssh` command options.

All the options are listed below.

## General

* `--print`: Print generated ssh_config.

* `--gen`: Only generate ssh_config.

* `--working-dir <dir>`: Change working directory.

* `--config <file>`: Load configuration from the file.

* `--color`: Force ANSI output.

* `--no-color`: Disable ANSI output.

* `--debug`: Output debug log.

## Manage Hosts, Tags And Tasks

* `--hosts`: List hosts.

* `--select <tag|host>`: (Using with `--hosts` option) Get only the hosts filtered with tags or hosts.

* `--filter <tag|host>`: (Using with `--hosts` option) Filter selected hosts with tags or hosts.

* `--namespace <namespace>`: (Using with `--hosts` option) Get hosts from specific namespace.

* `--ssh-config`: (Using with `--hosts` option) Output selected hosts as ssh_config format.

* `--tasks`: List tasks.

* `--all`: (Using with `--tasks` option) Show all that include hidden objects.

* `--tags`: List tags.

* `--namespaces`: List namespaces.

* `--quiet`: (Using with `--hosts`, `--tasks` or `--tags` option) Show only names.

## Manage Modules

* `--update`: Update modules.

* `--clean-modules`: Clean downloaded modules.

* `--clean-tmp`: Clean temporary data.

* `--clean-all`: Clean all data.

* `--with-global`: (Using with `--update`, `--clean-modules`, `--clean-tmp` or `--clean-all` option) Update or clean modules in the local and global both registry.

## Execute Commands

* `--exec`: Execute commands with the hosts.

* `--target <tag|host>`: (Using with `--exec` option) Target hosts to run the commands.

* `--filter <tag|host>`: (Using with `--exec` option) Filter target hosts with tags or hosts.

* `--backend remote|local`: (Using with `--exec` option) Run the commands on local or remote hosts.

* `--prefix`: (Using with `--exec` option) Enable outputing prefix.

* `--prefix-string <prefix>` (Using with `--exec` option) Custom string of the prefix.

* `--privileged`: (Using with `--exec` option) Run by the privileged user.

* `--parallel`: (Using with `--exec` option) Run in parallel.

* `--pty`: (Using with `--exec` option) Allocate pseudo-terminal. (add ssh option "-t -t" internally)

* `--script-file`: (Using with `--exec` option) Load commands from a file.

* `--driver`: (Using with `--exec` option) Specify a driver.

## Completion

* `--zsh-completion`: Output zsh completion code.

* `--aliases`: Output aliases code.

## Help

* `--version`: Print version.

* `--help`: Print help.
