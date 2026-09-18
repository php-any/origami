<?php

namespace tests\php;

/**
 * preg_replace replacement 中未知 \X 须保留反斜杠（对齐 PHP），
 * 否则 Livewire 编译视图会把 \Livewire\Features\... 写成 LivewireFeatures...。
 */

$c = '<?php $component->withAttributes([]); ?>';
$pat = '/(<\?php\s+\$component->withAttributes\(\[.*?\]\);\s*\?>)/s';
$rep = "$1\n<?php \Livewire\Features\SupportCompiledWireKeys\SupportCompiledWireKeys::processComponentKey(\$component); ?>\n";
$out = preg_replace($pat, $rep, $c);

if (! str_contains($out, '\\Livewire\\Features\\SupportCompiledWireKeys\\SupportCompiledWireKeys')) {
    Log::fatal('preg_replace replacement 丢失命名空间反斜杠: ' . $out);
}
if (str_contains($out, 'LivewireFeaturesSupportCompiledWireKeys')) {
    Log::fatal('preg_replace 不应吞掉 FQN 反斜杠');
}

Log::info('preg_replace replacement 反斜杠保留测试通过');
