<?php

namespace tests\php;

/**
 * extract(..., EXTR_SKIP) 必须导入 ComponentSlot 类对象到当前作用域。
 */
class ExtractSkip_Slot
{
    public $attributes;
    public function __construct($html = '')
    {
        $this->html = $html;
        $this->attributes = (object) ['class' => 'x'];
    }
    public $html;
}

$data = [
    'trigger' => new ExtractSkip_Slot('OpenMe'),
    'slot' => new ExtractSkip_Slot('Body'),
    'placement' => null,
];

(function () use ($data) {
    extract($data, EXTR_SKIP);
    if (!isset($trigger) || !($trigger instanceof \tests\php\ExtractSkip_Slot)) {
        \Log::fatal('extract 未导入 trigger: '.var_export(isset($trigger)?$trigger:null, true));
    }
    if ($trigger->html !== 'OpenMe') {
        \Log::fatal('trigger.html 错误');
    }
    if (!isset($slot) || !($slot instanceof \tests\php\ExtractSkip_Slot)) {
        \Log::fatal('extract 未导入 slot');
    }
})();

// EXTR_SKIP 不覆盖已有
$trigger = 'EXISTING';
(function () use ($data) {
    // 注意：闭包内 $trigger 未定义，应导入
    extract($data, EXTR_SKIP);
    if (!($trigger instanceof \tests\php\ExtractSkip_Slot)) {
        \Log::fatal('闭包内应导入对象 trigger');
    }
})();

\Log::info('extract_skip_object_slot 测试通过');
