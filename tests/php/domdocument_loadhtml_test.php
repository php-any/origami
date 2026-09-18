<?php

namespace tests\php;

/**
 * DOMDocument::loadHTML 不应挂死（Livewire debug 多根检测会调用）。
 */
$html = '<div wire:id="x">hello</div>';
$dom = new \DOMDocument();
$ok = @$dom->loadHTML($html, LIBXML_NOERROR);
if ($ok !== true && $ok !== false) {
    \Log::fatal('loadHTML 返回异常: '.var_export($ok, true));
}
$body = $dom->getElementsByTagName('body')->item(0);
if ($body === null) {
    \Log::fatal('body 节点缺失');
}
$len = 0;
if (isset($body->childNodes)) {
    $len = $body->childNodes->length;
}
\Log::info('domdocument_loadhtml 测试通过 body_children='.$len);
