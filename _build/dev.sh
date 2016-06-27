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

# Get the parent directory of where this script is.
SOURCE="${BASH_SOURCE[0]}"
while [ -h "$SOURCE" ] ; do SOURCE="$(readlink "$SOURCE")"; done
DIR="$( cd -P "$( dirname "$SOURCE" )/.." && pwd )"

cd "$DIR"
source ./_build/config

PRODUCT_VERSION="dev"

echo "--> Building dev binary..."
echo "    PRODUCT_NAME: $PRODUCT_NAME"
echo "    PRODUCT_VERSION: $PRODUCT_VERSION"
echo "    COMMIT_HASH: $COMMIT_HASH"

go build \
    -ldflags=" -w \
        -X github.com/kohkimakimoto/$PRODUCT_NAME/$PRODUCT_NAME.CommitHash=$COMMIT_HASH \
        -X github.com/kohkimakimoto/$PRODUCT_NAME/$PRODUCT_NAME.Version=$PRODUCT_VERSION" \
    -o="_build/dev/$PRODUCT_NAME" \
    ./cmd/$PRODUCT_NAME/$PRODUCT_NAME.go

echo "--> Results:"
ls -hl "_build/dev/" | indent

echo "--> Done."
