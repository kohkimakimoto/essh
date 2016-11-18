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
      set -e
      echo "foobar" | indent
      echo $(upper "foobar")
      ]=],
    },
}

task "user-bash-module-privilaged" {
    backend = "remote",
    prefix = true,
    privileged = true,
    targets = {"webserver-01", "webserver-02"},
    script = {
      bash.indent,
      bash.prefix,
      bash.lock,
      [=[
      set -e
      echo "foobar" | indent
      ]=],
    },
}

task "user-bash-module-local" {
    backend = "local",
    prefix = true,
    targets = {"webserver-01", "webserver-02"},
    script = {
      bash.indent,
      bash.prefix,
      bash.lock,
      [=[
      set -e
      echo "foobar" | indent
      ]=],
    },
}

task "user-bash-module-local-privilaged" {
    backend = "local",
    prefix = true,
    privileged = true,
    targets = {"webserver-01", "webserver-02"},
    script = {
      bash.indent,
      bash.prefix,
      bash.lock,
      [=[
      set -e
      echo "foobar" | indent
      ]=],
    },
}

task "load-stdin-remote" {
    backend = "remote",
    prefix = true,
    targets = {"webserver-01", "webserver-02"},
    script = {
      [=[
      set -e
      cat -
      ]=],
    },
}

task "load-stdin-remote-privilaged" {
    backend = "remote",
    prefix = true,
    privileged = true,
    targets = {"webserver-01", "webserver-02"},
    script = {
      [=[
      set -e
      whoami
      cat -
      ]=],
    },
}


local docker = import "github.com/kohkimakimoto/essh/modules/docker"

driver "docker-centos6" {
    engine = docker.driver,
    image = "centos:centos6",
    privileged = true,
    remove_terminated_containers = true,
}

task "docker-module" {
    backend = "remote",
    prefix = true,
    targets = {"webserver-01", "webserver-02"},
    driver = "docker-centos6",
    script = {
      [=[
      echo "aaaa"
      ]=],
    },
}

task "task-in-task" {
    backend = "local",
    prefix = true,
    targets = {"webserver-01", "webserver-02"},
    script = {
      [=[
        echo "start"
        essh task-in-task2
        echo "end"
      ]=],
    },
}

task "task-in-task2" {
    backend = "remote",
    prefix = true,
    targets = {"webserver-01", "webserver-02"},
    script = {
      [=[
        echo "foo"
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

# ----
echo "==> test-4:"
ret=`essh user-bash-module-privilaged`
exp=`cat << 'EOF'
[remote:webserver-01]     foobar
[remote:webserver-02]     foobar
EOF`
if [ "$ret" = "$exp" ]; then
    testok "$ret"
else
    testfail "$ret"
fi

# ----
echo "==> test-5:"
ret=`essh user-bash-module-local`
exp=`cat << 'EOF'
[local:webserver-01]     foobar
[local:webserver-02]     foobar
EOF`
if [ "$ret" = "$exp" ]; then
    testok "$ret"
else
    testfail "$ret"
fi

# # ----
# echo "==> test-5.1:"
# ret=`essh user-bash-module-local-privilaged`
# exp=`cat << 'EOF'
# [local:webserver-01]     foobar
# [local:webserver-02]     foobar
# EOF`
# if [ "$ret" = "$exp" ]; then
#     testok "$ret"
# else
#     testfail "$ret"
# fi

# ----
echo "==> test-6:"
ret=`echo hogehoge | essh load-stdin-remote`
exp=`cat << 'EOF'
[remote:webserver-01] hogehoge
[remote:webserver-02] hogehoge
EOF`
if [ "$ret" = "$exp" ]; then
    testok "$ret"
else
    testfail "$ret"
fi

# ----
echo "==> test-7:"
ret=`echo hogehoge | essh load-stdin-remote-privilaged`
exp=`cat << 'EOF'
[remote:webserver-01] root
[remote:webserver-01] hogehoge
[remote:webserver-02] root
[remote:webserver-02] hogehoge
EOF`
if [ "$ret" = "$exp" ]; then
    testok "$ret"
else
    testfail "$ret"
fi

# ----
echo "==> test-8: checks just running"
essh docker-module
if [ ! "$?" = "0" ]; then
    testfail "exited with non zero"
else
    testok "exited with zero"
fi

# ----
echo "==> test-9:"
ret=`essh task-in-task`
exp=`cat << 'EOF'
[local:webserver-01] start
[local:webserver-01] [remote:webserver-01] foo
[local:webserver-01] [remote:webserver-02] foo
[local:webserver-01] end
[local:webserver-02] start
[local:webserver-02] [remote:webserver-01] foo
[local:webserver-02] [remote:webserver-02] foo
[local:webserver-02] end
EOF`
if [ "$ret" = "$exp" ]; then
    testok "$ret"
else
    testfail "$ret"
fi
