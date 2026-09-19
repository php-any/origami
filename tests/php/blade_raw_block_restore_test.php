<?php

namespace tests\php;

/**
 * 对齐 BladeCompiler：@php 抽出 raw block 后须用 preg_replace_callback
 * 把 @__raw_block_N__@ 还原，且回调里 $this->rawBlocks[$matches[1]] 可用。
 */

class BladeRawBlock_Compiler
{
    protected $rawBlocks = [];

    public function getRawPlaceholder($replace)
    {
        return str_replace('#', $replace, '@__raw_block_#__@');
    }

    public function storeRawBlock($value)
    {
        return $this->getRawPlaceholder(
            array_push($this->rawBlocks, $value) - 1
        );
    }

    public function storePhpBlocks($value)
    {
        return preg_replace_callback('/(?<!@)@php(.*?)@endphp/s', function ($matches) {
            return $this->storeRawBlock("<?php{$matches[1]}?>");
        }, $value);
    }

    public function restoreRawContent($result)
    {
        $result = preg_replace_callback('/'.$this->getRawPlaceholder('(\d+)').'/', function ($matches) {
            return $this->rawBlocks[$matches[1]];
        }, $result);

        $this->rawBlocks = [];

        return $result;
    }

    public function compile($value)
    {
        $value = $this->storePhpBlocks($value);
        if (!empty($this->rawBlocks)) {
            $value = $this->restoreRawContent($value);
        }

        return $value;
    }
}

$src = <<<'BLADE'
@php
    use Filament\Support\Enums\Alignment;
    $hasNotifications = false;
@endphp

<div class="fi-no-database">
    {{ Alignment::Center }}
</div>
BLADE;

$c = new BladeRawBlock_Compiler();
$out = $c->compile($src);

if (!is_string($out)) {
    \Log::fatal('compile 应返回 string, got='.gettype($out));
}
if (str_contains($out, '@__raw_block_')) {
    \Log::fatal('raw block 占位符未被还原: '.$out);
}
if (!str_contains($out, 'use Filament\\Support\\Enums\\Alignment;')) {
    \Log::fatal('应还原 @php 内的 use Alignment: '.$out);
}
if (!str_contains($out, '<div class="fi-no-database">')) {
    \Log::fatal('HTML 根节点丢失: '.$out);
}

// 多块：占位符下标 0/1/2 都能还原
$c2 = new BladeRawBlock_Compiler();
$multi = "@php echo 0; @endphp\n@php echo 1; @endphp\n@php echo 2; @endphp";
$out2 = $c2->compile($multi);
if (str_contains($out2, '@__raw_block_')) {
    \Log::fatal('多 raw block 未还原: '.$out2);
}
if (!str_contains($out2, '<?php echo 0; ?>') || !str_contains($out2, '<?php echo 2; ?>')) {
    \Log::fatal('多 raw block 内容错误: '.$out2);
}

\Log::info('blade_raw_block_restore 测试通过');
