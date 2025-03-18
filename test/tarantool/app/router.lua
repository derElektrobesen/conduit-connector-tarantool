#!/usr/bin/env tarantool

require('strict').on()
local fiber = require('fiber')
rawset(_G, 'fiber', fiber) -- set fiber as global

local console = require 'console'
console.listen('127.0.0.1:3301')

box.cfg{
    listen = "127.0.0.1:12002"
}

box.once('access:v1', function()
    box.schema.user.grant('guest', 'read,write,execute', 'universe')
end)

replicasets = {'cbf06940-0790-498b-948d-042b62cf3d29',
               'ac522f65-aa94-4134-9f64-51ee384f1a54'}

-- Call a configuration provider
cfg = dofile('localcfg.lua')
if arg[1] == 'discovery_disable' then
    cfg.discovery_mode = 'off'
end

-- Start the database with sharding
local vshard = require 'vshard'
rawset(_G, 'vshard', vshard) -- set vshard as global

vshard.router.cfg(cfg)

vshard.router.bootstrap({
    timeout = 4,
    if_not_bootstrapped = true
})
