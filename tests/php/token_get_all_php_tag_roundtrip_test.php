<?php

namespace tests\php;

/**
 * Blade compileString 依赖 token_get_all 再把 token 拼回去。
 * <?php / ?> 必须原样往返，短标签不能把 <?php 拆成 <?p + hp。
 */

function TokenRoundtrip_join($tokens)
{
    $out = '';
    foreach ($tokens as $t) {
        $out .= is_array($t) ? $t[1] : $t;
    }
    return $out;
}

$samples = [
    '<?php echo 1; ?>',
    "<?php\n\$a = 1;\n?>",
    '<?php if(true): ?><!--[if BLOCK]><![endif]--><?php endif; ?>',
    "<?php\n    \$resultCount = (\$records instanceof LengthAwarePaginator)\n        ? \$records->total()\n        : count(\$records);\n?>",
    '<?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><!--[if BLOCK]><![endif]--><?php endif; ?>',
];

foreach ($samples as $i => $src) {
    $got = TokenRoundtrip_join(token_get_all($src));
    if ($got !== $src) {
        \Log::fatal("token_get_all 往返 #{$i} 失败 got=".var_export($got, true).' want='.var_export($src, true));
    }
}

$tokens = token_get_all('<?php echo 1; ?>');
$open = $tokens[0];
if (!is_array($open) || $open[0] !== T_OPEN_TAG || !str_starts_with($open[1], '<?php')) {
    \Log::fatal('T_OPEN_TAG 应为 <?php..., 实际 '.var_export($open, true));
}
if (str_contains($open[1], 'hp') && $open[1] !== '<?php' && $open[1] !== '<?php ' && $open[1] !== "<?php\n") {
    \Log::fatal('T_OPEN_TAG 含异常 hp: '.var_export($open[1], true));
}

\Log::info('token_get_all php tag roundtrip 测试通过');
