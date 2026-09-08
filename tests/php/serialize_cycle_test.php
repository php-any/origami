<?php

namespace tests\php;

/**
 * serialize 对对象自引用应输出 r:N;，不能死递归占内存。
 */
class SerializeCycle_Node
{
	public $me;
}

$n = new SerializeCycle_Node();
$n->me = $n;
$s = serialize($n);
if (!is_string($s) || !str_contains($s, 'r:1;')) {
	Log::fatal('循环引用 serialize 应包含 r:1; 实际: ' . var_export($s, true));
}
if (strlen($s) > 200) {
	Log::fatal('循环引用 serialize 结果过长，可能未截断环: ' . strlen($s));
}
Log::info('serialize 循环引用测试通过');
