
driver "essh-lua" {
    engine = [=[
{{template "environment" .}}

# Detect if essh is installed
if ! command -v essh &> /dev/null
then
    # Try to extend the PATH
    export PATH=$PATH:$HOME/.bin
    if ! command -v essh &> /dev/null
    then
        echo "essh could not be found. Please install it."
        exit 1
    fi
fi

set -e

essh --eval <<EOF
{{range $i, $script := .Scripts -}}
{{$script.code}} 
{{end -}}
EOF

        set +e
    ]=],
}