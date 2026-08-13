<?php
namespace tests\php;

class StaticByRef_Util
{
    public static function set(&$array, $key, $value)
    {
        $array[$key] = $value;
    }
}

$a = ['x' => 1];
StaticByRef_Util::set($a, 'y', 2);
if (($a['y'] ?? null) !== 2) {
    Log::fatal('static byref 失败: '.var_export($a, true));
}

Log::info('static_byref_array 测试通过');
