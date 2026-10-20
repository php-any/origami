<?php

namespace tests\php;

/**
 * tap() 的回调若按引用接收参数，修改必须反映到返回值上。
 */
$hash = '40bb37d8133165b2235301141bda11cf44da023d51649b13753338bb2156cd0a';
$obj = new \stdClass();
$map = new \WeakMap();
$map[$obj] = tap([], fn (&$stack) => $stack[$hash] = 'ok');
$stored = $map[$obj][$hash] ?? null;
if ($stored !== 'ok') {
    Log::fatal('tap 引用回调未写回: got='.var_export($stored, true));
}

$plain = tap(1, function ($n) {
    return $n + 1;
});
if ($plain !== 1) {
    Log::fatal('tap 应返回原值: got='.var_export($plain, true));
}

Log::info('tap_by_ref 测试通过');
