#!/usr/bin/env bash

#
# bash boilerplate.
# You should edit it for your purpose.
#

##################################################
# General functions, initializing and etc..
# You probably should not edit this section.
##################################################
READLINK=$(type -p greadlink readlink | head -1)
if [ -z "$READLINK" ]; then
  echo "cannot find readlink - are you missing GNU coreutils?" >&2
  exit 1
fi

resolve_link() {
  $READLINK "$1"
}

# get absolute path.
abs_dirname() {
  local cwd="$(pwd)"
  local path="$1"

  while [ -n "$path" ]; do
    # cd "${path%/*}" does not work in "$ bash script.sh"
    # cd "${path%/*}"
    cd "$(dirname $path)"
    local name="${path##*/}"
    path="$(resolve_link "$name" || true)"
  done

  pwd -P
  cd "$cwd"
}

# inspired by https://github.com/heroku/heroku-buildpack-php
# Usage:
#   echo "message" | prefix "hoge"
prefix() {
  local p="${1:-prefix}"
  local c="s/^/$p/"
  case $(uname) in
    Darwin) sed -l "$c";; # mac/bsd sed: -l buffers on line boundaries
    *)      sed -u "$c";; # unix/gnu sed: -u unbuffered (arbitrary) chunks of data
  esac
}

# inspired by https://github.com/heroku/heroku-buildpack-php
# Usage:
#   echo "message" | indent 4
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

# inspired by http://stackoverflow.com/questions/3231804/in-bash-how-to-add-are-you-sure-y-n-to-any-command-or-alias
# Usage:
#   confirm "message"
confirm() {
    local response
    # call with a prompt string or use a default
    read -r -p "${1:-Are you sure? [y/N]:} " response
    case $response in
        [yY][eE][sS]|[yY])
            return 0
            ;;
        *)
            return 1
            ;;
    esac
}

# Usage:
#   var=$(ask "message")
ask() {
    local response
    # call with a prompt string or use a default
    read -r -p "${1:->} " response
    echo $response
}

# write a horizontal line
# http://wiki.bash-hackers.org/snipplets/print_horizontal_line
# Usage:
#   hr
#   hr "="
#   hr "=" 10
hr() {
    printf '%*s\n' "${2:-$(tput cols)}" '' | tr ' ' "${1:--}"
}

# inspired by http://dharry.hatenablog.com/entry/20110122/1295681180
# Usage:
#   sleep 3 & progress
progress() {
  local _bar="${1:-.}"
  while :
  do
    # about "&&:" http://qiita.com/ngyuki/items/aefd47700a9522fada75
    jobs %1 > /dev/null 2>&1 &&:
    [ $? -eq 0 ] || break
    echo -n ${_bar}
    sleep 0.2
  done;
}

# Usage:
#   sleep 3 & loading
loading() {
  local _ptn=0
  while :
  do
    jobs %1 > /dev/null 2>&1 &&:
    [ $? -eq 0 ] || break
    if [ ${_ptn} -eq 0 ]; then
        _ptn=1
        echo -ne '-\r'
    elif [ ${_ptn} -eq 1 ]; then
        _ptn=2
        echo -ne '\\\r'
    elif [ ${_ptn} -eq 2 ]; then
        _ptn=3
        echo -ne '|\r'
    else
        _ptn=0
        echo -ne '/\r'
    fi
    sleep 0.1
  done;
}

# Usage:
#   upper "abcdefg"
#   $(upper "abcdefg")
upper() {
    echo -n "$1" | tr '[a-z]' '[A-Z]'
}

# https://github.com/rbenv/rbenv
# Usage: abort "error message"
abort() {
  { if [ "$#" -eq 0 ]; then cat -
    else echo "${txtred}${progname}: $*${txtreset}"
    fi
  } >&2
  exit 1
}

# bold and color text utility
# https://linuxtidbits.wordpress.com/2008/08/11/output-color-on-bash-scripts/
# http://stackoverflow.com/questions/2924697/how-does-one-output-bold-text-in-bash
# usage:
#   ${txtred}foobar${txtreset}
txtunderline=$(tput sgr 0 1)     # Underline
txtbold=$(tput bold)             # Bold

txtred=$(tput setaf 1)           # red
txtgreen=$(tput setaf 2)         # green
txtyellow=$(tput setaf 3)        # yellow
txtblue=$(tput setaf 4)          # blue
txtreset=$(tput sgr0)            # Reset

##################################################
# set useful variables
##################################################
set -eu
script_dir="$(abs_dirname "$0")"
progname=$(basename $0)
progversion="0.1.0"

##################################################
# Actions.
##################################################
name="zssh"

usage() {
    echo "Usage: $progname [OPTIONS] ${txtunderline}COMMAND${txtreset}"
    echo
    echo "Options:"
    echo "  -h, --help       show help."
    echo
    echo "Commands:"
    echo "  help        show help."
    echo "  dev         build a dev binary."
    echo "  dist        build all platform binaries for distribution."
    echo "  packaging   build rpm package."
    echo
}

do_dev() {
    cd $script_dir
    echo "--> Building ${txtyellow}${txtbold}dev${txtreset} binary."

    cd ..
    go build -o="build/dev/${name}" cmd/${name}/${name}.go
    if [ $? -eq 0 ]; then
        echo "Updated $(pwd)/build/dev/${name}" | indent
    fi
    echo "--> Done."
}

do_dist() {
    cd $script_dir
    echo "--> Building ${txtyellow}${txtbold}dist${txtreset} binaries."

    echo "--> Removing old files..."
    rm -rf dist/*

    echo "--> Building binaries..."
    cd ..
    gom exec gox \
        -os="linux darwin" \
        -output "build/dist/${name}_{{.OS}}_{{.Arch}}" \
        ./cmd/${name} \
        | indent & loading

    echo "--> Packaging to zip archives..."
    cd build/dist

    echo "Packaging ${txtyellow}${txtbold}darwin${txtreset}" binaries | indent
    mv ${name}_darwin_386 ${name} && zip ${name}_darwin_386.zip ${name} | indent && rm ${name}
    mv ${name}_darwin_amd64 ${name} && zip ${name}_darwin_amd64.zip ${name}  | indent && rm ${name}

    echo "Packaging ${txtyellow}${txtbold}linux${txtreset}" binaries | indent
    mv ${name}_linux_386 ${name} && zip ${name}_linux_386.zip ${name}  | indent && rm ${name}
    mv ${name}_linux_amd64 ${name} && zip ${name}_linux_amd64.zip ${name}  | indent && rm ${name}
    mv ${name}_linux_arm ${name} && zip ${name}_linux_arm.zip ${name}  | indent && rm ${name}

    echo "Updated $(pwd)" | indent
    echo "--> Done."
}

do_packaging() {
    cd $script_dir
    echo "--> Start ${txtyellow}${txtbold}packaigng${txtreset}."
    cd rpm
    vagrant up | indent
}

#do_destroy() {
#
##    cd rpm
##    vagrant destroy -f | indent
#}

##################################################
# Main logic.
##################################################
cd $script_dir


# parse arguments and options.
# inspired by http://qiita.com/b4b4r07/items/dcd6be0bb9c9185475bb
declare -a params=()

for OPT in "$@"
do
    case "$OPT" in
        '-h'|'--help' )
            usage
            exit 0
            ;;
        '-v'|'--version' )
            printversion
            exit 0
            ;;
        '--'|'-' )
            shift 1
            params+=( "$@" )
            break
            ;;
        -*)
            echo "$progname: illegal option -- '$(echo $1 | sed 's/^-*//')'" 1>&2
            exit 1
            ;;
        *)
            if [[ ! -z "$1" ]] && [[ ! "$1" =~ ^-+ ]]; then
                params+=( "$1" )
                shift 1
            fi
            ;;
    esac
done

command="" && [ ${#params[@]} -ne 0 ] && command=${params[0]}
case $command in
    'help' )
        usage
        exit 0
        ;;
    'dev' )
        do_dev
        exit 0
        ;;
    'dist' )
        do_dist
        exit 0
        ;;
    'packaging' )
        do_packaging
        exit 0
        ;;
    'destroy' )
        do_destroy
        exit 0
        ;;
    '' )
        usage
        exit 0
        ;;
    *)
        echo "$progname: illegal command '$command'" 1>&2
        exit 1
        ;;
esac
