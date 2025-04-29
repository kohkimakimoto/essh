#!./esshdev --eval-file

print("Hello from Lua executed by Go!")

-- Print the arguments passed to the script
print("Arguments received by the Lua script:")
if arg then
    for i = 0, #arg do
        print("arg[" .. i .. "] = " .. tostring(arg[i]))
    end
else
    print("No arguments found in the 'arg' table.")
end

-- Example of using a Lua function
local function greet(name)
    return "Hello, " .. name .. "!"
end

print(greet("World"))

-- Example of error (uncomment to test error handling)
-- error("This is a test error in Lua")