<?php

namespace tests\php;

/**
 * Laravel Eloquent getTable() 用 class_basename($this)。
 * $this 必须当对象：get_class 得到 FQCN，basename 为短类名。
 * 若当成 Object(App\Models\User) 字符串，basename 会变成 User)。
 */

function ClassBasenameThis_basename($class)
{
    $class = is_object($class) ? get_class($class) : $class;

    return basename(str_replace('\\', '/', $class));
}

class ClassBasenameThis_Box
{
    public function check(): void
    {
        if (!is_object($this)) {
            Log::fatal('is_object($this) 应为 true，实际 false');
        }
        $gc = get_class($this);
        if ($gc !== self::class) {
            Log::fatal('get_class($this) 错误: '.$gc);
        }
        if (strpos($gc, 'Object(') !== false) {
            Log::fatal('get_class($this) 不能是 AsString: '.$gc);
        }
        $base = ClassBasenameThis_basename($this);
        if ($base !== 'ClassBasenameThis_Box') {
            Log::fatal('class_basename($this) 错误: '.$base);
        }
        $fromStr = ClassBasenameThis_basename(self::class);
        if ($base !== $fromStr) {
            Log::fatal('对象与 FQCN 的 basename 应相同: '.$base.' vs '.$fromStr);
        }
    }
}

$box = new ClassBasenameThis_Box();
$box->check();

$fromInstance = ClassBasenameThis_basename($box);
if ($fromInstance !== 'ClassBasenameThis_Box') {
    Log::fatal('对实例 class_basename 错误: '.$fromInstance);
}

Log::info('class_basename_this 测试通过');
