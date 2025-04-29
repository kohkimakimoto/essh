#!/usr/bin/env bash
set -eu

# Get the directory path.
SOURCE="${BASH_SOURCE[0]}"
while [ -h "$SOURCE" ] ; do SOURCE="$(readlink "$SOURCE")"; done
scripts_dir="$( cd -P "$( dirname "$SOURCE" )/" && pwd )"
build_dir="$(cd $scripts_dir/.. && pwd)"
outputs_dir="$(cd $build_dir/outputs && pwd)"
repo_dir="$(cd $build_dir/.. && pwd)"

# Move the parent (repository) directory
cd "$repo_dir"

# Load config
source $scripts_dir/config

echo "Removing old files."
rm -rf $outputs_dir/dist/*

COMMIT_HASH=`git log --pretty=format:%H -n 1`

echo "Building dev binary..."
echo "PRODUCT_NAME: $PRODUCT_NAME"
echo "PRODUCT_VERSION: $PRODUCT_VERSION"
echo "COMMIT_HASH: $COMMIT_HASH"

echo "Building binaries for windows and macosx..."
gox \
    -os="darwin windows" \
    -arch="amd64 arm64" \
    -ldflags=" -w \
        -X github.com/kohkimakimoto/$PRODUCT_NAME/$PRODUCT_NAME.CommitHash=$COMMIT_HASH \
        -X github.com/kohkimakimoto/$PRODUCT_NAME/$PRODUCT_NAME.Version=$PRODUCT_VERSION" \
    -output "$outputs_dir/dist/${PRODUCT_NAME}_{{.OS}}_{{.Arch}}" \
    ./cmd/${PRODUCT_NAME}

export CGO_ENABLED=0
export CGO_CFLAGS="-static"
export CGO_LDFLAGS="--static"

echo "Building linux static binary..."
gox -os="linux" -arch="amd64 arm64" -ldflags "-s -w \
        -X github.com/kohkimakimoto/$PRODUCT_NAME/$PRODUCT_NAME.CommitHash=$COMMIT_HASH \
        -X github.com/kohkimakimoto/$PRODUCT_NAME/$PRODUCT_NAME.Version=$PRODUCT_VERSION" \
    -gcflags " -l -l -l" \
    -output "$outputs_dir/dist/${PRODUCT_NAME}_{{.OS}}_{{.Arch}}" ./cmd/${PRODUCT_NAME}

echo "Packaging to zip archives..."

cd "$outputs_dir/dist"
for f in *; do
    if [ -f "$f" ]; then
        zip -r "$f.zip" "$f"
        rm -rf "$f"
    fi
done

cd "$repo_dir"

echo "Results:"
ls -hl "$outputs_dir/dist"

