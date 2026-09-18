<?php

namespace tests\php;

/**
 * 对齐 PHP：$this->stack[] = $v 应对属性数组做 [] 追加。
 */
class PropertyArrayAppend_Holder
{
    protected $buildStack = [];

    public function push($value)
    {
        $this->buildStack[] = $value;
        return $this->buildStack;
    }
}

$h = new PropertyArrayAppend_Holder();
$out = $h->push('hash-one');
if (!is_array($out) || count($out) !== 1 || $out[0] !== 'hash-one') {
    Log::fatal('protected 属性 [] 追加失败: ' . var_export($out, true));
}

$out = $h->push('hash-two');
if (count($out) !== 2 || $out[1] !== 'hash-two') {
    Log::fatal('第二次 [] 追加失败: ' . var_export($out, true));
}

Log::info('property_array_append 测试通过');
