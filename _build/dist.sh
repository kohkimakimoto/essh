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

COMMIT_HASH=`git log --pretty=format:%H -n 1`

echo "--> Building dist binaries..."
echo "    PRODUCT_NAME: $PRODUCT_NAME"
echo "    PRODUCT_VERSION: $PRODUCT_VERSION"
echo "    COMMIT_HASH: $COMMIT_HASH"

echo "--> Removing old files..."
rm -rf _build/dist/*

echo "--> Building..."
gox \
    -os="linux darwin windows" \
    -arch="amd64" \
    -ldflags=" -w \
        -X github.com/kohkimakimoto/$PRODUCT_NAME/$PRODUCT_NAME.CommitHash=$COMMIT_HASH \
        -X github.com/kohkimakimoto/$PRODUCT_NAME/$PRODUCT_NAME.Version=$PRODUCT_VERSION" \
    -output "_build/dist/${PRODUCT_NAME}_{{.OS}}_{{.Arch}}" \
    ./cmd/${PRODUCT_NAME} \
     | indent

echo "--> Packaging to zip archives..."

cd "_build/dist"
echo "Packaging darwin binaries" | indent
mv ${PRODUCT_NAME}_darwin_amd64 ${PRODUCT_NAME} && zip ${PRODUCT_NAME}_darwin_amd64.zip ${PRODUCT_NAME}  | indent && rm ${PRODUCT_NAME}
echo "Packaging linux binaries" | indent
mv ${PRODUCT_NAME}_linux_amd64 ${PRODUCT_NAME} && zip ${PRODUCT_NAME}_linux_amd64.zip ${PRODUCT_NAME}  | indent && rm ${PRODUCT_NAME}
echo "Packaging windows binaries" | indent
mv ${PRODUCT_NAME}_windows_amd64.exe ${PRODUCT_NAME}.exe && zip ${PRODUCT_NAME}_windows_amd64.zip ${PRODUCT_NAME}.exe | indent && rm ${PRODUCT_NAME}.exe

cd "../.."

echo "--> Results:"
ls -hl "_build/dist/" | indent

echo "--> Done."
