<?php

namespace tests\php;

/**
 * unset($arr['strKey'])：通过 $arr['k']= 写入的字符串键必须能被删除
 *（StringValue 实现 AsInt，不能在 AsInt 失败后直接放弃）。
 */

$a = [];
$a['temp'] = 'temporary';
unset($a['temp']);
if (array_key_exists('temp', $a) || count($a) !== 0) {
    Log::fatal('赋值写入的字符串键 unset 失败: ' . var_export($a, true));
}

$b = [];
$b['temp'] = 't';
$b['other'] = 'x';
unset($b['temp']);
if (array_key_exists('temp', $b) || !array_key_exists('other', $b)) {
    Log::fatal('多键数组 unset 字符串键失败: ' . var_export($b, true));
}

class UnsetArrayProp_Holder
{
    public array $data = [];
}

$h = new UnsetArrayProp_Holder();
$h->data['temp'] = 'temporary';
unset($h->data['temp']);
if (isset($h->data['temp']) || array_key_exists('temp', $h->data)) {
    Log::fatal('对象属性数组 unset 失败: ' . var_export($h->data, true));
}

Log::info('unset_string_array_key 测试通过');
