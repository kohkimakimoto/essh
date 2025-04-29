# This is zsh completion code.
# If you want to use it. write the following code in your '.zshrc'
#   eval "$(essh --bash-completion)"

_essh_hosts() {
    COMPREPLY=( $(compgen -W "$({{.Executable}} --bash-completion-hosts)" -- $cur) )
}

_essh_tasks() {
    COMPREPLY=( $(compgen -W "$({{.Executable}} --bash-completion-tasks)" -- $cur) )
}

_essh_hosts_and_tasks() {
    COMPREPLY=( $(compgen -W "$({{.Executable}} --bash-completion-hosts) $({{.Executable}} --bash-completion-tasks)" -- $cur) )
}

_essh_hosts_and_tags() {
    COMPREPLY=( $(compgen -W "$({{.Executable}} --bash-completion-hosts) $({{.Executable}} --bash-completion-tags)" -- $cur) )
}

_essh_registry_options() {
    COMPREPLY=( $(compgen -W "
        --with-global
    " -- $cur) )
}

_essh_backends() {
    COMPREPLY=( $(compgen -W "
        local
        remote
    " -- $cur) )
}

_essh_hosts_options() {
    COMPREPLY=( $(compgen -W "
        --debug
        --quiet
        --all
        --select
        --filter
        --ssh-config
    " -- $cur) )
}

_essh_tasks_options() {
    COMPREPLY=( $(compgen -W "
        --debug
        --quiet
        --all
    " -- $cur) )
}

_essh_tags_options() {
    COMPREPLY=( $(compgen -W "
        --debug
        --quiet
    " -- $cur) )
}

_essh_exec_options() {
    local -a __essh_options
    __essh_options=(
        '--debug:Output debug log.'
		'--eval:Evaluate lua script.'
        '--backend:Run the commands on local or remote hosts.'
        '--target:Target hosts to run the commands.'
        '--filter:Filter target hosts with tags or hosts.'
        '--prefix:Disable outputing prefix.'
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

_essh_options() {
    COMPREPLY=( $(compgen -W "
        --version
        --help
        --print
        --color
        --no-color
        --gen
        --global
        --working-dir
        --config
        --hosts
        --tags
        --tasks
		--eval
		--eval-file
        --debug
        --exec
        --zsh-completion
        --bash-completion
        --aliases
    " -- $cur) )
}

_essh() {
    COMP_WORDBREAKS=${COMP_WORDBREAKS//:}

    local last_arg arg execMode hostsMode tasksMode tagsMode

    local cur=${COMP_WORDS[COMP_CWORD]}
    case "$COMP_CWORD" in
        1)
            case "$cur" in
                -*)
                    _essh_options
                    ;;
                *)
                    _essh_hosts_and_tasks
                    ;;
            esac
            ;;
        *)
            last_arg="${COMP_WORDS[COMP_CWORD-1]}"
            for arg in ${COMP_WORDS[@]}; do
                case $arg in
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

            case "$last_arg" in
                --print|--help|--version|--gen)
                    ;;
                --script-file|--config|--eval-file)
                    ;;
                --select|--target|--filter)
                    _essh_hosts_and_tags
                    ;;
                --backend)
                    _essh_backends
                    ;;
                *)
                    if [ "$execMode" = "on" ]; then
                        _essh_hosts
                    elif [ "$hostsMode" = "on" ]; then
                        _essh_hosts_options
                    elif [ "$tasksMode" = "on" ]; then
                        _essh_tasks_options
                    elif [ "$tagsMode" = "on" ]; then
                        _essh_tags_options
                    else
                        _essh_options
                    fi
                    ;;
            esac
            ;;

    esac
}

complete -o default -o nospace -F _essh essh