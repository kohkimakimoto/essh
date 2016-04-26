local template = require "essh.template"

local sshrc = function(override_config)
    config = {
        sshhome = "~",
    }
    if override_config ~= nil then
        -- merge
        for k,v in pairs(override_config) do config[k] = v end
    end

    return function()
        -- check existing xxd command.
        local status = os.execute("command -v xxd >/dev/null 2>&1 || exit 1")
        if status ~= 0 then
            error("sshrc requires xxd to be installed locally, but it's not. Aborting.")
        end

        essh.debug("sshrc sshhome: ".. config.sshhome)

        local command_v = [=[
function _essh_sshrc() {
    local SSHHOME=${SSHHOME:=]=] .. config.sshhome ..[=[}
    if [ -f $SSHHOME/.sshrc ]; then
        local files=.sshrc
        if [ -d $SSHHOME/.sshrc.d ]; then
            files="$files .sshrc.d"
        fi
        SIZE=$(tar cz -h -C $SSHHOME $files | wc -c)
        if [ $SIZE -gt 65536 ]; then
            echo >&2 $'.sshrc.d and .sshrc files must be less than 64kb\ncurrent size: '$SIZE' bytes'
            exit 1
        fi

        # echo $(tar cz -h -C $SSHHOME $files | xxd -p)
    else
        echo "No such file: $SSHHOME/.sshrc" >&2
        exit 1
    fi
}
_essh_sshrc
]=]
        -- validation
        local status = os.execute(command_v)
        if status ~= 0 then
            error("get a error by _essh_sshrc")
        end

        local command = [=[
function _essh_sshrc() {
    local SSHHOME=${SSHHOME:=]=] .. config.sshhome ..[=[}
    if [ -f $SSHHOME/.sshrc ]; then
        local files=.sshrc
        if [ -d $SSHHOME/.sshrc.d ]; then
            files="$files .sshrc.d"
        fi
        SIZE=$(tar cz -h -C $SSHHOME $files | wc -c)
        if [ $SIZE -gt 65536 ]; then
            echo >&2 $'.sshrc.d and .sshrc files must be less than 64kb\ncurrent size: '$SIZE' bytes'
            exit 1
        fi

        echo $(tar cz -h -C $SSHHOME $files | xxd -p)
    else
        echo "No such file: $SSHHOME/.sshrc" >&2
        exit 1
    fi
}
_essh_sshrc
]=]

        local f = io.popen(command)
        local result = f:read("*a")
        f:close()

        if result == nil or result == "" then
            error("empty result")
        end

        local dict = {
            sshrc_content = result,
        }

        local severside_script = template.dostring([=[
command -v xxd >/dev/null 2>&1 || { echo >&2 "sshrc requires xxd to be installed on the server, but it's not. Aborting."; exit 1; }
if [ -e /etc/motd ]; then cat /etc/motd; fi
if [ -e /etc/update-motd.d ]; then run-parts /etc/update-motd.d/ 2>/dev/null; fi
export SSHHOME=$(mktemp -d -t .$(whoami).sshrc.XXXX)
export SSHRCCLEANUP=$SSHHOME
trap "rm -rf $SSHRCCLEANUP; exit" 0

cat << 'EOF' > $SSHHOME/sshrc.bashrc
if [ -r /etc/profile ]; then source /etc/profile; fi
if [ -r ~/.bash_profile ]; then source ~/.bash_profile
elif [ -r ~/.bash_login ]; then source ~/.bash_login
elif [ -r ~/.profile ]; then source ~/.profile
fi
source $SSHHOME/.sshrc;
EOF

echo "{{.sshrc_content}}" | xxd -p -r | tar mxz -C $SSHHOME

bash --rcfile $SSHHOME/sshrc.bashrc
exit $?
]=], dict)

        return severside_script
    end
end

return sshrc
