# Integrating Other Tools

Essh can be used with `scp`, `rsync` and `git`.

* `git`: To use Essh inside of the git command. Write the following line in your `~/.zshrc`.

  ```
  export GIT_SSH=essh
  ```

* `scp`: Essh support to use with scp.

  ```
  $ essh --scp <scp command args...>
  ```

  For more easy to use, you can run `eval "$(essh --aliases)"` in your `~/.zshrc`, the above code can be written as the following.

  ```
  $ escp <scp command args...>
  ```

* `rsync`: Essh support to use with rsync.

  ```
  $ essh --rsync <rsync command args...>
  ```

  For more easy to use, you can run `eval "$(essh --aliases)"` in your `~/.zshrc`, the above code can be written as the following.

  ```
  $ ersync <rsync command args...>
  ```
