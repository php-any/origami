<?php
function redisCheck($condition, $message) {
    if (!$condition) { throw new RuntimeException($message); }
}
function redisFails($callback) {
    try { $callback(); } catch (RedisException $e) {
        redisCheck(strlen($e->getMessage()) > 0, 'empty RedisException');
        return;
    }
    throw new RuntimeException('expected RedisException');
}
redisCheck(extension_loaded('redis'), 'explicit registration');
redisCheck(in_array('redis', get_loaded_extensions()), 'extension list');
redisCheck(phpversion('redis') === REDIS_VERSION, 'module version');
$a = new Redis();
$b = new Redis();
redisCheck(!$a->isConnected(), 'constructor opened connection');
redisCheck($a->getMode() === Redis::ATOMIC, 'initial mode');
redisCheck($a->setOption(Redis::OPT_PREFIX, 'a:'), 'prefix option');
redisCheck($a->_prefix('key') === 'a:key', 'prefix');
redisCheck($b->getOption(Redis::OPT_PREFIX) === '', 'instance isolation');
redisCheck($a->setOption(Redis::OPT_READ_TIMEOUT, 0.25), 'read timeout');
redisCheck($a->getOption(Redis::OPT_READ_TIMEOUT) === 0.25, 'fractional timeout');
redisCheck($a->_pack("a\0b") === "a\0b", 'binary value');
redisCheck($a->getLastError() === null, 'initial error');
redisFails(function () use ($a) { $a->get('unconnected'); });
redisFails(function () use ($a) { $a->setOption(Redis::OPT_SERIALIZER, 999); });
redisFails(function () use ($a) { $a->setOption(Redis::OPT_READ_TIMEOUT, -1); });
redisFails(function () use ($a) { $a->pconnect('127.0.0.1', -1); });
redisFails(function () use ($a) { $a->connect('127.0.0.1', -1); });
redisFails(function () use ($a) { $a->rawCommand('HELLO', 3); });
redisCheck($a->getOption(Redis::OPT_READ_TIMEOUT) === 0.25, 'invalid option changed state');
echo "redis_driver_contract_test: OK\n";
