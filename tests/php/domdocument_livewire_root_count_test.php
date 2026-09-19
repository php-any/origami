<?php

namespace tests\php;

/**
 * 对齐 PHP DOMDocument::loadHTML 片段：body 下根元素个数。
 * Livewire debug 用它检测多根；loadHTML 不得因文档中的 ?> 截断 HTML。
 */

function livewire_root_count(string $html): int
{
    $html = preg_replace('/<script\b[^>]*>.*?<\/script>/si', '', $html);
    $html = preg_replace('/<style\b[^>]*>.*?<\/style>/si', '', $html);
    $dom = new \DOMDocument();
    $dom->loadHTML($html, defined('LIBXML_NOERROR') ? \LIBXML_NOERROR : 0);
    $body = $dom->getElementsByTagName('body')->item(0);
    if ($body === null) {
        \Log::fatal('body 缺失');
    }
    $count = 0;
    foreach ($body->childNodes as $child) {
        if ($child->nodeType == \XML_ELEMENT_NODE) {
            $count++;
        }
    }
    return $count;
}

$one = '<div class="fi-page"><div class="inner"><span>x</span></div></div>';
$c1 = livewire_root_count($one);
if ($c1 !== 1) {
    \Log::fatal('单根嵌套 div 应计 1，实际 ' . $c1);
}

$withPi = '<?xml encoding="UTF-8"?><div class="fi-page"><p>ok</p></div>';
$c2 = livewire_root_count($withPi);
if ($c2 !== 1) {
    \Log::fatal('XML PI 后单根应计 1，实际 ' . $c2);
}

$withClose = '<div class="fi-page"><button x-data="state ? true : false">t</button><span>tail</span></div>';
$c3 = livewire_root_count($withClose);
if ($c3 !== 1) {
    \Log::fatal('属性含问号的单根应计 1，实际 ' . $c3);
}

$fakeClose = '<div class="fi-page"><div x-show="a ?>">keep</div><p>still inside</p></div>';
$c4 = livewire_root_count($fakeClose);
if ($c4 !== 1) {
    \Log::fatal('HTML 中部 ?> 不得截断外层，应计 1，实际 ' . $c4);
}

\Log::info('domdocument_livewire_root_count 测试通过');
