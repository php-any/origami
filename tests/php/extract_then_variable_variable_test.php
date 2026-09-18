<?php

namespace tests\php;

/**
 * extract 导入的变量必须能被 $$name 读到（Blade @props 依赖此语义）。
 */
class ExtractVV_Slot
{
    public $html;
    public function __construct($html)
    {
        $this->html = $html;
    }
}

$data = ['trigger' => new ExtractVV_Slot('OpenMe')];

(function () use ($data) {
    extract($data, EXTR_SKIP);

    if (!isset($trigger)) {
        \Log::fatal('extract 后 isset($trigger) 为 false');
    }

    $__key = 'trigger';
    $__value = null;

    // 直接读
    $direct = $trigger;
    if (!($direct instanceof \tests\php\ExtractVV_Slot)) {
        \Log::fatal('直接读 $trigger 失败');
    }

    // 变量变量读
    $via = $$__key;
    if (!($via instanceof \tests\php\ExtractVV_Slot)) {
        \Log::fatal('$$__key 读不到 extract 变量: '.var_export($via, true).' type='.gettype($via));
    }

    // Blade @props 行
    $$__key = $$__key ?? $__value;
    if (!($trigger instanceof \tests\php\ExtractVV_Slot)) {
        \Log::fatal('@props 行后 $trigger 被破坏: '.var_export($trigger, true));
    }
})();

\Log::info('extract_then_variable_variable 测试通过');
