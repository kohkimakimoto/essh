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

# Checks if it has loaded .envrc by direnv.
if [ -z ${DIRENV_DIR+x} ]; then
    if [ -e "./.envrc" ]; then
        source ./.envrc
    fi
fi

source ./_build/config

echo "--> Building RPM packages..."
# building PRMs by using docker.
cd _build/packaging/rpm
for image in 'kohkimakimoto/rpmbuild:el5' 'kohkimakimoto/rpmbuild:el6' 'kohkimakimoto/rpmbuild:el7'; do
    docker run \
        --env DOCKER_IMAGE=${image}  \
        --env PRODUCT_NAME=${PRODUCT_NAME}  \
        --env PRODUCT_VERSION=${PRODUCT_VERSION}  \
        --env COMMIT_HASH=${COMMIT_HASH}  \
        -v $DIR:/tmp/repo \
        -w /tmp/repo \
        --rm \
        ${image} \
        bash ./_build/packaging/rpm/run.sh
done

cd "$DIR"

echo "--> Results:"
ls -hl "_build/dist/" | indent

echo "--> Done."
