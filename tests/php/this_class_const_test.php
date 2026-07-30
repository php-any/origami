<?php

namespace tests\php;

/**
 * $this::class / $obj::class 应返回类名字符串，且不因对象 AsString dump 触发类加载 fatal。
 */

class ThisClassConst_Model
{
    public $timestamps = true;

    public function usesTimestampsLike()
    {
        // 模拟 Eloquent HasTimestamps::usesTimestamps
        return $this->timestamps && !static::isIgnoringTimestamps($this::class);
    }

    public static function isIgnoringTimestamps($class = null)
    {
        $class ??= static::class;
        $ignored = [];
        foreach ($ignored as $ignoredClass) {
            if ($class === $ignoredClass || is_subclass_of($class, $ignoredClass)) {
                return true;
            }
        }
        return false;
    }

    public function className()
    {
        return $this::class;
    }
}

$m = new ThisClassConst_Model();
$name = $m->className();
$expect = ThisClassConst_Model::class;
if ($name !== $expect && $name !== 'tests\\php\\ThisClassConst_Model') {
    Log::fatal('$this::class 期望类名，得到: ' . $name);
}

if ($m->usesTimestampsLike() !== true) {
    Log::fatal('usesTimestampsLike 期望 true');
}

$viaVar = $m;
$n2 = $viaVar::class;
if ($n2 !== $expect && $n2 !== 'tests\\php\\ThisClassConst_Model') {
    Log::fatal('$obj::class 期望类名，得到: ' . $n2);
}

Log::info('$this::class 测试通过');
