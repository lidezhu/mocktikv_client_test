local rawkv = require 'rawkv'


--print('before create client')

--print(rawkv.double(2))
client = rawkv.newClient()
--print('after create client')
print('client', client["id"])
rawkv.put(client, 'foo2', 'bar2')
v = rawkv.get(client, 'foo2')
print('expect bar2', 'got', v)
--
rawkv.assertEquals(v, 'bar2')
rawkv.assertEquals(v, 'bar')