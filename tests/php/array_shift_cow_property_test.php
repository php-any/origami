<?php

namespace tests\php;

/**
 * PHP 数组 COW：`$parsed = $tokens` 后 array_shift($parsed) 不得清空 $tokens。
 * Symfony ArgvInput::bind → parse() 依赖此语义，否则 artisan 命令名丢失、退回 list。
 */

class ArrayShiftCow_Input
{
    private array $tokens;
    private array $parsed;

    public function __construct(array $tokens)
    {
        $this->tokens = $tokens;
    }

    public function bindShift(): ?string
    {
        $this->parsed = $this->tokens;
        $first = array_shift($this->parsed);
        if ($this->tokens === []) {
            return null;
        }
        return $first;
    }

    public function tokens(): array
    {
        return $this->tokens;
    }
}

$in = new ArrayShiftCow_Input(['inspire']);
$first = $in->bindShift();
if ($first !== 'inspire') {
    Log::fatal('array_shift 返回值错误: '.var_export($first, true));
}
if ($in->tokens() !== ['inspire']) {
    Log::fatal('array_shift 破坏了拷贝源 tokens: '.var_export($in->tokens(), true));
}

$a = ['artisan', 'serve'];
$b = $a;
array_shift($b);
if ($a !== ['artisan', 'serve']) {
    Log::fatal('局部数组 COW 失败: '.var_export($a, true));
}
if ($b !== ['serve']) {
    Log::fatal('array_shift 后副本不正确: '.var_export($b, true));
}

Log::info('array_shift_cow_property 测试通过');
