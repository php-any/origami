<?php

namespace tests\php;

/**
 * 复现 Livewire 布局 nowdoc：named argument 与 slots[-1] ?? []。
 */
class NamedArgSlot_Env
{
    public $last = null;
    public function slot($name, $attributes = null)
    {
        $this->last = [$name, $attributes];
        return 'ok';
    }
}

class NamedArgSlot_Attrs
{
    public function getAttributes()
    {
        return ['class' => 'x'];
    }
}

class NamedArgSlot_Slot
{
    public $attributes;
    public function __construct()
    {
        $this->attributes = new NamedArgSlot_Attrs();
    }
}

$__env = new NamedArgSlot_Env();
$name = 'sidebar';
$slot = new NamedArgSlot_Slot();
$__env->slot($name, attributes: $slot->attributes->getAttributes());
if ($__env->last[0] !== 'sidebar' || ($__env->last[1]['class'] ?? '') !== 'x') {
    Log::fatal('named argument attributes: 调用失败: '.var_export($__env->last, true));
}

$slots = [];
$n = 0;
foreach ($slots[-1] ?? [] as $k => $v) {
    $n++;
}
if ($n !== 0) {
    Log::fatal('slots[-1] ?? [] foreach 未按空数组处理');
}

Log::info('named arg 与负下标 foreach 测试通过');
