<?php
function featureCheck($ok, $message) {
    if (!$ok) { throw new RuntimeException($message); }
}
$host = getenv('ORIGAMI_REDIS_TEST_HOST') ?: '127.0.0.1';
$port = (int) getenv('ORIGAMI_REDIS_TEST_PORT');
if ($port <= 0) { throw new RuntimeException('ORIGAMI_REDIS_TEST_PORT required'); }
$r = new Redis(['host' => $host, 'port' => $port, 'connectTimeout' => 2, 'readTimeout' => 2,
    'backoff' => ['algorithm' => Redis::BACKOFF_ALGORITHM_CONSTANT, 'base' => 1, 'cap' => 5]]);
featureCheck(!$r->isConnected(), 'constructor must connect lazily');
featureCheck($r->ping() === true && $r->isConnected(), 'lazy connection');
featureCheck($r->getOption(Redis::OPT_BACKOFF_BASE) === 1, 'constructor backoff');
featureCheck($r->setOption(Redis::OPT_MAX_RETRIES, 3), 'set retries');
featureCheck($r->getOption(Redis::OPT_MAX_RETRIES) === 3, 'read retries');
featureCheck($r->setOption(Redis::OPT_MAX_RETRIES, -1) === false, 'reject invalid retries');
featureCheck($r->getOption(Redis::OPT_MAX_RETRIES) === 3, 'invalid option mutated state');
$prefix = 'origami-redis-features:' . uniqid('', true) . ':';
$r->setOption(Redis::OPT_PREFIX, $prefix);
try {
    $p = new Redis();
    try {
        featureCheck($p->pconnect($host, $port, 2, 'feature-test', 0, 2), 'pconnect');
        featureCheck($p->getPersistentID() === 'feature-test', 'persistent ID');
        featureCheck($p->ping() === true, 'persistent socket usable');
        featureCheck($p->client('ID') !== $r->client('ID'), 'active sockets must not be shared');
    } finally { $p->close(); }

    for ($i = 0; $i < 30; $i++) { $r->set('scan:' . $i, $i); }
    $cursor = null;
    $seen = [];
    do {
        $keys = $r->scan($cursor, $prefix . 'scan:*', 2, 'string');
        featureCheck(is_array($keys), 'SCAN result');
        foreach ($keys as $key) { $seen[$key] = true; }
    } while ($cursor !== 0);
    featureCheck(count($seen) === 30, 'SCAN lost keys');
    featureCheck($r->scan($cursor) === false, 'completed cursor must not restart');
    $r->setOption(Redis::OPT_SCAN, Redis::SCAN_RETRY);
    $r->setOption(Redis::OPT_SCAN, Redis::SCAN_PREFIX);
    featureCheck($r->getOption(Redis::OPT_SCAN) === 3, 'SCAN flags combine');
    $holder = ['cursor' => null];
    $all = [];
    do {
        $batch = $r->scan($holder['cursor'], 'scan:*', 3);
        foreach ($batch as $key) { $all[$key] = true; }
    } while ($holder['cursor'] !== 0);
    featureCheck(count($all) === 30, 'array-slot reference cursor');
    $r->setOption(Redis::OPT_SCAN, Redis::SCAN_NOPREFIX);
    featureCheck($r->getOption(Redis::OPT_SCAN) === 1, 'SCAN_NOPREFIX preserves RETRY');

    featureCheck($r->setOption(Redis::OPT_SERIALIZER, Redis::SERIALIZER_PHP), 'PHP serializer');
    $payload = ['text' => "a\0b", 'list' => [1, true, null, false]];
    featureCheck($r->_serialize($payload) === serialize($payload), 'wire format');
    featureCheck($r->_unserialize('plain text') === 'plain text', 'raw strings retained');
    featureCheck($r->_unserialize('b:0;') === false, 'serialized false');
    $r->set('payload', $payload);
    featureCheck($r->get('payload') === $payload, 'SET/GET serialization');
    featureCheck($r->rawCommand('GET', $prefix . 'payload') === serialize($payload), 'server stores PHP format');
    featureCheck($r->get('missing') === false, 'missing key remains false');
    $r->mSet(['m1' => $payload, 'm2' => null]);
    featureCheck($r->mGet(['m1', 'm2', 'missing']) === [$payload, null, false], 'MGET serialization');
    $r->hMSet('hash', ['x' => $payload, 'y' => false]);
    featureCheck($r->hGetAll('hash') === ['x' => $payload, 'y' => false], 'hash serialization');
    featureCheck($r->hMGet('hash', ['x', 'missing']) === ['x' => $payload, 'missing' => false], 'HMGET serialization');
    $r->rPush('list', $payload, false);
    featureCheck($r->lRange('list', 0, -1) === [$payload, false], 'list serialization');
    $r->sAdd('set', $payload);
    featureCheck($r->sIsMember('set', $payload), 'serialized set membership');
    featureCheck($r->sMembers('set') === [$payload], 'set decoding');
    featureCheck($r->pipeline()->set('pipe', $payload)->get('pipe')->exec() === [true, $payload], 'serialized pipeline');
    featureCheck($r->multi()->set('tx', $payload)->get('tx')->exec() === [true, $payload], 'serialized transaction');
    class RedisFeatureObject {
        public $value;
        public function __serialize() { return ['value' => $this->value]; }
        public function __unserialize($data) { $this->value = $data['value']; }
    }
    $obj = new RedisFeatureObject(); $obj->value = $payload;
    $r->set('object', $obj);
    $copy = $r->get('object');
    featureCheck($copy instanceof RedisFeatureObject && $copy !== $obj && $copy->value === $payload, 'object hydration');
    class RedisReentrantSerialization {
        public $redis;
        public function __serialize() { $this->redis->ping(); return ['ok' => true]; }
    }
    $reentrant = new RedisReentrantSerialization(); $reentrant->redis = $r;
    featureCheck($r->set('reentrant', $reentrant), 'serialization reentrancy');
    class RedisThrowingSerialization {
        public function __serialize() { throw new LogicException('serialize failure'); }
    }
    try { $r->set('throwing', new RedisThrowingSerialization()); throw new RuntimeException('exception swallowed'); }
    catch (LogicException $e) { featureCheck($e->getMessage() === 'serialize failure', 'exception identity'); }
    featureCheck($r->get('throwing') === false, 'failed serialization sent a write');
} finally {
    $keys = $r->keys('*');
    if ($keys) { $r->setOption(Redis::OPT_PREFIX, ''); $r->del($keys); }
    $r->close();
}
echo "redis_phpredis_features_test: OK\n";
