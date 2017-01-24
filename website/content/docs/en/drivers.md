+++
title = "Drivers | Documentation"
type = "docs"
category = "docs"
lang = "en"
basename = "drivers.html"
+++

# Dirvers

Drivers in Essh are template system to construct shell scripts in tasks execution. You can use a driver to modify behavior of tasks.

## Example

~~~lua
driver "custom_driver" { 
    engine = [=[
    
    
    ]=],
    
    foo = "foo",
    
    bar = "bar",
}
~~~

WIP...