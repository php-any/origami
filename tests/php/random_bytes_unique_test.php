<?php

namespace tests\php;

/**
 * random_bytes 长度须接受 int 与算术结果（可能是 float），且每次调用内容不同。
 * Ramsey CombGenerator::generate(16-6) 依赖此语义。
 */

$a = random_bytes(16);
$b = random_bytes(16);
if (!is_string($a) || strlen($a) !== 16) {
    Log::fatal('random_bytes(16) 长度错误: '.strlen((string) $a));
}
if ($a === $b) {
    Log::fatal('连续 random_bytes(16) 不应相同');
}

$len = 16 - 6;
$c = random_bytes($len);
if (strlen($c) !== 10) {
    Log::fatal('random_bytes(16-6) 长度应为 10，实际 '.strlen($c).' type='.gettype($len));
}

$d = random_bytes(10.0);
if (strlen($d) !== 10) {
    Log::fatal('random_bytes(10.0) 长度应为 10，实际 '.strlen($d));
}

Log::info('random_bytes_unique 测试通过');
