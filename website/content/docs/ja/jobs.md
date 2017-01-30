+++
title = "ジョブ | ドキュメント"
type = "docs"
category = "docs"
lang = "ja"
basename = "jobs.html"
+++

# ジョブ

Jobs in Essh encapsulate tasks, hosts and drivers. Hosts and drivers that are defined in a job can be used only by the tasks in the same job.

## Example

~~~lua
job "myjob" {
    -- define description of the job
    description = "this is my job",
    
    -- If you set it true, defined tasks in this job are hidden at default.
    hidden = false,
    
    -- prepare function 
    prepare = function()
    
    end,
    
    -- defining hosts
    host "web01.localhost" {
        -- ...
    },

    host "web02.localhost" {
        -- ...
    },
    
    -- defining drivers
    driver "default" {
        -- ...
    },

    driver "custom" {
        -- ...
    },

    -- defining tasks
    task "foo" {
        -- ...
    },
    
    task "bar" {
        -- ...
    },
    
}
~~~

You can also define hosts, drivers and tasks by using tables. see below example:

~~~lua
job "myjob" {
    hosts = {
        ["web01.localhost"] = {
            -- ...
        },
        ["web02.localhost"] = {
            -- ...
        },
    },
    drivers = {
        ["default"] = {
            -- ...
        },
        ["custom"] = {
            -- ...
        },
    }
    tasks = {
        ["foo"] = {
            --- ...
        },
        ["bar"] = {
            --- ...
        }
    }
    
}
~~~

## Running job's task

A Job’s tasks have a prefix that is their job’s name, So you can run the task like the following

~~~
$ essh myjob:foo
~~~

## Default job 

If you define `default` job like the following. This job's tasks can be run without prefix.

~~~lua
job "default" {
    task "foo" {
        -- ...
    }
}
~~~

You can run the `foo` task without prefix.

~~~
$ essh foo
~~~


