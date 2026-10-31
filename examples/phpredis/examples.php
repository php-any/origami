<?php
// Executable phpredis API examples, with strict checks and visible results.
// API reference: https://github.com/phpredis/phpredis/blob/6.1.0/README.md
// Requires a standalone Redis; only this run's randomly prefixed keys are removed.
function showExample($name, $actual, $expected) {
    if ($actual !== $expected) {
        throw new RuntimeException($name . ': unexpected result ' . json_encode($actual));
    }
    echo 'PASS ' . $name . ': ' . json_encode($actual) . "\n";
}
function probeExample($name, $callback) {
    try {
        $value = $callback();
        echo 'AVAILABLE ' . $name . ': ' . json_encode($value) . "\n";
    } catch (Throwable $e) {
        echo 'UNSUPPORTED ' . $name . ': ' . get_class($e) . ': ' . $e->getMessage() . "\n";
    }
}
$host = getenv('ORIGAMI_REDIS_TEST_HOST') ?: '127.0.0.1';
$port = (int) (getenv('ORIGAMI_REDIS_TEST_PORT') ?: 6379);
$redis = new Redis();
showExample('connect', $redis->connect($host, $port, 2.0, null, 0, 2.0), true);
$redis->setOption(Redis::OPT_PREFIX, 'origami-phpredis-example:' . uniqid('', true) . ':');
try {
    echo 'Redis server: ' . $redis->info('server')['redis_version'] . "\n";
    $redis->set('greeting', 'Hello Redis');
    showExample('SET/GET', $redis->get('greeting'), 'Hello Redis');
    showExample('missing GET', $redis->get('missing'), false);
    $redis->set('counter', 10);
    showExample('INCR', $redis->incr('counter'), 11);
    showExample('SET NX EX', $redis->set('lock', 'owner', ['nx', 'ex' => 60]), true);
    showExample('SET NX conflict', $redis->set('lock', 'other', ['nx']), false);
    $redis->mSet(['first' => 'A', 'second' => 'B']);
    showExample('MSET/MGET', $redis->mGet(['first', 'second']), ['A', 'B']);
    $redis->hMSet('profile', ['name' => 'Alice', 'role' => 'developer']);
    showExample('HMSET/HGETALL', $redis->hGetAll('profile'), ['name' => 'Alice', 'role' => 'developer']);
    $redis->rPush('jobs', 'email', 'report');
    showExample('RPUSH/LRANGE', $redis->lRange('jobs', 0, -1), ['email', 'report']);
    $redis->sAdd('tags', 'php', 'redis');
    showExample('SADD/SISMEMBER', $redis->sIsMember('tags', 'php'), true);
    $redis->zAdd('scores', 20, 'Bob', 10, 'Alice');
    showExample('ZADD/ZRANGE', $redis->zRange('scores', 0, -1, true), ['Alice' => 10.0, 'Bob' => 20.0]);
    showExample('Lua EVAL', $redis->eval('return redis.call("get", KEYS[1])', ['greeting'], 1), 'Hello Redis');
    showExample('pipeline', $redis->pipeline()->set('batch', 'one')->get('batch')->exec(), [true, 'one']);
    showExample('MULTI/EXEC', $redis->multi()->set('total', 1)->incr('total')->get('total')->exec(), [true, 2, '2']);

    // These execute real calls and print actual failures; they are capability
    // probes, not assertions that unsupported behavior is phpredis-compatible.
    probeExample('pconnect', function () use ($host, $port) {
        $persistent = new Redis();
        try { return $persistent->pconnect($host, $port, 2.0); }
        finally { $persistent->close(); }
    });
    probeExample('constructor options', function () use ($host, $port) {
        $configured = new Redis(['host' => $host, 'port' => $port]);
        try { return $configured->ping(); } finally { $configured->close(); }
    });
    probeExample('SCAN', function () use ($redis) {
        $cursor = null;
        return $redis->scan($cursor);
    });
    probeExample('PHP serializer', function () use ($redis) {
        return $redis->setOption(Redis::OPT_SERIALIZER, Redis::SERIALIZER_PHP);
    });
    probeExample('retry option', function () use ($redis) {
        return $redis->setOption(Redis::OPT_MAX_RETRIES, 3);
    });
} finally {
    $keys = $redis->keys('*');
    if ($keys) {
        $redis->setOption(Redis::OPT_PREFIX, '');
        $redis->del($keys);
    }
    $redis->close();
}
echo "Examples finished; temporary keys removed.\n";
