local docker = {}

docker.driver = function(config)
    local sudo = ""
    if config.privileged then
        sudo = "sudo "
    end

    return [=[
{{template "environment" .}}
set -e

{{ $sudo := "]=] .. sudo .. [=[" }}

__essh_var_docker_working_dir=$(pwd)
__essh_var_docker_image={{if .Driver.Props.image}}{{.Driver.Props.image | ShellEscape}}{{end}}
__essh_var_docker_build="{{if .Driver.Props.build}}1{{end}}"
__essh_var_docker_build_url={{if .Driver.Props.build.url}}{{.Driver.Props.build.url | ShellEscape}}{{end}}
__essh_var_docker_build_dockerfile={{if .Driver.Props.build.dockerfile}}{{.Driver.Props.build.dockerfile | ShellEscape}}{{end}}

{{if ne .Task.Registry.TypeString "local"}}
echo "error: docker driver engine supports only running a task that is defined in 'local' registry." 1>&2
exit 1
{{end}}

__essh_var_status=0
echo 'Starting task by using docker driver engine.'
echo "Checking docker version."
{{$sudo}}docker version
__essh_var_status=$?
if [ $__essh_var_status -ne 0 ]; then
    echo "error: got a error when it checks the docker environment. exited with $__essh_var_status." 1>&2
    exit $__essh_var_status
fi

echo ""

if [ -z "$__essh_var_docker_image" ]; then
    echo "error: docker driver engine requires 'image' config." 1>&2
    exit 1
fi

# checks existence of the image
echo "Using image '$__essh_var_docker_image'"
if [ -z $({{$sudo}}docker images -q $__essh_var_docker_image) ]; then
    # There is not the image in the host.
    if [ -n "$__essh_var_docker_build" ]; then
        echo "There is not the image '$__essh_var_docker_image' in the running machine."
        echo "Building a docker image '$__essh_var_docker_image'..."

        if [ -n "$__essh_var_docker_build_url" ]; then
            echo "{{$sudo}}docker build -t $__essh_var_docker_image $__essh_var_docker_build_url"
            {{$sudo}}docker build -t $__essh_var_docker_image $__essh_var_docker_build_url
            __essh_var_status=$?
            if [ $__essh_var_status -ne 0 ]; then
                echo "error: got a error in docker build." 1>&2
                exit $__essh_var_status
            fi
        elif [ -n "$__essh_var_docker_build_dockerfile" ]; then
            echo "{{$sudo}}docker build -t $__essh_var_docker_image -"

            # note: double quote is needed to output multi lines
            echo "$__essh_var_docker_build_dockerfile" | docker build -t $__essh_var_docker_image -
            __essh_var_status=$?
            if [ $__essh_var_status -ne 0 ]; then
                echo "error: got a error in docker build." 1>&2
                exit $__essh_var_status
            fi
        else
            echo "error: got a error in docker build. require 'url' or 'dockerfile'" 1>&2
            exit 1
        fi
    fi
fi

# create temporary directory
{{if .Task.IsRemoteTask}}
    __essh_var_docker_tmp_dir=$({{$sudo}}mktemp -d /tmp/.essh_docker.XXXXXXXX)
{{else}}
    __essh_var_docker_tmp_dir=$({{$sudo}}mktemp -d {{.Task.Registry.TmpDir}}/.essh_docker.XXXXXXXX)
{{end}}

trap "{{$sudo}}rm -rf $__essh_var_docker_tmp_dir; exit" 0
{{$sudo}}chmod 777 $__essh_var_docker_tmp_dir

# create runfile
__essh_var_docker_runfile=$__essh_var_docker_tmp_dir/run.sh
{{$sudo}}touch $__essh_var_docker_runfile
{{$sudo}}chmod 777 $__essh_var_docker_runfile

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
{{$sudo}}docker run \
    --privileged \
    -v $__essh_var_docker_working_dir:/essh \
    -v $__essh_var_docker_tmp_dir:/tmp/essh \
    -w /essh \
    $__essh_var_docker_image \
    bash /tmp/essh/run.sh
__essh_var_status=$?

echo "Removing tarminated containers."
{{$sudo}}docker rm `{{$sudo}}docker ps -a -q`

echo "Task exited with $__essh_var_status."
exit $__essh_var_status
]=]
end

return docker
