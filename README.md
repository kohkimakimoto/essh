# Essh [![Build Status](https://travis-ci.org/kohkimakimoto/essh.svg?branch=master)](https://travis-ci.org/kohkimakimoto/essh)

Extended ssh command. 

* [Website](https://kohkimakimoto.github.io/essh/)
* [Documentation](https://kohkimakimoto.github.io/essh/docs/en/index.html)
* [Gettting Started](https://kohkimakimoto.github.io/essh/intro/en/index.html)

## Overview

Essh is an extended `ssh` command. If you use `essh` command instead of `ssh`, Your SSH operation becomes more efficient and convenient. Essh is a single binary CLI tool and simply wraps ssh command. You can use it in the same way as ssh. And it has useful features over ssh.

![example01.gif](https://raw.githubusercontent.com/kohkimakimoto/essh/master/example01.gif)

## Features

* **Configuration As Code**: You can write SSH client configuration (aka:`~/.ssh/config`) in [Lua](https://www.lua.org/) code. So your ssh_config can become more dynamic.

* **Hooks**: Essh supports hooks that execute commands when it connects a remote server.

* **Servers List Management**: Essh provides utilities for managing hosts, that list and classify servers by using tags.

* **Per-Project Configuration**: Essh supports per-project configuration. This allows you to change SSH hosts config by changing current working directory.

* **Task Runner**: Task is a script that runs on remote and local servers. You can use it to automate your system administration tasks.

## Installation

Essh is provided as a single binary. You can download it and drop it in your $PATH.
After installing Essh, run the `essh` without any options in your terminal to check working.

### Homebrew

```
$ brew install kohkimakimoto/essh/essh
```

### Download the binary from releases page

[Download latest version](https://github.com/kohkimakimoto/essh/releases/latest)

## Development

Requirements

* Go 1.23 or later (my development env)

## Tasks

Tasks for xcfile.

### build

Building distributed binaries.

```sh
make dist
```

### test-lua-code

interactive: true

```sh
go run cmd/essh/essh.go --eval <<EOF
$(cat ./test/test.lua)
EOF
```
### install:deps

Installing dependences

```
make deps
```

### dev

Building dev binary.

```
make dev
```

### tag

Deploys a new tag for the repo.

Specify major/minor/patch with VERSION

Env: PRERELEASE=0, VERSION=minor, FORCE_VERSION=0
Inputs: VERSION, PRERELEASE, FORCE_VERSION


```
# https://github.com/unegma/bash-functions/blob/main/update.sh

CURRENT_VERSION=`git describe --abbrev=0 --tags 2>/dev/null`
CURRENT_VERSION_PARTS=(${CURRENT_VERSION//./ })
VNUM1=${CURRENT_VERSION_PARTS[0]}
# remove v
VNUM1=${VNUM1:1}
VNUM2=${CURRENT_VERSION_PARTS[1]}
VNUM3=${CURRENT_VERSION_PARTS[2]}

if [[ $VERSION == 'major' ]]
then
  VNUM1=$((VNUM1+1))
  VNUM2=0
  VNUM3=0
elif [[ $VERSION == 'minor' ]]
then
  VNUM2=$((VNUM2+1))
  VNUM3=0
elif [[ $VERSION == 'patch' ]]
then
  VNUM3=$((VNUM3+1))
else
  echo "Invalid version"
  exit 1
fi

NEW_TAG="v$VNUM1.$VNUM2.$VNUM3"

# if command convco is available, use it to check the version
if command -v convco &> /dev/null
then
  # if the version is a prerelease, add the prerelease tag
  if [[ $PRERELEASE == '1' ]]
  then
    NEW_TAG=v$(convco version -b --prerelease)
  else
    NEW_TAG=v$(convco version -b)
  fi
fi

# if $FORCE_VERSION is different to 0 then use it as the version
if [[ $FORCE_VERSION != '0' ]]
then
  NEW_TAG=v$FORCE_VERSION
fi

echo Adding git tag with version ${NEW_TAG}
git tag ${NEW_TAG}
git push origin ${NEW_TAG}
```

### changelog

Generate a changelog for the repo.

```
convco changelog > CHANGELOG.md
git add CHANGELOG.md
git commit -m "Update changelog"
git push
```

### release

Releasing a new version into the repo.

```
goreleaser release --clean --skip sign
```

### release-snapshot

Releasing a new snapshot version into the repo.

```
goreleaser release --snapshot --skip sign --clean
```

### package-rpm

Building packages (now support only RPM)

require: build

```
make packaging
```

### publish-docs

Publishing docs to gh-pages

```
cd website
rm -rf public
./hugo
./scripts/gh-pages-publish
```

## Author

Original Author (before 3.5.0 version):
Kohki Makimoto <kohki.makimoto@gmail.com>

Current Maintainer:
José Francisco Rives <jose@sevir.org>

## License

The MIT License (MIT)

