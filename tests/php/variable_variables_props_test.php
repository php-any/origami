<?php

namespace tests\php;

/**
 * 变量变量 $$name 读写，对齐 Blade @props 的 $$__key = $$__key ?? $__value。
 */
$trigger = 'KEEP_ME';
$__key = 'trigger';
$__value = null;

// 读取：$$__key 应为 KEEP_ME
$read = $$__key;
if ($read !== 'KEEP_ME') {
    \Log::fatal('变量变量读取失败: '.var_export($read, true));
}

// Blade @props 默认填充：已存在则保留
$$__key = $$__key ?? $__value;
if ($trigger !== 'KEEP_ME') {
    \Log::fatal('?? 后被覆盖: '.var_export($trigger, true));
}

// 未定义的变量
$__key2 = 'missingSlot';
$__value2 = 'DEFAULT';
$$__key2 = $$__key2 ?? $__value2;
if (!isset($missingSlot) || $missingSlot !== 'DEFAULT') {
    \Log::fatal('未定义变量变量赋默认失败');
}

// 对象槽位模拟
class VV_Slot
{
    public $attributes;
    public function __construct()
    {
        $this->attributes = (object) ['class' => 'btn'];
    }
}
$trigger = new VV_Slot();
$__key = 'trigger';
$__value = null;
$$__key = $$__key ?? $__value;
if (!($trigger instanceof VV_Slot)) {
    \Log::fatal('对象槽被 ?? 覆盖成: '.var_export($trigger, true));
}

\Log::info('variable_variables_props 测试通过');
