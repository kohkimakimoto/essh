#!/usr/bin/env bash
set -eu

indent() {
    local n="${1:-4}"
    local p=""
    for i in `seq 1 $n`; do
        p="$p "
    done;

    local c="s/^/$p/"
    case $(uname) in
      Darwin) sed -l "$c";; # mac/bsd sed: -l buffers on line boundaries
      *)      sed -u "$c";; # unix/gnu sed: -u unbuffered (arbitrary) chunks of data
    esac
}

if [ "${TERM:-dumb}" != "dumb" ]; then
    txtunderline=$(tput sgr 0 1)     # Underline
    txtbold=$(tput bold)             # Bold
    txtred=$(tput setaf 1)           # red
    txtgreen=$(tput setaf 2)         # green
    txtyellow=$(tput setaf 3)        # yellow
    txtblue=$(tput setaf 4)          # blue
    txtreset=$(tput sgr0)            # Reset
else
    txtunderline=""
    txtbold=""
    txtred=""
    txtgreen=""
    txtyellow=""
    txtblue=$""
    txtreset=""
fi

GOTEST_FLAGS=${GOTEST_FLAGS:--cover -timeout=360s}
DOCKER_IMAGE=${DOCKER_IMAGE:-"kohkimakimoto/ssh"}

test_dir=$(cd $(dirname $0); pwd)
cd "$test_dir/.."

echo "--> Running tests (flags: $GOTEST_FLAGS)..."

echo "--> Starting a docker container as a test SSH server..."
docker run -d -P --name essh_test_ssh_server $DOCKER_IMAGE 2>&1 | indent
trap "echo '--> Removing tarminated containers...' && docker rm `docker ps -a -q` 2>&1 | indent" EXIT HUP INT QUIT TERM

GOBIN="`which go`"
$GOBIN test $GOTEST_FLAGS $($GOBIN list ./... | grep -v vendor) 2>&1 | perl -pe "s/^ok/${txtgreen}ok${txtreset}/; s/^FAIL/${txtred}FAIL${txtreset}/;" | indent

echo "--> Done."