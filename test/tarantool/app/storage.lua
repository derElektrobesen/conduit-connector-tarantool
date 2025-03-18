#!/usr/bin/env tarantool

require('strict').on()

-- Get instance name
local fio = require('fio')
local NAME = fio.basename(arg[0], '.lua')
local fiber = require('fiber')

-- Call a configuration provider
cfg = dofile('localcfg.lua')
-- Name to uuid map
names = {
	['storage_1_a'] = '8a274925-a26d-47fc-9e1b-af88ce939412',
	['storage_1_b'] = '3de2e3e1-9ebe-4d0d-abb1-26d301b84633',
	['storage_1_c'] = '3de2e3e1-9ebe-4d0d-abb1-26d301b84635',
	['storage_2_a'] = '1e02ae8a-afc0-4e91-ba34-843a356b8ed7',
	['storage_2_b'] = '001688c3-66f8-4a31-8e19-036c17d489c2',
	['storage_2_c'] = 'c23516d5-22de-4ef4-8918-73d52e7661e2',
}

replicasets = {'cbf06940-0790-498b-948d-042b62cf3d29',
'ac522f65-aa94-4134-9f64-51ee384f1a54'}

local console = require 'console'
console.listen('127.0.0.1:3301')

-- Start the database with sharding
local vshard = require 'vshard'
rawset(_G, 'vshard', vshard) -- set as global variable

vshard.storage.cfg(cfg, names[NAME])

box.once('access:v1', function()
	box.schema.user.grant('guest', 'read,write,execute', 'universe')
end)

box.once("testapp:schema:1", function()
	local user = box.schema.space.create('user')
	user:format({
		{ 'ID', 'unsigned' },
		{ 'CAGId', 'number' },
		{ 'Domain', 'unsigned' },
		{ 'Username', 'string' },
		{ 'Flags', 'unsigned' },
		{ 'MailboxLimit', 'unsigned' },
		{ 'Forward', 'string' },
		{ 'RealName', 'string' },
		{ 'Comment', 'string' },
		{ 'Storage', 'string' },
		{ 'Target', 'string' },
		{ 'PwdCode', 'unsigned' },
		{ 'AttachStorage', 'string' },
		{ 'Flags2', 'unsigned' },
	})

	user:create_index('primary',    { parts = { 'ID' }, sequence = true, unique = true })
	user:create_index('user_uniq',  { parts = { 'Username', 'Domain' }, unique = true })
	user:create_index('domain',     { parts = { 'Domain', 'Storage' } })
	user:create_index('domain_att', { parts = { 'Domain', 'AttachStorage' } })

	local img = box.schema.space.create('security_image')
	img:format({
		{ 'id', 'number' },
		{ 'word', 'string' },
		{ 'time', 'string' },
		{ 'ip', 'unsigned' },
	})

	img:create_index('primary', { parts = { 'id' }, sequence = true, unique = true })
	img:create_index('tm',      { parts = { 'time' } })
	img:create_index('word',    { parts = { 'word' }, unique = true })

	box.snapshot()
end)

function put_tuple(space, tuple_map)
	return box.space[space]:put(box.space[space]:frommap(tuple_map))
end

function delete_tuple(space, key)
	return box.space[space]:delete(key)
end
