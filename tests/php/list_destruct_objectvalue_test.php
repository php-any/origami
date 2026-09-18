<?php

namespace tests\php;

/**
 * 列表解构 [$a,$b]=$arr 须支持 ObjectValue（字符串键 "0","1"），对齐 Collection::toArray()。
 */
$arr = ['0' => '$content', '1' => '$logo, $isDarkMode = false'];
echo "gettype=".gettype($arr)."\n";

[$name, $args] = $arr;
if ($name !== '$content') {
    Log::fatal('ObjectValue 列表解构 name 失败: ' . var_export($name, true));
}
if ($args !== '$logo, $isDarkMode = false') {
    Log::fatal('ObjectValue 列表解构 args 失败: ' . var_export($args, true));
}

$out = "<?php {$name} = (function (\$a) { return function ({$args}) use (\$a) {";
if (!str_contains($out, '<?php $content =')) {
    Log::fatal('插值失败: ' . json_encode($out));
}

Log::info('list_destruct_objectvalue 测试通过');
