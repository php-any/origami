<?php

namespace tests\php;

/**
 * 对齐 Illuminate\Container：类实现 ArrayAccess 时 $this->buildStack[] 仍应写属性数组。
 */
class ArrayAccessAppend_Container implements \ArrayAccess
{
    protected $buildStack = [];

    public function offsetExists(mixed $offset): bool
    {
        return false;
    }

    public function offsetGet(mixed $offset): mixed
    {
        return null;
    }

    public function offsetSet(mixed $offset, mixed $value): void
    {
    }

    public function offsetUnset(mixed $offset): void
    {
    }

    public function build($concrete)
    {
        $this->buildStack[] = spl_object_hash($concrete);
        return $this->buildStack;
    }
}

$c = new ArrayAccessAppend_Container();
$fn = function () {
    return 1;
};
$stack = $c->build($fn);
if (!is_array($stack) || count($stack) !== 1 || !is_string($stack[0]) || $stack[0] === '') {
    Log::fatal('ArrayAccess 类上属性 [] 追加失败: ' . var_export($stack, true));
}

Log::info('arrayaccess_property_append 测试通过');
