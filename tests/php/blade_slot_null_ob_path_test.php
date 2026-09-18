<?php

namespace tests\php;

/**
 * Factory::slot('name', null, []) 必须走 ob 路径，不能因 !== null 误判把 null 写成 slot。
 */
class SlotNull_Env
{
    public $slots = [];
    public $slotStack = [];
    public $componentStack = ['comp'];

    public function currentComponent()
    {
        return count($this->componentStack) - 1;
    }

    public function slot($name, $content = null, $attributes = [])
    {
        if (func_num_args() === 2 || $content !== null) {
            $this->slots[$this->currentComponent()][$name] = $content;
            return 'assigned';
        } elseif (ob_start()) {
            $this->slots[$this->currentComponent()][$name] = '';
            $this->slotStack[$this->currentComponent()][] = [$name, $attributes];
            return 'buffered';
        }
        return 'fail';
    }

    public function endSlot()
    {
        $currentSlot = array_pop($this->slotStack[$this->currentComponent()]);
        [$currentName, $currentAttributes] = $currentSlot;
        $this->slots[$this->currentComponent()][$currentName] = [
            'html' => trim(ob_get_clean()),
            'attributes' => $currentAttributes,
        ];
    }
}

$e = new SlotNull_Env();
$mode = $e->slot('trigger', null, []);
if ($mode !== 'buffered') {
    \Log::fatal('slot(null) 应 buffered，实际: '.$mode.' slots='.var_export($e->slots, true));
}
echo 'IN_SLOT';
$e->endSlot();
$slot = $e->slots[0]['trigger'];
if (!is_array($slot) || $slot['html'] !== 'IN_SLOT') {
    \Log::fatal('endSlot 内容错误: '.var_export($slot, true));
}

// 显式内容路径
$e2 = new SlotNull_Env();
$mode2 = $e2->slot('trigger', 'STATIC', []);
if ($mode2 !== 'assigned') {
    \Log::fatal('slot(STATIC) 应 assigned，实际: '.$mode2);
}

\Log::info('blade_slot_null_ob_path 测试通过');
