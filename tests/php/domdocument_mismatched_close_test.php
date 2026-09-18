<?php

namespace tests\php;

/**
 * 含不匹配闭合标签的 HTML 解析不得死循环（Livewire debug DOMDocument 路径）。
 */
$html = '<div class="a"><span>x</span></p><p>y</p></div>';
$dom = new \DOMDocument();
$ok = @$dom->loadHTML($html, LIBXML_NOERROR);
if ($ok !== true) {
    \Log::fatal('loadHTML 应成功');
}
$body = $dom->getElementsByTagName('body')->item(0);
if ($body === null) {
    \Log::fatal('缺少 body');
}
$count = 0;
foreach ($body->childNodes as $child) {
    $count++;
    if ($count > 100) {
        \Log::fatal('childNodes 遍历疑似死循环');
    }
}
\Log::info('domdocument_mismatched_close 测试通过 count='.$count);
