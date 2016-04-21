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

bash.driver = function()
    return [=[
__essh_var_status=0
{{range $i, $script := .Scripts}}
if [ $__essh_var_status -eq 0 ]; then
echo '--> Script step {{$i}}: {{if $script.name}}{{$script.name}}{{end}}'
{{$script.code}}
__essh_var_status=$?
fi
{{end}}
exit $__essh_var_status
]=]
end


return bash
