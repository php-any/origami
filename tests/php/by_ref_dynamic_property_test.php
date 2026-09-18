<?php

namespace tests\php;

/**
 * 用户函数 by-ref 必须接受 $obj->{$key} 动态属性（对齐 PHP / data_set）。
 */
class ByRefDynProp_Obj
{
    public array $data = [];
}

function by_ref_dyn_prop_set(&$target, $key, $value): void
{
    if (!is_array($target)) {
        $target = [];
    }
    $target[$key] = $value;
}

$obj = new ByRefDynProp_Obj();
$seg = 'data';
by_ref_dyn_prop_set($obj->{$seg}, 'email', 'a@b.c');

if (($obj->data['email'] ?? null) !== 'a@b.c') {
    Log::fatal('动态属性 by-ref 写回失败: ' . var_export($obj->data, true));
}

Log::info('by_ref_dynamic_property 测试通过');
