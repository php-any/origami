<?php

namespace tests\php;

/**
 * PHP：子类未声明构造函数时，new 与 ReflectionClass::getConstructor 都应使用父类构造函数。
 * Livewire Redirector 继承 Illuminate Redirector 依赖此语义才能注入 UrlGenerator。
 */

class InheritedCtor_Parent
{
    public $n;

    public function __construct(int $n)
    {
        $this->n = $n;
    }
}

class InheritedCtor_Child extends InheritedCtor_Parent
{
}

$ref = new \ReflectionClass(InheritedCtor_Child::class);
$ctor = $ref->getConstructor();
if ($ctor === null) {
    Log::fatal('子类 getConstructor 应为父类构造函数，得到 null');
}
$params = $ctor->getParameters();
if (count($params) !== 1 || $params[0]->getName() !== 'n') {
    Log::fatal('继承构造函数参数不正确: ' . json_encode(array_map(fn ($p) => $p->getName(), $params)));
}

$obj = new InheritedCtor_Child(7);
if (($obj->n ?? null) !== 7) {
    Log::fatal('new 子类应调用父类构造函数, n=' . var_export($obj->n ?? null, true));
}

Log::info('继承构造函数反射与 new 测试通过');
