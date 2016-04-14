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
      Darwin) sed -l "$c";;
      *)      sed -u "$c";;
    esac
}
]=]

return bash
