local description = "Open essh repository"

if essh.module.var.description then
    description = essh.module.var.description
end

task "open-essh-repository" {
    script = "open https://github.com/kohkimakimoto/essh",
    description = description,
}

