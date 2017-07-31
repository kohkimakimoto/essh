

-- print("Hello example")

-- print(essh.module.var.hoge)

host "webaaa-01" {
    description = essh.module.var.hoge,
}

host "webaaa-02" {
    description = essh.module.var.hoge,
}

task "task999" {
    script = "ls -la",
}


abc = "hogehgoe abc!!"