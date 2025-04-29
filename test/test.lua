-- extension strings
local strings = require("strings")

print(strings.ToUpper("abc"))

for i, s in ipairs(strings.Split("aa,b,,c", ",")) do
	print(i, s)
end

-- extension time
local time = require("time")

print("Wating 1 second")
time.sleep(1)

local begin = math.floor(time.unix())
print("begin:", begin)
local date, err = time.format(begin, "Jan 2 15:04:05 2006", "Europe/Moscow")
print("result:", date)
local result, err = time.parse(date, "Jan 2 15:04:05 2006")
print("result:", result)

-- extension tac
-- read file from the end to the beginning
local file = io.open("data.txt", "w")
file:write("1", "\n")
file:write("2", "\n")
file:write("3", "\n")

local tac = require("tac")
local scanner, err = tac.open("data.txt")
if err then error(err) end

while true do
    local line = scanner:line()
    if line == nil then break end
    print(line)
end
scanner:close()

local fs = require("fs")
-- remove file
fs.remove("data.txt")

-- extension storage

local storage = require("storage")

-- storage.open
local s, err = storage.open("db.json")
if err then error(err) end

-- storage:set(): key, value, ttl (default = 60s)
local err = s:set("key", {"one", "two", 1}, 10)
if err then error(err) end

-- storage:get()
local value, found, err = s:get("key")
if found then
    -- display all elements of the array
    for i, v in ipairs(value) do
        print(i, v)
    end
end

-- storage:set(): override with set max ttl
local err = s:set("key", "override", nil)
local value, found, err = s:get("key")
print("value:", value)

-- storage:keys()
local list = s:keys()
print("keys:", table.concat(list, ", "))


-- storage:dump()
local dump, err = s:dump()
print("dump:", dump)
for k, v in pairs(dump) do
    print(k, v)
end

fs.remove("db.json")

-- extension log
local log = require 'loglevel'

-- Logs go to /dev/null
log.DEBUG:print('foobar')
log.DEBUG:printf('foo bar %s', 'baz')

-- Logs go to some-file-on-disk (the default_output setting)
log.set_level('DEBUG') -- defaults to "INFO"
log.DEBUG:print('foobar')
log.DEBUG:printf('foo bar %s', 'baz')

-- Other loggers
log.DEBUG:print('debug')
log.INFO:print('info')
log.WARN:print('warn')
log.ERROR:print('error')

-- Log levels
log.set_level('DEBUG')
log.set_level('INFO')
log.set_level('WARN')
log.set_level('ERROR')

local log = require("log")
local info = log.new()

info:print("ok", " ", 1.2)
-- ok 1.2

info:println("ok", 1.2)
-- ok 1.2

info:printf("%s %f", "ok", 1.2)

-- extension crypto

local crypto = require("crypto")

print(crypto.md5("1\n"))
print(crypto.sha256("1\n"))

-- extension base64
local base64 = require("base64")

s = base64.RawStdEncoding:encode_to_string("foo\01bar")
print(s)
s = base64.StdEncoding:encode_to_string("foo\01bar")
print(s)
s = base64.RawURLEncoding:encode_to_string("this is a <tag> and should be encoded")
print(s)
s = base64.URLEncoding:encode_to_string("this is a <tag> and should be encoded")
print(s)

-- extension runtime
local runtime = require("runtime")
print(runtime.goos())
print(runtime.goarch())

-- extension inspect
local inspect = require("inspect")

local table = {"one", "two", {1, 2, 3}}
local result = inspect(table)
print(result)

-- extension aws
local aws = require("aws")

-- List EC2 instances by specifying region and profile
local instances = aws.listEC2Instances("eu-west-1", "default")

-- Iterate through instances
for _, instance in ipairs(instances) do
    -- Access instance properties
    print("Instance ID:", instance.instanceId)
    print("Instance Type:", instance.instanceType)
    print("State:", instance.state)
    
    -- Access IP addresses (if available)
    if instance.privateIp then
        print("Private IP:", instance.privateIp)
    end
    if instance.publicIp then
        print("Public IP:", instance.publicIp)
    end
    
    -- Access instance tags
    for tagKey, tagValue in pairs(instance.tags) do
        print("Tag:", tagKey, "=", tagValue)
    end
end

-- extension mdns

local mdns = require("mdns")

local results = mdns.browse("_workstation._tcp", "local", 2)

-- Results is a table of entries, each containing:
--   - name: service name
--   - service: service type
--   - domain: domain name
--   - port: port number
--   - ips: table of IP addresses
--   - serviceInstanceName: full service instance name
--   - serviceTypeName: full service type name

for _, entry in ipairs(results) do
    print("Found service:", entry.name)
    print("Port:", entry.port)
    -- output all IP addresses
    for _, ip in ipairs(entry.ips) do
        print("IP:", ip)
    end
    

    print("Service Type:", entry.service)
    print("Domain:", entry.domain)
    print("Service Instance Name:", entry.serviceInstanceName)
    print("Service Type Name:", entry.serviceTypeName)
end

-- extension watch

local watch = require("watch")

-- Watch a file for changes
local files = {
    "./test/"
}
local w, err = watch.watch(files, function(filename)
    print("File changed:", filename)
end)

