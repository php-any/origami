<?php

namespace tests\php;

use function hash_init;
use function hash_update;
use function hash_update_stream;
use function hash_final;

/**
 * hash_init / hash_update / hash_update_stream / hash_final 增量哈希测试。
 */

$data = 'hello origami hash';
$expected = hash('md5', $data);

$ctx = hash_init('md5');
if (!is_resource($ctx)) {
    Log::fatal('hash_init 应返回 resource: ' . gettype($ctx));
}

$type = get_resource_type($ctx);
if (strpos($type, 'Hash') === false) {
    Log::fatal('get_resource_type 应含 Hash: ' . var_export($type, true));
}

hash_update($ctx, 'hello ');
hash_update($ctx, 'origami hash');
$got = hash_final($ctx);
if ($got !== $expected) {
    Log::fatal("hash_update/final 结果不一致: got={$got} expected={$expected}");
}

$tmp = __DIR__ . '/_hash_init_stream.tmp';
file_put_contents($tmp, $data);
$fp = fopen($tmp, 'rb');
$ctx2 = hash_init('md5');
$n = hash_update_stream($ctx2, $fp);
fclose($fp);
@unlink($tmp);

if ($n <= 0) {
    Log::fatal("hash_update_stream 字节数异常: {$n}");
}
$got2 = hash_final($ctx2);
if ($got2 !== $expected) {
    Log::fatal("hash_update_stream/final 结果不一致: got={$got2} expected={$expected}");
}

$ctx3 = hash_init('sha256');
hash_update($ctx3, $data);
$raw = hash_final($ctx3, true);

$ctx4 = hash_init('sha256');
hash_update($ctx4, $data);
$hexExpected = hash_final($ctx4);
if (bin2hex($raw) !== $hexExpected) {
    Log::fatal('hash_final(raw=true) 与 hex 不一致');
}

Log::info('hash 增量 API 测试通过');
