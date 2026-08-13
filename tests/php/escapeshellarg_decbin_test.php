<?php

namespace tests\php;

/**
 * escapeshellarg / decbin 函数测试。
 */

$a = escapeshellarg('hello');
if ($a !== "'hello'") {
    Log::fatal('escapeshellarg basic failed: ' . $a);
}

$b = escapeshellarg("it's");
if ($b !== "'it'\\''s'") {
    Log::fatal('escapeshellarg quote failed: ' . $b);
}

$d = decbin(10);
if ($d !== '1010') {
    Log::fatal('decbin(10) failed: ' . $d);
}

$d0 = decbin(0);
if ($d0 !== '0') {
    Log::fatal('decbin(0) failed: ' . $d0);
}

Log::info('escapeshellarg/decbin 函数测试通过');
