local docker = {}

docker.driver = function(config)
    return [=[
{{template "environment" .}}
set -e

prefix() {
  local p="${1:-prefix}"
  local c="s/^/$p/"
  case $(uname) in
    Darwin) sed -l "$c";;
    *)      sed -u "$c";;
  esac
}

{{if not .Task.Registry}}
echo "error: docker driver engine can not be used with '--exec' option. You should define a task." 1>&2
exit 1
{{end}}

__essh_var_docker_working_dir=$(pwd)
__essh_var_docker_image={{if .Driver.Props.image}}{{.Driver.Props.image | ShellEscape}}{{end}}
__essh_var_docker_container_name={{if .Driver.Props.container_name}}{{.Driver.Props.container_name | ShellEscape}}{{end}}
__essh_var_docker_remove_terminated_container="{{if .Driver.Props.remove_terminated_container}}1{{end}}"
__essh_var_docker_print_docker_version="{{if .Driver.Props.print_docker_version}}1{{end}}"

__essh_var_status=0

if [ -n "$__essh_var_docker_print_docker_version" ]; then
    echo "Checking docker version..."
    docker version
    __essh_var_status=$?
    if [ $__essh_var_status -ne 0 ]; then
        echo "error: got a error when it checks the docker environment. exited with $__essh_var_status." 1>&2
        exit $__essh_var_status
    fi
fi

if [ -z "$__essh_var_docker_image" ]; then
    echo "error: docker driver engine requires 'image' config." 1>&2 
    exit 1
fi

if [ -z "$__essh_var_docker_container_name" ]; then
    __essh_var_docker_container_name_part=$(echo $__essh_var_docker_image | perl -pe "s/[:\/]/-/g;")
    __essh_var_docker_container_name="essh-${__essh_var_docker_container_name_part}-$(date +%s)"
fi

# create temporary directory
{{if .Task.IsRemoteTask}}
    __essh_var_docker_tmp_dir=$(mktemp -d /tmp/.essh_docker.XXXXXXXX)
{{else}}
    {{if .Task.Registry}}
    __essh_var_docker_tmp_dir=$(mktemp -d {{.Task.Registry.TmpDir}}/.essh_docker.XXXXXXXX)
    {{end}}
{{end}}

trap "rm -rf $__essh_var_docker_tmp_dir; exit" 0
chmod 777 $__essh_var_docker_tmp_dir

# create runfile
__essh_var_docker_runfile=$__essh_var_docker_tmp_dir/run.sh
touch $__essh_var_docker_runfile
chmod 777 $__essh_var_docker_runfile

# input content to the runfile.
cat << 'EOF-ESSH-DOCKER_SCRIPT' > $__essh_var_docker_runfile

{{template "environment" .}}

__essh_var_status=0
{{range $i, $script := .Scripts}}
if [ $__essh_var_status -eq 0 ]; then
{{$script.code}}
__essh_var_status=$?
fi
{{end}}
exit $__essh_var_status

EOF-ESSH-DOCKER_SCRIPT

echo "Running task in a docker container..."
echo "container name: '$__essh_var_docker_container_name'"
echo "image: '$__essh_var_docker_image'"

docker run \
    --privileged \
    -v $__essh_var_docker_working_dir:/essh \
    -v $__essh_var_docker_tmp_dir:/tmp/essh \
    -w /essh \
    --name $__essh_var_docker_container_name \
    $__essh_var_docker_image \
    bash /tmp/essh/run.sh
__essh_var_status=$?

if [ -n "$__essh_var_docker_remove_terminated_container" ]; then
    echo "Removing terminated container..."
    docker rm $__essh_var_docker_container_name 2>&1 | prefix 'Removed: '
fi

echo "Task exited with status $__essh_var_status."
exit $__essh_var_status
]=]
end

return docker
