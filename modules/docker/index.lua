local docker = {}

docker.driver = function()
    t = [=[
echo "Using docker driver engine."
    ]=]

    return t
end

return docker
