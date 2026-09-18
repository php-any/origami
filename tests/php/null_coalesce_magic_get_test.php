<?php

namespace tests\php;

/**
 * $obj->missing ?? $default 不得触发 Undefined array key（__get 内访问 $this->data[$key]）。
 */

class CoalesceMagic_ViewLike
{
    protected $data = ['x' => 1];

    public function &__get($key)
    {
        return $this->data[$key];
    }

    public function __isset($key)
    {
        return isset($this->data[$key]);
    }
}

$v = new CoalesceMagic_ViewLike();

$got = $v->layoutConfig ?? 'DEFAULT';
if ($got !== 'DEFAULT') {
    Log::fatal('?? 结果错误: ' . var_export($got, true));
}

// 对照：直接访问应 Warning（不 Fatal）
$direct = @$v->layoutConfig;
if ($direct !== null) {
    Log::fatal('直接 __get 缺失键应为 null: ' . var_export($direct, true));
}

Log::info('对象属性 ?? 不触发 __get 测试通过');
