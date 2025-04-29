# Changelog

## [Unreleased](https://github.com/sevir/essh/compare/v3.8.2...a6daa43f95cffd4014422cc6314dedfea1e232b8) (2025-04-24)

### [v3.8.2](https://github.com/sevir/essh/compare/v3.8.1...v3.8.2) (2025-04-12)

#### Fixes

* **zsh-completion:** Fix zsh completion script
([8544b04](https://github.com/sevir/essh/commit/8544b04b1d747d751d39022f2a8e08b508bff1ba))

### [v3.8.1](https://github.com/sevir/essh/compare/v3.8.0...v3.8.1) (2025-03-22)

## [v3.8.0](https://github.com/sevir/essh/compare/v3.7.2...v3.8.0) (2025-03-20)

### Features

* **lua:** Add new lua functions for aws, mdns and watch file changes
([cd1dfe8](https://github.com/sevir/essh/commit/cd1dfe81c95a880f74928beca537ea4a190e4641))

### [v3.7.2](https://github.com/sevir/essh/compare/v3.7.1...v3.7.2) (2025-03-10)

#### Fixes

* Problem with the correct release version in build process
([786c2cc](https://github.com/sevir/essh/commit/786c2cc46ea313cc9649ed736bbc852cad53ee79))

### [v3.7.1](https://github.com/sevir/essh/compare/v3.7.0...v3.7.1) (2025-03-10)

## [v3.7.0](https://github.com/sevir/essh/compare/v3.6.2...v3.7.0) (2025-03-10)

### Features

* **menu:** Add tui menu for hosts with new flag --menu
([fac2ab6](https://github.com/sevir/essh/commit/fac2ab68e58199e1698ffa80befff4dd6fc53dc8))
* **build:** Add upx compression binary for linux amd64 version (the standard
server system)
([58809a9](https://github.com/sevir/essh/commit/58809a9fdd22482211f5bc26c64ea7a1d0d537f7))

### Fixes

* **menu:** Improve the menu of tasks and hosts showing the type and detect
the screen size
([9f5d386](https://github.com/sevir/essh/commit/9f5d3866d68b36e8810fe5b9306c4b5e4ee073b0))

### [v3.6.2](https://github.com/sevir/essh/compare/v3.6.1...v3.6.2) (2025-03-08)

#### Features

* **luavm:** Add new flag --eval-file passing a lua file for execute it
([478a748](https://github.com/sevir/essh/commit/478a74857417798e47f1aab9c1f328f741efad55))

### [v3.6.1](https://github.com/sevir/essh/compare/v3.6.0...v3.6.1) (2025-03-08)

## [v3.6.0](https://github.com/sevir/essh/compare/v3.5.4...v3.6.0) (2025-03-08)

### [v3.5.4](https://github.com/sevir/essh/compare/v3.5.0...v3.5.4) (2025-03-08)

#### Features

* **lua:** Add a  set of lua libraries embedded like log, time, runtime,
base64, crypto, storage, strings
([7f4f71b](https://github.com/sevir/essh/commit/7f4f71b1ec1b9307d25b518c95228cb683c911ef))

#### Fixes

* **github-pages:** Update github pages
([a4855f9](https://github.com/sevir/essh/commit/a4855f9fe6aab0ac6ffa49f8bdc4e68fe1f84979))

## v3.5.0 (2025-02-24)

### Features

* **build:** Build linux version without dependencies, static compiling
([706b4cb](https://github.com/sevir/essh/commit/706b4cbd5d0efcd5c1e5f0b0a6739654d11e2e01))
* **eval-lua:** New flag parameter allows to eval lua code directly from
command line, env variable or standard input
([fc14199](https://github.com/sevir/essh/commit/fc1419974c58d732e610d45e488477522dfd64cd))
* **update:** Add go.mod and go.sum for compile with Go 1.23.4 and other
disabled `go get` versions
([a03aa20](https://github.com/sevir/essh/commit/a03aa20a4b7de5d97b60883842fdeff7c54c8cd6))
