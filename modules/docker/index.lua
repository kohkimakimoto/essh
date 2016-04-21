local docker = {}

docker.driver = function()
    t = [=[
echo 'Starting {{.Task.Name}}'
echo 'Using docker driver engine.'
    ]=]

    return t
end

return docker
