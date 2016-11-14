+++
title = "Running Commands"
type = "docs"
category = "intro"
lang = "en"
+++

# Running Commands

Essh allow you to run commands on the selected remote hosts by using `--exec`, `--backend` and `--target` options.

~~~sh
$ essh --exec --backend=remote --target=web uptime
 22:48:31 up  7:58,  0 users,  load average: 0.00, 0.01, 0.03
 22:48:31 up  7:58,  0 users,  load average: 0.00, 0.02, 0.04
~~~

Use `--prefix` option, Essh outputs result of command with hostname prefix.

~~~sh
$ essh --exec --backend=remote --target=web --prefix uptime
[web01.localhost]  22:48:31 up  7:58,  0 users,  load average: 0.00, 0.01, 0.03
[web02.localhost]  22:48:31 up  7:58,  0 users,  load average: 0.00, 0.02, 0.04
~~~

Let's read next section: [Running Tasks](running-tasks.html)
