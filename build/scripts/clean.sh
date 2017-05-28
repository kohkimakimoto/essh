#!/usr/bin/env bash
set -eu

# Get the directory path.
SOURCE="${BASH_SOURCE[0]}"
while [ -h "$SOURCE" ] ; do SOURCE="$(readlink "$SOURCE")"; done
scripts_dir="$( cd -P "$( dirname "$SOURCE" )/" && pwd )"
outputs_dir="$(cd $scripts_dir/../outputs && pwd)"
repo_dir="$(cd $scripts_dir/../.. && pwd)"

echo "Cleaning old files."
rm -rf $outputs_dir/*
