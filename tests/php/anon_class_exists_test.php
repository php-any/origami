<?php

namespace tests\php;

/**
 * PHP：匿名类 get_class 后 class_exists 应为 true，且可 new $name()。
 */
$base = new class {
    public $v = 1;
};
$name = get_class($base);
if ($name === '' || $name === false) {
    \Log::fatal('get_class(anonymous) 为空');
}
if (!class_exists($name)) {
    \Log::fatal('class_exists(匿名类名) 应为 true, name='.$name);
}

$again = new $name();
if (!is_object($again)) {
    \Log::fatal('new $anonName 失败');
}
if (!($again instanceof $name)) {
    \Log::fatal('instanceof 匿名类名失败');
}

// 带父类
class AnonClassExists_Parent {}
$child = new class extends AnonClassExists_Parent {
    public $x = 2;
};
$cn = get_class($child);
if (!class_exists($cn)) {
    \Log::fatal('extends 匿名类 class_exists 失败: '.$cn);
}
if (!is_subclass_of($cn, AnonClassExists_Parent::class)) {
    \Log::fatal('is_subclass_of 匿名类失败');
}
$c2 = new $cn();
if (!($c2 instanceof AnonClassExists_Parent)) {
    \Log::fatal('new 匿名子类 instanceof parent 失败');
}

\Log::info('anon_class_exists 测试通过');
