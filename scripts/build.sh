#!/bin/bash

#
# Build all supported platform binaries.
#

name="zssh"

abs_dirname() {
  local cwd="$(pwd)"
  local path="$1"

  while [ -n "$path" ]; do
    cd "${path%/*}"
    local name="${path##*/}"
    path="$(readlink "$name" || true)"
  done

  pwd -P
  cd "$cwd"
}

root_dir="$(abs_dirname $(abs_dirname "$0"))"
cd $root_dir

echo "==> Removing old files..."
rm -rf build/*

echo "==> Building..."

gox \
    -os="linux darwin windows" \
    -output "build/${name}_{{.OS}}_{{.Arch}}" \
    ./cmd/${name}

echo "==> ZIP Packaging..."

cd build/
mv ${name}_darwin_386 ${name} && zip ${name}_darwin_386.zip ${name} && rm ${name}
mv ${name}_darwin_amd64 ${name} && zip ${name}_darwin_amd64.zip ${name} && rm ${name}

mv ${name}_linux_386 ${name} && zip ${name}_linux_386.zip ${name} && rm ${name}
mv ${name}_linux_amd64 ${name} && zip ${name}_linux_amd64.zip ${name} && rm ${name}
mv ${name}_linux_arm ${name} && zip ${name}_linux_arm.zip ${name} && rm ${name}

mv ${name}_windows_386.exe ${name}.exe && zip ${name}_windows_386.zip ${name}.exe && rm ${name}.exe
mv ${name}_windows_amd64.exe ${name}.exe && zip ${name}_windows_amd64.zip ${name}.exe && rm ${name}.exe
