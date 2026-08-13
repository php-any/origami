<?php

namespace tests\php;

/**
 * foreach 对象时只遍历 public 属性（与 PHP 一致；Symfony TableSeparator 依赖此点）。
 */

class ForeachVisibility_PrivateProps
{
    public string $pub = 'public';
    protected string $pro = 'protected';
    private string $priv = 'private';
}

$obj = new ForeachVisibility_PrivateProps();
$keys = [];
foreach ($obj as $k => $v) {
    $keys[] = $k;
}

if ($keys !== ['pub']) {
    Log::fatal('foreach should only see public props, got ' . var_export($keys, true));
}

Log::info('foreach_object_visibility_test 测试通过');
