<?php

namespace tests\php;

$values = [];
$first = new \stdClass;
$second = new \stdClass;

array_splice($values, 0, 0, [$first]);
array_splice($values, 1, 0, [$second]);

if (count($values) !== 2 || $values[0] !== $first || $values[1] !== $second) {
    Log::fatal('array_splice 未按引用修改空数组');
}

function spliceLocalArray(array $parameters, $value)
{
    array_splice($parameters, 0, 0, [$value]);

    return $parameters;
}

$local = spliceLocalArray([], $first);
if (count($local) !== 1 || $local[0] !== $first) {
    Log::fatal('array_splice 未按引用修改函数局部数组参数');
}

class ArraySpliceReferenceTarget
{
    protected function splice(array &$parameters, $value)
    {
        array_splice($parameters, 0, 0, [$value]);
    }

    public function resolve(array $parameters, $value)
    {
        $this->splice($parameters, $value);

        return $parameters;
    }
}

$methodLocal = (new ArraySpliceReferenceTarget)->resolve([], $first);
if (count($methodLocal) !== 1 || $methodLocal[0] !== $first) {
    Log::fatal('array_splice 未通过实例方法引用修改调用方数组');
}

Log::info('array_splice_reference 测试通过');
