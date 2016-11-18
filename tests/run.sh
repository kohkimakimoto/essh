#!/usr/bin/env bash

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

testfail() {
  { if [ "$#" -eq 0 ]; then cat -
    else echo "${txtred}fail:${txtreset}" && echo "${txtred}$*${txtreset}"
    fi
  } >&2
  exit 1
}

testok() {
  { if [ "$#" -eq 0 ]; then cat -
    else echo "${txtgreen}ok:${txtreset}" && echo "${txtgreen}$*${txtreset}"
    fi
  } >&2
}

tests_dir=$(cd $(dirname $0); pwd)
cd $tests_dir

export PATH="$PATH:$tests_dir/../_build/dev"

# ----------------------------------------------------------------
# configuration
# ----------------------------------------------------------------
# convert ssh-config to essh config
vagrant ssh-config | perl -pe 's/Host (.+)$/s = private_host "$1" /; s/^(  )(\w)/s.$2/;  s/^(s\.\w+)( )/$1 = /; s/= ([\w\.\/]+)$/= "$1"/; ' > esshconfig.lua

# add tasks
cat << 'EOF' >> esshconfig.lua

local bash = import "github.com/kohkimakimoto/essh/modules/bash"

task "list-hosts" {
    backend = "local",
    targets = {"webserver-01", "webserver-02"},
    script = {
      [=[
      echo "$ESSH_HOST_HOSTNAME"
      ]=],
    },
}

task "echo-on-remote" {
    backend = "remote",
    prefix = true,
    targets = {"webserver-01", "webserver-02"},
    script = {
      [=[
      echo "foobar"
      ]=],
    },
}

task "user-bash-module" {
    backend = "remote",
    prefix = true,
    targets = {"webserver-01", "webserver-02"},
    script = {
      bash.indent,
      bash.prefix,
      bash.upper,
      bash.lock,
      [=[
      echo "foobar" | indent
      echo $(upper "foobar")
      ]=],
    },
}
EOF




# ----------------------------------------------------------------
# start testing...
# ----------------------------------------------------------------
echo "tested binary is: $(which essh)"
echo "tasks:"
essh --tasks

# ----
echo "==> test-1:"
ret=`essh list-hosts`
exp=`cat << 'EOF'
webserver-01
webserver-02
EOF`
if [ "$ret" = "$exp" ]; then
    testok "$ret"
else
    testfail "$ret"
fi

# ----
echo "==> test-2:"
ret=`essh echo-on-remote`
exp=`cat << 'EOF'
[remote:webserver-01] foobar
[remote:webserver-02] foobar
EOF`
if [ "$ret" = "$exp" ]; then
    testok "$ret"
else
    testfail "$ret"
fi

# ----
echo "==> test-3:"
ret=`essh user-bash-module`
exp=`cat << 'EOF'
[remote:webserver-01]     foobar
[remote:webserver-01] FOOBAR
[remote:webserver-02]     foobar
[remote:webserver-02] FOOBAR
EOF`
if [ "$ret" = "$exp" ]; then
    testok "$ret"
else
    testfail "$ret"
fi
