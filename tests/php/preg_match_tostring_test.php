<?php

namespace tests\php;

/**
 * preg_match 的 subject 应对带 __toString 的对象调用魔法方法。
 * Livewire insertAttributesIntoHtmlRoot 会对 HtmlString 做 preg_match 找根标签。
 */

class PregMatch_HtmlStringLike
{
    public function __toString(): string
    {
        return "<div class=\"root\">hello</div>";
    }
}

$n = preg_match('/(?:\n\s*|^\s*)<([a-zA-Z0-9\-]+)/', new PregMatch_HtmlStringLike(), $matches, PREG_OFFSET_CAPTURE);
if ($n !== 1) {
    \Log::fatal('preg_match 未对 __toString 对象匹配根标签, n='.var_export($n, true));
}
if (!isset($matches[1][0]) || $matches[1][0] !== 'div') {
    \Log::fatal('根标签名错误: '.var_export($matches, true));
}

$patched = substr_replace(new PregMatch_HtmlStringLike(), ' wire:id="x"', (int) $matches[1][1] + strlen($matches[1][0]), 0);
if (strpos($patched, '<div wire:id="x"') === false) {
    \Log::fatal('substr_replace 未对 __toString 对象生效: '.$patched);
}

\Log::info('preg_match __toString 测试通过');
