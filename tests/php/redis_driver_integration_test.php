<?php
function redisAssert($condition, $message) {
    if (!$condition) { throw new RuntimeException($message); }
}
$port = (int) getenv('ORIGAMI_REDIS_TEST_PORT');
if ($port <= 0) { throw new RuntimeException('ORIGAMI_REDIS_TEST_PORT is required'); }
$host = getenv('ORIGAMI_REDIS_TEST_HOST') ?: '127.0.0.1';
$r = new Redis();
$other = new Redis();
redisAssert($r->connect($host, $port, 2.0, null, 0, 2.0), 'connect');
redisAssert($other->connect($host, $port, 2.0, null, 0, 2.0), 'second connection');
$prefix = 'origami-redis-test:' . uniqid('', true) . ':';
$r->setOption(Redis::OPT_PREFIX, $prefix);
$other->setOption(Redis::OPT_PREFIX, $prefix);
try {
    redisAssert($r->ping() === true && $r->ping('echo') === 'echo', 'ping');
    redisAssert($r->get('absent') === false, 'missing key');
    redisAssert($r->set('empty', '') && $r->get('empty') === '', 'empty value');
    redisAssert($r->set('bin', "a\0b\r\n") && $r->get('bin') === "a\0b\r\n", 'binary value');
    redisAssert($r->rawCommand('GET', $prefix . 'bin') === "a\0b\r\n", 'raw command must not prefix');
    redisAssert($r->set('lock', 'one', ['nx', 'px' => 5000]) === true, 'SET NX PX');
    redisAssert($r->set('lock', 'two', ['nx']) === false, 'SET NX failure');
    redisAssert($r->pttl('lock') > 0, 'expiry');
    redisAssert($r->set('lock', 'two', ['xx', 'get']) === 'one', 'SET GET old value');
    redisAssert($r->mSet(['a' => 'A', 'b' => 'B']), 'mset');
    redisAssert($r->mGet(['a', 'b', 'missing']) === ['A', 'B', false], 'mget');
    redisAssert($r->incr('counter') === 1 && $r->incr('counter', 2) === 3, 'counter');
    redisAssert($r->type('counter') === Redis::REDIS_STRING, 'type');
    redisAssert($r->hMSet('hash', ['x' => '1', 'y' => '2']), 'hash set');
    redisAssert($r->hMGet('hash', ['x', 'missing']) === ['x' => '1', 'missing' => false], 'hash missing fields');
    redisAssert($r->hGetAll('hash') === ['x' => '1', 'y' => '2'], 'hash mapping');
    redisAssert($r->get('hash') === false, 'WRONGTYPE returns false');
    redisAssert(strpos($r->getLastError(), 'WRONGTYPE') === 0, 'last error');
    redisAssert($r->ping() === true, 'error must not desynchronize connection');
    redisAssert($r->clearLastError() && $r->getLastError() === null, 'clear error');
    redisAssert($r->rPush('list', 'x', 'y', 'x') === 3, 'rpush');
    redisAssert($r->lRem('list', 'x', 1) === 1, 'PHP lrem argument order');
    redisAssert($r->lRange('list', 0, -1) === ['y', 'x'], 'list');
    redisAssert($r->sAdd('set', 'x', 'y') === 2 && $r->sIsMember('set', 'x'), 'set');
    redisAssert($r->zAdd('sorted', 2, 'b', 1, 'a') === 2, 'zadd');
    redisAssert($r->zRange('sorted', 0, -1, true) === ['a' => 1.0, 'b' => 2.0], 'scores');
    redisAssert($r->eval('return {KEYS[1], ARGV[1]}', ['bin', 'arg'], 1) === [$prefix . 'bin', 'arg'], 'Lua prefix only keys');
    redisAssert($r->eval('return redis.call("get", KEYS[1])', ['bin'], 1) === "a\0b\r\n", 'Lua value');
    redisAssert($r->pipeline() === $r, 'pipeline identity');
    redisAssert($r->set('pipe', 'p') === $r, 'queue identity');
    $r->get('pipe');
    $r->get('hash');
    $r->ping();
    redisAssert($other->get('pipe') === false, 'pipeline executed before exec');
    redisAssert($r->exec() === [true, 'p', false, true], 'pipeline results');
    redisAssert($r->ping() === true, 'pipeline error draining');
    $r->multi()->set('tx', 'v')->get('tx');
    redisAssert($other->get('tx') === false, 'transaction visible before exec');
    redisAssert($r->exec() === [true, 'v'], 'transaction results');
    $r->multi()->set('discard', 'v');
    redisAssert($r->discard() && $r->get('discard') === false, 'discard');
    $r->pipeline()->set('discard', 'v');
    redisAssert($r->discard() && $r->get('discard') === false, 'discard pipeline');
    $r->set('watch', 'before');
    $r->watch('watch');
    $other->set('watch', 'changed');
    $r->multi()->set('watch', 'wrong');
    redisAssert($r->exec() === false && $r->get('watch') === 'changed', 'WATCH abort');
    redisAssert($r->getMode() === Redis::ATOMIC, 'mode reset');
    // The second connection starts in DB 0. SELECT on the first must remain
    // isolated, including when queued in a pipeline.
    redisAssert($r->select(1) === true && $r->getDbNum() === 1, 'select database');
    try {
        redisAssert($r->get('bin') === false, 'database isolation');
        $r->set('db-one', 'v');
        redisAssert($other->get('db-one') === false, 'connection database isolation');
        redisAssert($r->del('db-one') === 1, 'db-one cleanup');
    } finally { $r->select(0); }
    redisAssert($r->pipeline()->select(1)->get('bin')->select(0)->exec() === [true, false, true], 'pipeline SELECT');
    redisAssert($r->getDbNum() === 0 && $r->get('bin') === "a\0b\r\n", 'pipeline database tracking');
    // Authentication is tested only when the caller explicitly enables ACL
    // administration for its disposable server.
    if (getenv('ORIGAMI_REDIS_TEST_ACL') === '1') {
        $user = 'origami-' . uniqid();
        $password = 'test-' . uniqid('', true);
        $r->rawCommand('ACL', 'SETUSER', $user, 'on', '>' . $password, '~' . $prefix . '*', '+@all');
        try {
            try { $other->auth([$user, 'incorrect']); throw new RuntimeException('bad AUTH accepted'); }
            catch (RedisException $e) { redisAssert(strpos($e->getMessage(), 'WRONGPASS') === 0, 'AUTH error'); }
            redisAssert($other->auth([$user, $password]) === true, 'ACL AUTH');
            redisAssert($other->get('bin') === "a\0b\r\n", 'authenticated read');
        } finally { $r->rawCommand('ACL', 'DELUSER', $user); }
    }
    redisAssert($r->info('server')['redis_version'] !== '', 'server version');
} finally {
    // Only remove keys with this invocation's random prefix; never flush a DB.
    $keys = $r->keys('*');
    if ($keys) { $r->setOption(Redis::OPT_PREFIX, ''); $r->del($keys); }
    $r->close();
    $other->close();
}
redisAssert(!$r->isConnected(), 'close');
echo "redis_driver_integration_test: OK\n";
