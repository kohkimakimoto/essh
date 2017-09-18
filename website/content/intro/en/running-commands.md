+++
title = "Running Commands | Introduction"
type = "docs"
category = "intro"
lang = "en"
basename = "running-commands.html"
+++

# Running Commands

Essh allow you to run commands on the selected remote hosts by using `--exec`, `--backend` and `--target` options.

~~~sh
$ essh --exec --backend=remote --target=web uptime
 22:48:31 up  7:58,  0 users,  load average: 0.00, 0.01, 0.03
 22:48:31 up  7:58,  0 users,  load average: 0.00, 0.02, 0.04
~~~

`--target` option can be used multiple times.

~~~sh
$ essh --exec --backend=remote --target=web --target=db uptime
 16:47:02 up 270 days, 13:29,  0 users,  load average: 0.11, 0.18, 0.11
 16:47:02 up 270 days, 13:26,  0 users,  load average: 0.00, 0.01, 0.00
 16:47:02 up 10 days,  1:02,  0 users,  load average: 0.01, 0.03, 0.00
 16:47:03 up 2 days, 22:24,  1 user,  load average: 0.00, 0.01, 0.05
~~~

Use `--prefix` option, Essh outputs result of command with hostname prefix.

~~~sh
$ essh --exec --backend=remote --target=web --prefix uptime
[remote:web01.localhost]  22:48:31 up  7:58,  0 users,  load average: 0.00, 0.01, 0.03
[remote:web02.localhost]  22:48:31 up  7:58,  0 users,  load average: 0.00, 0.02, 0.04
~~~

Use `--parallel` option, Essh runs commands in parallel.

~~~sh
$ essh --exec --backend=remote --target=web --prefix --parallel uptime
[remote:web01.localhost]  22:48:31 up  7:58,  0 users,  load average: 0.00, 0.01, 0.03
[remote:web02.localhost]  22:48:31 up  7:58,  0 users,  load average: 0.00, 0.02, 0.04
~~~

Use `--privileged` option, Essh runs commands by privileged (root) user. You have to configure your machine to be able to be used `sudo` without password.

~~~sh
$ essh --exec --backend=remote --target=web --prefix --privileged whoami
[remote:web01.localhost] root
[remote:web01.localhost] root
~~~

Set `--backend=local` option, Essh runs commands locally.

~~~sh
$ essh --exec --backend=local --target=web --parallel --prefix 'echo $ESSH_HOSTNAME'
[local:web01.localhost] web01.localhost
[local:web02.localhost] web02.localhost
~~~

In the above example, I use `ESSH_HOSTNAME` environment variable.
Essh runs commands by using a temporary [task](/essh/docs/en/tasks.html) internally. So you can use some predefined variables that defined by a task in the commands. For detail, see [Tasks](/essh/docs/en/tasks.html)


Let's read next section: [Running Tasks](running-tasks.html)
