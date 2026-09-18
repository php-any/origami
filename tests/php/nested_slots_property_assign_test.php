<?php

namespace tests\php;

/**
 * 对象属性上的嵌套数组赋值：$this->slots[$i][$name] = ComponentSlot 形态。
 */
class NestedSlots_Bag
{
    public $slots = [];
    public $componentStack = ['a'];

    public function currentComponent()
    {
        return count($this->componentStack) - 1;
    }

    public function store($name, $value)
    {
        $this->slots[$this->currentComponent()][$name] = $value;
    }

    public function read($name)
    {
        return $this->slots[$this->currentComponent()][$name] ?? null;
    }
}

$o = new NestedSlots_Bag();
$o->store('trigger', ['html' => 'X', 'attributes' => ['class' => 'btn']]);
$got = $o->read('trigger');
if (!is_array($got) || ($got['html'] ?? null) !== 'X') {
    \Log::fatal('嵌套 slots 赋值失败: '.var_export($o->slots, true));
}

// 模拟 componentData merge 后 $$key
$data = array_merge(['placement' => 'bottom'], $o->slots[0]);
if (!isset($data['trigger'])) {
    \Log::fatal('merge 后无 trigger');
}

// 变量变量
$__key = 'trigger';
$__value = null;
$$__key = $$__key ?? $__value;
// 此时 $trigger 可能尚未定义，应变成 null——先设再测
$trigger = $data['trigger'];
$$__key = $$__key ?? $__value;
if (!is_array($trigger) || $trigger['html'] !== 'X') {
    \Log::fatal('变量变量覆盖坏了: '.var_export($trigger, true));
}

\Log::info('nested_slots_property_assign 测试通过');
