local bash = {}

bash.indent = [=[
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
]=]

bash.prefix = [=[
prefix() {
  local p="${1:-prefix}"
  local c="s/^/$p/"
  case $(uname) in
    Darwin) sed -l "$c";; # mac/bsd sed: -l buffers on line boundaries
    *)      sed -u "$c";; # unix/gnu sed: -u unbuffered (arbitrary) chunks of data
  esac
}
]=]

bash.upper = [=[
upper() {
    echo -n "$1" | tr '[a-z]' '[A-Z]'
}
]=]

bash.xterm = "TERM=xterm"

bash.errexit_on = "set -e"

bash.version = "bash --version"

bash.lock = [=[
essh_bash_lockdir=${TMPDIR:-/tmp}
essh_bash_lockdir=${essh_bash_lockdir%/}/essh_lock.${ESSH_TASK_NAME:-unknown}
essh_bash_trylocking=0
essh_bash_trylocking_count=0

while [ ${essh_bash_trylocking} -eq 0 ]
do
    if ( mkdir ${essh_bash_lockdir} ) 2> /dev/null; then
        echo $$ > ${essh_bash_lockdir}/pid
        # break loop
        essh_bash_trylocking=1
        trap 'rm -rf "$essh_bash_lockdir"; exit $?' INT TERM EXIT
    else
        # could not get a lock, try to check pid.
        if ps -p $(cat ${essh_bash_lockdir}/pid) > /dev/null 2>&1; then
            # pid exists.
            echo "Lock exists: $essh_bash_lockdir owned by $(cat ${essh_bash_lockdir}/pid)" >&2
            exit 1
        else
            # pid does not exist. remove lockdir
            echo "Lock exists: $essh_bash_lockdir owned by $(cat ${essh_bash_lockdir}/pid). But the process terminated." >&2
            echo "Trying to clean the lock and start it again..." >&2
            rm -rf "$essh_bash_lockdir"

            essh_bash_trylocking_count=$((essh_bash_trylocking_count+1))

            if [ "$essh_bash_trylocking_count" -gt 5 ]; then
                echo "Could not get the lock." >&2
                exit 1
            fi
        fi
    fi
done
]=]

bash.driver = [=[
{{template "environment" .}}
{{if .Driver.Props.show_version -}}
echo "==> bash version:"
bash --version
{{end}}

__essh_var_status=0
{{range $i, $script := .Scripts}}
if [ $__essh_var_status -eq 0 ]; then
echo "==> script: {{$i}}{{if $script.description}} ({{$script.description}}){{end}}"
{{$script.code}}
__essh_var_status=$?
fi
{{end}}
exit $__essh_var_status
]=]

return bash
