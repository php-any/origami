<?php
namespace tests\basic;

/**
 * 关联数组字面量中的 ... 展开（Laravel ConfiguresPrompts 等场景）
 */

$opts = ['foo' => 'bar', 'baz' => 'qux'];
$result = ['' => 'None', ...$opts];

if (!isset($result['']) || $result[''] !== 'None') {
    Log::fatal("空字符串键 'None' 丢失");
}
if (!isset($result['foo']) || $result['foo'] !== 'bar') {
    Log::fatal("展开键 foo 失败");
}
if (!isset($result['baz']) || $result['baz'] !== 'qux') {
    Log::fatal("展开键 baz 失败");
}
Log::info("关联数组尾部展开测试通过");

$mixed = ['a' => 1, ...[10, 20], 'b' => 2];
if ($mixed['a'] !== 1 || $mixed['0'] !== 10 || $mixed['1'] !== 20 || $mixed['b'] !== 2) {
    Log::fatal("关联数组混合索引展开失败");
}
Log::info("关联数组混合索引展开测试通过");

$head = [...$opts, 'y' => 2];
if ($head['foo'] !== 'bar' || $head['y'] !== 2) {
    Log::fatal("展开后追加键值失败: " . json_encode($head));
}
Log::info("展开后追加键值测试通过");

$viaArray = array('' => 'None', ...$opts);
if ($viaArray[''] !== 'None' || $viaArray['foo'] !== 'bar') {
    Log::fatal("array() 关联展开失败");
}
Log::info("array() 关联展开测试通过");

Log::info("✅ 关联数组展开语法测试完成");
