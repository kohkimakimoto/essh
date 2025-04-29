# This is zsh completion code.
# If you want to use it. write the following code in your '.zshrc'
#   eval "$(essh --zsh-completion)"
_essh_hosts() {
    local -a __essh_hosts
    PRE_IFS=$IFS
    IFS=$'\n'
    __essh_hosts=($({{.Executable}} --zsh-completion-hosts | awk -F'\t' '{print $1":"$2}'))
    IFS=$PRE_IFS
    _describe -t host "host" __essh_hosts
}

_essh_hosts_global() {
    local -a __essh_hosts
    PRE_IFS=$IFS
    IFS=$'\n'
    __essh_hosts=($({{.Executable}} --global --zsh-completion-hosts | awk -F'\t' '{print $1":"$2}'))
    IFS=$PRE_IFS
    _describe -t host "host" __essh_hosts
}

_essh_tasks() {
    local -a __essh_tasks
    PRE_IFS=$IFS
    IFS=$'\n'
    __essh_tasks=($({{.Executable}} --zsh-completion-tasks | awk -F'\t' '{print $1":"$2}'))
    IFS=$PRE_IFS
    _describe -t task "task" __essh_tasks
}

_essh_tasks_global() {
    local -a __essh_tasks
    PRE_IFS=$IFS
    IFS=$'\n'
    __essh_tasks=($({{.Executable}} --global --zsh-completion-tasks | awk -F'\t' '{print $1":"$2}'))
    IFS=$PRE_IFS
    _describe -t task "task" __essh_tasks
}

_essh_tags() {
    local -a __essh_tags
    PRE_IFS=$IFS
    IFS=$'\n'
    __essh_tags=($({{.Executable}} --zsh-completion-tags))
    IFS=$PRE_IFS
    _describe -t tag "tag" __essh_tags
}

_essh_tags_global() {
    local -a __essh_tags
    PRE_IFS=$IFS
    IFS=$'\n'
    __essh_tags=($({{.Executable}} --global --zsh-completion-tags))
    IFS=$PRE_IFS
    _describe -t tag "tag" __essh_tags
}

_essh_options() {
    local -a __essh_options
    __essh_options=(
        '--version:Print version.'
        '--help:Print help.'
        '--print:Print generated ssh config.'
        '--color:Force ANSI output.'
        '--no-color:Disable ANSI output.'
        '--gen:Only generate ssh config.'
        '--working-dir:Change working directory.'
        '--config:Load per-project configuration from the file.'
        '--hosts:List hosts.'
        '--tags:List tags.'
        '--tasks:List tasks.'
		'--eval:Evaluate lua script.'
		'--eval-file:Evaluate lua script from file.'
        '--debug:Output debug log.'
        '--global:Force using global config.'
        '--exec:Execute commands with the hosts.'
        '--zsh-completion:Output zsh completion code.'
        '--bash-completion:Output bash completion code.'
        '--aliases:Output aliases code.'
     )
    _describe -t option "option" __essh_options
}

_essh_hosts_options() {
    local -a __essh_options
    __essh_options=(
        '--debug:Output debug log.'
        '--quiet:Show only names.'
        '--all:Show all that includes hidden hosts.'
        '--select:Get only the hosts filtered with tags or hosts.'
        '--filter:Filter selected hosts with tags or hosts.'
        '--ssh-config:Output selected hosts as ssh_config format.'
     )
    _describe -t option "option" __essh_options
}

_essh_tasks_options() {
    local -a __essh_options
    __essh_options=(
        '--debug:Output debug log.'
        '--quiet:Show only names.'
        '--all:Show all that includes hidden tasks.'
     )
    _describe -t option "option" __essh_options
}

_essh_tags_options() {
    local -a __essh_options
    __essh_options=(
        '--debug:Output debug log.'
        '--quiet:Show only names.'
     )
    _describe -t option "option" __essh_options
}

_essh_exec_options() {
    local -a __essh_options
    __essh_options=(
        '--debug:Output debug log.'
        '--backend:Run the commands on local or remote hosts.'
        '--target:Target hosts to run the commands.'
        '--filter:Filter target hosts with tags or hosts.'
        '--prefix:Disable outputting prefix.'
        '--prefix-string:Custom string of the prefix.'
        '--privileged:Run by the privileged user.'
        '--user:Run by the specific user.'
        '--parallel:Run in parallel.'
        '--pty:Allocate pseudo-terminal. (add ssh option "-t -t" internally)'
        '--script-file:Load commands from a file.'
        '--driver:Specify a driver.'
     )
    _describe -t option "option" __essh_options
}

_essh_registry_options() {
    local -a __essh_options
    __essh_options=(
        '--with-global:Update or clean modules in the local, global both registry.'
     )
    _describe -t option "option" __essh_options
}

_essh_backends() {
    local -a __essh_options
    __essh_options=(
        'local'
        'remote'
     )
    _describe -t option "option" __essh_options
}

_essh () {
    local curcontext="$curcontext" state line
    local last_arg arg execMode hostsMode tasksMode tagsMode globalMode

    typeset -A opt_args

    _arguments \
        '1: :->objects' \
        '*: :->args' \
        && ret=0

    case $state in
        objects)
            case $line[1] in
                -*)
                    _essh_options
                    ;;
                *)
                    _essh_tasks
                    _essh_hosts
                    ;;
            esac
            ;;
        args)
            last_arg="${line[${#line[@]}-1]}"

            for arg in ${line[@]}; do
                case $arg in
                    --global)
                        globalMode="on"
                        ;;
                    --exec)
                        execMode="on"
                        ;;
                    --hosts)
                        hostsMode="on"
                        ;;
                    --tasks)
                        tasksMode="on"
                        ;;
                    --tags)
                        tagsMode="on"
                        ;;
                    *)
                        ;;
                esac
            done

            case $last_arg in
                --global)
                    if [ "$globalMode" = "on" ]; then
                      _essh_tasks_global
                      _essh_hosts_global
                    else
                      _essh_tasks
                      _essh_hosts
                    fi
                    ;;
                --print|--help|--version|--gen)
                    ;;
                --script-file|--config|--eval-file)
                    _files
                    ;;
                --select|--target|--filter)
                    if [ "$globalMode" = "on" ]; then
                      _essh_hosts_global
                      _essh_tags_global
                    else
                      _essh_hosts
                      _essh_tags
                    fi

                    ;;

                --backend)
                    _essh_backends
                    ;;
                *)
                    if [ "$execMode" = "on" ]; then
                        _essh_exec_options
                    elif [ "$hostsMode" = "on" ]; then
                        _essh_hosts_options
                    elif [ "$tasksMode" = "on" ]; then
                        _essh_tasks_options
                    elif [ "$tagsMode" = "on" ]; then
                        _essh_tags_options
                    else
                        _essh_options
                        _files
                    fi
                    ;;
            esac
            ;;
        *)
            _files
            ;;
    esac
}

compdef _essh essh