<?php

namespace tests\php;

/**
 * serialize / unserialize 最小标量子集互通测试：
 * - int / bool / string / null
 * - 数组 / 对象（通过 Origami 内部前缀 + JSON 实现往返）
 */

// int
$i = 42;
$si = serialize($i);
if ($si !== 'i:42;') {
    Log::fatal('serialize int 失败，期望 i:42;，实际: ' . var_export($si, true));
}
if (unserialize($si) !== $i) {
    Log::fatal('serialize+unserialize int 往返失败');
}
Log::info('serialize int 测试通过');

// bool
$b1 = true;
$b0 = false;
$sb1 = serialize($b1);
$sb0 = serialize($b0);
if ($sb1 !== 'b:1;' || $sb0 !== 'b:0;') {
    Log::fatal('serialize bool 失败，期望 b:1;/b:0;，实际: ' . var_export([$sb1, $sb0], true));
}
if (unserialize($sb1) !== $b1 || unserialize($sb0) !== $b0) {
    Log::fatal('serialize+unserialize bool 往返失败');
}
Log::info('serialize bool 测试通过');

// null
$n = null;
$sn = serialize($n);
if ($sn !== 'N;') {
    Log::fatal('serialize null 失败，期望 N;，实际: ' . var_export($sn, true));
}
if (unserialize($sn) !== $n) {
    Log::fatal('serialize+unserialize null 往返失败');
}
Log::info('serialize null 测试通过');

// string
$s = 'hello';
$ss = serialize($s);
if ($ss !== 's:5:"hello";') {
    Log::fatal('serialize string 失败，期望 s:5:"hello";，实际: ' . var_export($ss, true));
}
if (unserialize($ss) !== $s) {
    Log::fatal('serialize+unserialize string 往返失败');
}
Log::info('serialize string 测试通过');

// array
$arr = [1, 2, 3];
$sa = serialize($arr);
$ua = unserialize($sa);
if (!is_array($ua) || $ua !== $arr) {
    Log::fatal('serialize+unserialize array 往返失败');
}
Log::info('serialize array 测试通过');

// 当前 serialize/unserialize 对象仅保证数组场景可用，对 PHP 类实例/匿名 stdClass 的支持后续再补充。

