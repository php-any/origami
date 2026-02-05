<?php

namespace tests\php;

echo "=== iconv() 函数测试 ===\n";

// 1. 相同编码，结果必须严格等于原字符串，且 !== false
$s = "Hello 世界";
$r1 = iconv('UTF-8', 'UTF-8', $s);
if ($r1 === $s && $r1 !== false) {
    Log::info("iconv 相同编码测试通过");
} else {
    Log::fatal("iconv 相同编码测试失败，期望: {$s}, 实际: " . var_export($r1, true));
}

// 2. 带 //TRANSLIT/IGNORE 等后缀，规范化后等价于 UTF-8，仍应严格等于原字符串，且 !== false
$r2 = iconv('utf-8', 'utf-8//IGNORE', $s);
if ($r2 === $s && $r2 !== false) {
    Log::info("iconv IGNORE 后缀测试通过");
} else {
    Log::fatal("iconv IGNORE 后缀测试失败，期望: {$s}, 实际: " . var_export($r2, true));
}

// 3. 目标不是 utf-8，目前实现必须返回 false，且结果 !== 原字符串
$r3 = iconv('UTF-8', 'ISO-8859-1', $s);
if ($r3 === false && $r3 !== $s) {
    Log::info("iconv 非 UTF-8 目标编码测试通过");
} else {
    Log::fatal("iconv 非 UTF-8 目标编码测试失败，期望: false, 实际: " . var_export($r3, true));
}

echo "=== iconv() 测试完成 ===\n";

