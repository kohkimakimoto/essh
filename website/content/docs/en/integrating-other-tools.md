+++
title = "Integrating Other Tools | Documentation"
type = "docs"
category = "docs"
lang = "en"
basename = "integrating-other-tools.html"
+++

# Integrating Other Tools

Essh can be used with `scp`, `rsync` and `git`.

## git

To use Essh inside of the git command. Write the following line in your `~/.zshrc`.

~~~
export GIT_SSH=essh
~~~

## scp

Essh supports to use with scp.

~~~
$ essh --exec 'scp -F $ESSH_SSH_CONFIG <scp command args...>'
~~~

For more easy to use, you can run `eval "$(essh --aliases)"` in your `~/.zshrc`, the above code can be written as the following.

~~~
$ escp <scp command args...>
~~~

## rsync

Essh supports to use with rsync.

~~~
$ essh --exec 'rsync -e "ssh -F $ESSH_SSH_CONFIG" <rsync command args...>'
~~~

For more easy to use, you can run `eval "$(essh --aliases)"` in your `~/.zshrc`, the above code can be written as the following.

~~~
$ ersync <rsync command args...>
~~~
