package lock

var (
	// lockscript lua script for acrequire a lock
	lockLuaScript = `
if redis.call("GET", KEYS[1]) == ARGV[1] then
    redis.call("SET", KEYS[1], ARGV[1], "PX", ARGV[2])
    return "OK"
else
    return redis.call("SET", KEYS[1], ARGV[1], "NX", "PX", ARGV[2])
end
`

	// unlockscript lua script for release a lock
	unlockLuaScript = `
if redis.call("GET", KEYS[1]) == ARGV[1] then
    return redis.call("DEL", KEYS[1])
else
    return 0
end
`
)
